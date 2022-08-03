#include <stdio.h>
#include <stdarg.h>
#include <string.h>
#include <time.h>

#include "common.h"
#include "compiler.h"
#include "debug.h"
#include "object.h"
#include "memory.h"
#include "vm.h"

VM vm;

static Value clockNative(int argCount, Value* args) {
    return NUMBER_VAL((double)clock() / CLOCKS_PER_SEC);
}

static void resetStack() {
    vm.stackTop = vm.stack;
    vm.frameCount = 0;
    vm.openUpvalues = NULL;
}

static void runtimeError(const char* format, ...) {
    CallFrame* frame = &vm.frames[vm.frameCount - 1];
    ObjFunction* function = frame->closure->function;
    size_t instruction = frame->ip - function->chunk.code - 1;
    fprintf(stderr, "[line %d] RuntimeError: ",
        function->chunk.lines[instruction]);

    va_list args;
    va_start(args, format);
    vfprintf(stderr, format, args);
    va_end(args);
    fputs("\n", stderr);

    resetStack();
}

static void defineNative(const char* name, NativeFn function) {
    push(OBJ_VAL(copyString(name, (int)strlen(name))));
    push(OBJ_VAL(newNative(function)));
    tableSet(&vm.globals, AS_STRING(vm.stack[0]), vm.stack[1]);
    pop();
    pop();
}

void initVM() {
    resetStack();
    vm.objects = NULL;
    vm.bytesAllocated = 0;
    vm.nextGC = 1024 * 1024;

    vm.grayCount = 0;
    vm.grayCapacity = 0;
    vm.grayStack = NULL;

    initTable(&vm.globals);
    initTable(&vm.strings);

    // Zero field before allocation to prevent GC
    vm.initString = NULL;
    vm.initString = copyString("init", 4);

    defineNative("clock", clockNative);
}

void freeVM() {
    freeTable(&vm.globals);
    freeTable(&vm.strings);
    vm.initString = NULL; // freeObjects() will take care of this
    freeObjects();
}

void push(Value value) {
    // Assign value to open top slot
    *vm.stackTop = value;

    // Advance top so points to where the next value to be pushed will go
    vm.stackTop++;
}

Value pop() {
    // Reduce top stack to last populated spot and indicate this slot is now
    // available (next slot to be filled)
    vm.stackTop--;

    // Return the popped value to caller
    return *vm.stackTop;
}

static Value peek(int distance) {
    return vm.stackTop[-1 - distance];
}

static bool call(ObjClosure* closure, int argCount) {
    if (argCount != closure->function->arity) {
        runtimeError("expected %d arguments but got %d",
            closure->function->arity, argCount);
        return false;
    }

    if (vm.frameCount == FRAMES_MAX) {
        runtimeError("stack overflow");
        return false;
    }

    CallFrame* frame = &vm.frames[vm.frameCount++];
    frame->closure = closure;
    frame->ip = closure->function->chunk.code;
    frame->slots = vm.stackTop - argCount - 1;
    return true;
}

static bool callValue(Value callee, int argCount) {
    if (IS_OBJ(callee)) {
        switch (OBJ_TYPE(callee)) {
        case OBJ_BOUND_METHOD: {
            ObjBoundMethod* bound = AS_BOUND_METHOD(callee);
            // Top of the stack will contain all of the method args. Directly
            // under the args is the closure of the called method. Within that
            // closure is how we can access "slot zero" in the call frame and
            // update its value to be the receiver. The -argcount and -1 is
            // pointer arithmetic since stackTop is a pointer to an array
            // element itself. -argcount lets us skip past all the arguments and
            // then -1 adjusts for stackTop pointing to 1 past the last used
            // stack slot.
            vm.stackTop[-argCount - 1] = bound->receiver;
            return call(bound->method, argCount);
        }
        case OBJ_CLASS: {
            ObjClass* klass = AS_CLASS(callee);
            vm.stackTop[-argCount - 1] = OBJ_VAL(newInstance(klass));

            Value initializer;
            if (tableGet(&klass->methods, vm.initString, &initializer)) {
                return call(AS_CLOSURE(initializer), argCount);
            } else if (argCount != 0) {
                runtimeError("expected 0 arguments but got %d", argCount);
                return false;
            }

            return true;
        }
        case OBJ_CLOSURE: // by default we treat all functions as closures
            return call(AS_CLOSURE(callee), argCount);
        case OBJ_NATIVE: {
            NativeFn native = AS_NATIVE(callee);
            Value result = native(argCount, vm.stackTop - argCount);
            vm.stackTop -= argCount + 1;
            push(result);
            return true;
        }
        default:
            break; // Non-callable object type.
        }
    }

    runtimeError("can only call functions and classes");
    return false;
}

