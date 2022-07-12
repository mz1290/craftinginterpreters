#ifndef clox_vm_h
#define clox_vm_h

#include <stdio.h>
#include <stdarg.h>

#include "chunk.h"
#include "common.h"
#include "debug.h"
#include "value.h"
#include "compiler.h"

#define STACK_MAX 256


typedef struct {
    Chunk* chunk;

    // IP always points to the next instruction, not the one currently being
    // handled
    uint8_t* ip;

    Value stack[STACK_MAX];
    Value* stackTop;
} VM;

// Results the VM will use to handle exiting scenarios
typedef enum {
    INTERPRET_OK,
    INTERPRET_COMPILE_ERROR,
    INTERPRET_RUNTIME_ERROR
} InterpretResult;

void initVM();
void freeVM();

// Entrypoint into the VM
InterpretResult interpret(const char*);

// Stack operations
void push(Value value);
Value pop();

#endif