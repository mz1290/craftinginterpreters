#include <stdio.h>
#include <stdarg.h>
#include <string.h>

#include "common.h"
#include "compiler.h"
#include "debug.h"
#include "object.h"
#include "memory.h"
#include "vm.h"

VM vm;

static void resetStack() {
    vm.stackTop = vm.stack;
    vm.frameCount = 0;
}

static void runtimeError(const char* format, ...) {
    // Get chunk and IP from topmost call frame on the stack
    CallFrame* frame = &vm.frames[vm.frameCount - 1];
    size_t instruction = frame->ip - frame->function->chunk.code - 1;
    int line = frame->function->chunk.lines[instruction];
    fprintf(stderr, "[line %d] RuntimeError: ", line);

    va_list args;
    va_start(args, format);
    vfprintf(stderr, format, args);
    va_end(args);
    fputs("\n", stderr);

    resetStack();
}

void initVM() {
    resetStack();
    vm.objects = NULL;

    initTable(&vm.globals);
    initTable(&vm.strings);
}

void freeVM() {
    freeTable(&vm.globals);
    freeTable(&vm.strings);
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

static bool isFalsey(Value value) {
    return IS_NIL(value) || (IS_BOOL(value) && !AS_BOOL(value));
}

static void concatenate() {
    ObjString* b = AS_STRING(pop());
    ObjString* a = AS_STRING(pop());

    int length = a->length + b->length;
    char* chars = ALLOCATE(char, length + 1);
    memcpy(chars, a->chars, a->length);
    memcpy(chars + a->length, b->chars, b->length);
    chars[length] = '\0';

    ObjString* result = takeString(chars, length);
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
// index, and looks up the corresponding Value in the chunk???s constant table.
#define READ_CONSTANT() (frame->function->chunk.constants.values[READ_BYTE()])

// Reads one-byte operand from the bytecode chunk. Treats that as an index into
// the chunk???s constant table and returns the string at that index. Does not
// check that the value is a string ??? it just indiscriminately casts it since
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
            disassembleInstruction(&frame->function->chunk,
                (int)(frame->ip - frame->function->chunk.code));
        }
/*
#ifdef DEBUG_TRACE_EXECUTION
    // Show the current contents of VM stack
    printf("          ");
    for (Value* slot = vm.stack; slot < vm.stackTop; slot++) {
      printf("[ ");
      printValue(*slot);
      printf(" ]");
    }
    printf("\n");

    // disassembleInstruction expects an integer byte offset, we must convert ip
    // back to a relative offset.
    disassembleInstruction(vm.chunk, (int)(vm.ip - vm.chunk->code));
#endif
*/

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
        case OP_RETURN: {
            // Exit interpreter
            return INTERPRET_OK;
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
    CallFrame* frame = &vm.frames[vm.frameCount++];
    frame->function = function;
    frame->ip = function->chunk.code;
    frame->slots = vm.stack;

    return run();
}