static bool invokeFromClass(ObjClass* klass, ObjString* name, int argCount) {
    Value method;

    if (!tableGet(&klass->methods, name, &method)) {
        runtimeError("undefined property \"%s\"", name->chars);
        return false;
    }

    return call(AS_CLOSURE(method), argCount);
}

static bool invoke(ObjString* name, int argCount) {
    // We can access the receiver from the stack by peeking back the number of
    // arguments (since those are in front of it on the stack)
    Value receiver = peek(argCount);

    if (!IS_INSTANCE(receiver)) {
        runtimeError("only instances have methods");
        return false;
    }

    ObjInstance* instance = AS_INSTANCE(receiver);

    Value value;
    if (tableGet(&instance->fields, name, &value)) {
        vm.stackTop[-argCount - 1] = value;
        return callValue(value, argCount);
    }

    return invokeFromClass(instance->klass, name, argCount);
}

static bool bindMethod(ObjClass* klass, ObjString* name) {
    Value method;

    // Check class's method table for a name match
    if (!tableGet(&klass->methods, name, &method)) {
        runtimeError("undefined property \"%s\"", name->chars);
        return false;
    }

    // Method found for name, wrap method in ObjBoundMethod
    ObjBoundMethod* bound = newBoundMethod(peek(0), AS_CLOSURE(method));

    // Remove the instance from top of stack and replace with the bound method
    pop();
    push(OBJ_VAL(bound));

    return true;
}

static ObjUpvalue* captureUpvalue(Value* local) {
    ObjUpvalue* prevUpvalue = NULL;
    ObjUpvalue* upvalue = vm.openUpvalues;

    while (upvalue != NULL && upvalue->location > local) {
        prevUpvalue = upvalue;
        upvalue = upvalue->next;
    }

    if (upvalue != NULL && upvalue->location == local) {
        return upvalue;
    }

    ObjUpvalue* createdUpvalue = newUpvalue(local);

    createdUpvalue->next = upvalue;
    if (prevUpvalue == NULL) {
        vm.openUpvalues = createdUpvalue;
    } else {
        prevUpvalue->next = createdUpvalue;
    }

    return createdUpvalue;
}

static void closeUpvalues(Value* last) {
    while (vm.openUpvalues != NULL && vm.openUpvalues->location >= last) {
        ObjUpvalue* upvalue = vm.openUpvalues;

        // Get the value by dereferencing current upvalue at location
        upvalue->closed = *upvalue->location;

        // Update upvalue location to be current "closed" address
        upvalue->location = &upvalue->closed;

        vm.openUpvalues = upvalue->next;
    }
}

static void defineMethod(ObjString* name) {
    // Get the method closure from top of stack
    Value method = peek(0);

    // Get the class the method will be bound to (second from top of stack)
    ObjClass* klass = AS_CLASS(peek(1));

    // Store the closure in the class's method table
    tableSet(&klass->methods, name, method);

    // Remove the closure from the stack
    pop();
}

static bool isFalsey(Value value) {
    return IS_NIL(value) || (IS_BOOL(value) && !AS_BOOL(value));
}

