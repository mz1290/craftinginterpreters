#ifndef clox_vm_h
#define clox_vm_h

#include "object.h"
#include "table.h"
#include "value.h"

#define STACK_MAX 256


typedef struct {
    Chunk* chunk;

    // IP always points to the next instruction, not the one currently being
    // handled
    uint8_t* ip;

    Value    stack[STACK_MAX];
    Value*   stackTop;
    Table    globals;
    Table    strings;
    Obj*     objects;
} VM;

// Results the VM will use to handle exiting scenarios
typedef enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR
} InterpretResult;

extern VM vm;

void initVM();
void freeVM();

// Entrypoint into the VM
InterpretResult interpret(const char*);

// Stack operations
void push(Value);
Value pop();

#endif