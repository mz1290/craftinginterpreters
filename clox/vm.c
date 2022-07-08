
#include "vm.h"

VM vm;

static void resetStack() {
    vm.stackTop = vm.stack;
}

void initVM() {
    resetStack();
}

void freeVM() {
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

// Internal helper function that runs the bytecode instructions.
static InterpretResult run() {
// Reads the byte currently pointed at by ip and then advances the instruction
// pointer
#define READ_BYTE() (*vm.ip++)

// Reads the next byte from the bytecode, treats the resulting number as an
// index, and looks up the corresponding Value in the chunk’s constant table
#define READ_CONSTANT() (vm.chunk->constants.values[READ_BYTE()])

// This is a nifty preprocessor trick. Passing an operator is valid because the
// preprocessor only cares about tokens. The do-while is necessary for a macro
// to include multiple statements inside a block and also allow a semicolon at
// the end.
#define BINARY_OP(op) \
    do { \
      double b = pop(); \
      double a = pop(); \
      push(a op b); \
    } while (false)

    // Read and execute a single bytecode instruction
    for (;;) {
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

        uint8_t instruction;

        // Decode the instruction
        switch (instruction = READ_BYTE()) {
        case OP_CONSTANT: {
            Value constant = READ_CONSTANT();
            push(constant);
            break;
        }
        case OP_ADD:      BINARY_OP(+); break;
        case OP_SUBTRACT: BINARY_OP(-); break;
        case OP_MULTIPLY: BINARY_OP(*); break;
        case OP_DIVIDE:   BINARY_OP(/); break;
        case OP_NEGATE: 
            // Pop a value from the stack, negate it, push back onto the stack
            push(-pop());
            break;
        case OP_RETURN: {
            printValue(pop());
            printf("\n");
            return INTERPRET_OK;
        }
        }
    }

#undef READ_BYTE
#undef READ_CONSTANT
#undef BINARY_OP
}

InterpretResult interpret(Chunk* chunk) {
    vm.chunk = chunk;
    vm.ip = vm.chunk->code;
    return run();
}