static void concatenate() {
    ObjString* b = AS_STRING(peek(0));
    ObjString* a = AS_STRING(peek(1));

    int length = a->length + b->length;
    char* chars = ALLOCATE(char, length + 1);
    memcpy(chars, a->chars, a->length);
    memcpy(chars + a->length, b->chars, b->length);
    chars[length] = '\0';

    ObjString* result = takeString(chars, length);
    pop();
    pop();
    push(OBJ_VAL(result));
}

// Internal helper function that runs the bytecode instructions.
static InterpretResult run() {
    // Store the current topmost call frame in a local variable
    CallFrame* frame = &vm.frames[vm.frameCount - 1];

// Reads the byte currently pointed at by ip and then advances the instruction
// pointer
#define READ_BYTE() (*frame->ip++)

// Reads 2 bytes from chunk and builds a 16-bit unsigned int out of them.
#define READ_SHORT() (frame->ip += 2, (uint16_t)((frame->ip[-2] << 8) | frame->ip[-1]))

// Reads the next byte from the bytecode, treats the resulting number as an
// index, and looks up the corresponding Value in the chunk’s constant table.
#define READ_CONSTANT() (frame->closure->function->chunk.constants.values[READ_BYTE()])

// Reads one-byte operand from the bytecode chunk. Treats that as an index into
// the chunk’s constant table and returns the string at that index. Does not
// check that the value is a string — it just indiscriminately casts it since
// the clox compiler never emits an instruction that refers to a non-string
// constant.
#define READ_STRING() AS_STRING(READ_CONSTANT())

// This is a nifty preprocessor trick. Passing an operator is valid because the
// preprocessor only cares about tokens. The do-while is necessary for a macro
// to include multiple statements inside a block and also allow a semicolon at
// the end.
#define BINARY_OP(valueType, op) \
    do { \
      if (!IS_NUMBER(peek(0)) || !IS_NUMBER(peek(1))) { \
        runtimeError("operands must be numbers"); \
        return INTERPRET_RUNTIME_ERROR; \
      } \
      double b = AS_NUMBER(pop()); \
      double a = AS_NUMBER(pop()); \
      push(valueType(a op b)); \
    } while (false)

    // Read and execute a single bytecode instruction
    for (;;) {
        if (DEBUG_LOX & DF_TRACE) {
            // Show the current contents of VM stack
            printf("          ");
            for (Value* slot = vm.stack; slot < vm.stackTop; slot++) {
                printf("[ ");
                printValue(*slot);
                printf(" ]");
            }
            printf("\n");

            // disassembleInstruction expects an integer byte offset, we must
            // convert ip back to a relative offset.
            disassembleInstruction(&frame->closure->function->chunk,
                (int)(frame->ip - frame->closure->function->chunk.code));
        }

        uint8_t instruction;

        // Decode the instruction
        switch (instruction = READ_BYTE()) {
        case OP_CONSTANT: {
            Value constant = READ_CONSTANT();
            push(constant);
            break;
        }
        case OP_NIL:      push(NIL_VAL); break;
        case OP_TRUE:     push(BOOL_VAL(true)); break;
        case OP_FALSE:    push(BOOL_VAL(false)); break;
        case OP_POP:      pop(); break;
        case OP_GET_LOCAL: {
            uint8_t slot = READ_BYTE();
            push(frame->slots[slot]);
            break;
        }
        case OP_SET_LOCAL: {
            uint8_t slot = READ_BYTE();
            frame->slots[slot] = peek(0);
            // Note that we do not pop the value from the VM stack. Assignment
            // is an expression so it must produce a value. The value of an
            // assginment *is* the assigned value, so the VM just leaves the
            // original value on the stack.
            break;
        }
        case OP_GET_GLOBAL: {
            ObjString* name = READ_STRING();
            Value value;
            if (!tableGet(&vm.globals, name, &value)) {
                runtimeError("undefined variable \"%s\"", name->chars);
                return INTERPRET_RUNTIME_ERROR;
            }

            push(value);
            break;
        }
        case OP_DEFINE_GLOBAL: {
            ObjString* name = READ_STRING();
            tableSet(&vm.globals, name, peek(0));
            pop();
            break;
        }
        case OP_SET_GLOBAL: {
            ObjString* name = READ_STRING();
            if (tableSet(&vm.globals, name, peek(0))) {
                tableDelete(&vm.globals, name); 
                runtimeError("undefined variable \"%s\"", name->chars);
                return INTERPRET_RUNTIME_ERROR;
            }
            break;
        }
        case OP_GET_UPVALUE: {
            uint8_t slot = READ_BYTE();
            push(*frame->closure->upvalues[slot]->location);
            break;
        }
        case OP_SET_UPVALUE: {
            uint8_t slot = READ_BYTE();
            *frame->closure->upvalues[slot]->location = peek(0);
            break;
        }
        case OP_GET_PROPERTY: {
            if (!IS_INSTANCE(peek(0))) {
                runtimeError("only instances have properties");
                return INTERPRET_RUNTIME_ERROR;
            }

            // Left expression has already been executed and resulting instance
            // is on top of the stack. Get it here.
            ObjInstance* instance = AS_INSTANCE(peek(0));

            // Get the field name from the constant pool
            ObjString* name = READ_STRING();

            // Look up field in instance's field table
            Value value;
            if (tableGet(&instance->fields, name, &value)) {
                // Found the field, pop the instance
                pop();

                // Push the found value on stack as result
                push(value);
                break;
            }

            // No field with provided property name was found. Check if name
            // refers to a method.
            if (!bindMethod(instance->klass, name)) {
                // Field does not exist in instance
                return INTERPRET_RUNTIME_ERROR;
            }

            break;
        }
        case OP_SET_PROPERTY: {
            if (!IS_INSTANCE(peek(1))) {
                runtimeError("only instances have fields");
                return INTERPRET_RUNTIME_ERROR;
            }

            // The state of the stack:
            // Top   = value to be stored as field
            // Top-1 = instance whose field is being set

            // Get the instance from stack
            ObjInstance* instance = AS_INSTANCE(peek(1));

            // Get the field name from the constant pool
            ObjString* name = READ_STRING();

            // Store value from top of stack into the instance's field table
            tableSet(&instance->fields, name, peek(0));

            // Pop stored value from top of stack
            Value value = pop();

            // Pop the instance
            pop();

            // Push value back on stack as top. Basically, we removed the second
            // stack element.
            push(value);
            break;
        }
        case OP_GET_SUPER: {
            ObjString* name = READ_STRING();
            ObjClass* superclass = AS_CLASS(pop());

            if (!bindMethod(superclass, name)) {
                return INTERPRET_RUNTIME_ERROR;
            }
            break;
        }
        case OP_EQUAL: {
            Value b = pop();
            Value a = pop();
            push(BOOL_VAL(valuesEqual(a, b)));
            break;
        }
        case OP_GREATER:  BINARY_OP(BOOL_VAL, >); break;
        case OP_LESS:     BINARY_OP(BOOL_VAL, <); break;
        case OP_ADD: {
            if (IS_STRING(peek(0)) && IS_STRING(peek(1))) {
                concatenate();
            } else if (IS_NUMBER(peek(0)) && IS_NUMBER(peek(1))) {
                double b = AS_NUMBER(pop());
                double a = AS_NUMBER(pop());
                push(NUMBER_VAL(a + b));
            } else {
                runtimeError(
                    "operands must be two numbers or two strings");
                return INTERPRET_RUNTIME_ERROR;
            }
            break;
        }
        case OP_SUBTRACT: BINARY_OP(NUMBER_VAL, -); break;
        case OP_MULTIPLY: BINARY_OP(NUMBER_VAL, *); break;
        case OP_DIVIDE:   BINARY_OP(NUMBER_VAL, /); break;
        case OP_NOT:
            push(BOOL_VAL(isFalsey(pop())));
            break;
        case OP_NEGATE:
            if (!IS_NUMBER(peek(0))) {
                runtimeError("operand must be a number");
                return INTERPRET_RUNTIME_ERROR;
            }
            push(NUMBER_VAL(-AS_NUMBER(pop())));
            break;
        case OP_PRINT: {
            // The code for evaluating the expression has already run, we just
            // need to pop the result from top of stack.
            printValue(pop());
            printf("\n");
            break;
        }
        case OP_JUMP: {
            uint16_t offset = READ_SHORT();
            frame->ip += offset;
            break;
        }
        case OP_JUMP_IF_FALSE: {
            uint16_t offset = READ_SHORT();
            // Check the condition value on top of stack and handle instruction
            // pointer accordingly with jump offset
            if (isFalsey(peek(0))) frame->ip += offset;
            break;
        }
        case OP_LOOP: {
            uint16_t offset = READ_SHORT();
            frame->ip -= offset;
            break;
        }
        case OP_CALL: {
            int argCount = READ_BYTE();
            if (!callValue(peek(argCount), argCount)) {
                return INTERPRET_RUNTIME_ERROR;
            }
            frame = &vm.frames[vm.frameCount - 1];
            break;
        }
        case OP_INVOKE: {
            ObjString* method = READ_STRING();
            int argCount = READ_BYTE();

            if (!invoke(method, argCount)) {
                return INTERPRET_RUNTIME_ERROR;
            }

            // If invocation succeeded then there is a new call frame on the
            // stack. We must refresh the VM run() cached copy.
            frame = &vm.frames[vm.frameCount - 1];

            break;
        }
        case OP_CLOSURE: {
            ObjFunction* function = AS_FUNCTION(READ_CONSTANT());
            ObjClosure* closure = newClosure(function);
            push(OBJ_VAL(closure));

            // Where closures come to life in runtime
            for (int i = 0; i < closure->upvalueCount; i++) {
                uint8_t isLocal = READ_BYTE();
                uint8_t index = READ_BYTE();

                if (isLocal) {
                    // Upvalue closes over local variable in enclosing function
                    closure->upvalues[i] = captureUpvalue(frame->slots + index);
                } else {
                    // Capture upvalue from surrounding (current) function
                    closure->upvalues[i] = frame->closure->upvalues[index];
                }
            }
            break;
        }
        case OP_CLOSE_UPVALUE:
            closeUpvalues(vm.stackTop - 1);
            pop();
            break;
        case OP_RETURN: {
            Value result = pop();
            closeUpvalues(frame->slots);
            vm.frameCount--;

            // Check if final stack frame and need to return
            if (vm.frameCount == 0) {
                pop();
                return INTERPRET_OK;
            }

            vm.stackTop = frame->slots;
            push(result);
            frame = &vm.frames[vm.frameCount - 1];
            break;
        }
        case OP_CLASS:
            // Load class name string from constant table and use to create new
            // class object with given name
            push(OBJ_VAL(newClass(READ_STRING())));
            break;
        case OP_METHOD:
            defineMethod(READ_STRING());
            break;
        case OP_INHERIT: {
            Value superclass = peek(1);
            if (!IS_CLASS(superclass)) {
                runtimeError("superclass must be a class");
                return INTERPRET_RUNTIME_ERROR;
            }

            ObjClass* subclass = AS_CLASS(peek(0));
            tableAddAll(&AS_CLASS(superclass)->methods, &subclass->methods);
            pop(); // Subclass
            break;
        }
        }
    }

#undef READ_BYTE
#undef READ_SHORT
#undef READ_CONSTANT
#undef READ_STRING
#undef BINARY_OP
}


InterpretResult interpret(const char* source) {
    // Pass source code to compiler and get back an object function containing
    // the compiled top-level code.
    ObjFunction* function = compile(source);

    // Check if the compiler encouterned error
    if (function == NULL) return INTERPRET_COMPILE_ERROR;

    // Store object function on the stack
    push(OBJ_VAL(function));

    // Set up call frame to execute object function instructions
    ObjClosure* closure = newClosure(function);
    pop();
    push(OBJ_VAL(closure));
    call(closure, 0);

    return run();
}