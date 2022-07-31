#ifndef clox_vm_h
#define clox_vm_h

#include "object.h"
#include "table.h"
#include "value.h"

// Maximum call depth supported
#define FRAMES_MAX 64
#define STACK_MAX (FRAMES_MAX * UINT8_COUNT)


// A call frame represents a single ongoing function call
typedef struct {
    // Pointer to clousre, containing the function, being called
    ObjClosure* closure;

    // IP for VM to return to upon function completion
    uint8_t*     ip;

    // Points into the VM’s value stack at the first slot that this function can
    // use.
    Value*       slots;
} CallFrame;

// Bytecode VM architecture for state:
// - Local variables and temporaries are on the stack
// - Globals are in a hash table
// - Variables in closures use upvalues
typedef struct {
    // We access each function's byte code chunk through it's call frame.
    CallFrame   frames[FRAMES_MAX];

    // Height of the call frame stack
    int         frameCount;

    Value       stack[STACK_MAX];
    Value*      stackTop;
    Table       globals;
    Table       strings;
    ObjString*  initString; // object initializer "init"
    ObjUpvalue* openUpvalues;

    size_t      bytesAllocated;
    size_t      nextGC;
    Obj*        objects;

    // Garbage collection worklist stack
    int grayCount;
    int grayCapacity;
    Obj** grayStack;
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