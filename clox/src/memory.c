#include <stdlib.h>

#include "compiler.h"
#include "memory.h"
#include "vm.h"

// For GC logging
#include <stdio.h>
#include "debug.h"

#define GC_HEAP_GROW_FACTOR 2


void* reallocate(void* pointer, size_t oldSize, size_t newSize) {
    vm.bytesAllocated += newSize - oldSize;

    if (newSize > oldSize && DEBUG_LOX & DF_STRESS_GC) {
        collectGarbage();
    }

    if (vm.bytesAllocated > vm.nextGC) {
        collectGarbage();
    }

    if (newSize == 0) {
        free(pointer);
        return NULL;
    }

    // https://man7.org/linux/man-pages/man3/realloc.3p.html
    void* result = realloc(pointer, newSize);

    // realloc() can fail in a scenario where there is not enough memory. We
    // must handle this scenario and exit accordingly.
    if (result == NULL) exit(1);

    return result;
}

void markObject(Obj* object) {
    if (object == NULL) return;

    // Prevent adding of already processed objects
    if (object->isMarked) return;

    if (DEBUG_LOX & DF_LOG_GC) {
        printf("%p mark ", (void*)object);
        printValue(OBJ_VAL(object));
        printf("\n");
    }

    object->isMarked = true;

    // Update worklist with object
    if (vm.grayCapacity < vm.grayCount + 1) {
        vm.grayCapacity = GROW_CAPACITY(vm.grayCapacity);
        vm.grayStack = (Obj**)realloc(vm.grayStack,
            sizeof(Obj*) * vm.grayCapacity);

        if (vm.grayStack == NULL) exit(1);
    }

    vm.grayStack[vm.grayCount++] = object;
}

void markValue(Value value) {
    if (IS_OBJ(value)) markObject(AS_OBJ(value));
}

static void markArray(ValueArray* array) {
    for (int i = 0; i < array->count; i++) {
        markValue(array->values[i]);
    }
}

static void blackenObject(Obj* object) {
    if (DEBUG_LOX & DF_LOG_GC) {
        printf("%p blacken ", (void*)object);
        printValue(OBJ_VAL(object));
        printf("\n");
    }

    switch (object->type) {
    case OBJ_BOUND_METHOD: {
        ObjBoundMethod* bound = (ObjBoundMethod*)object;
        markValue(bound->receiver);
        markObject((Obj*)bound->method);
        break;
    }
    case OBJ_CLASS: {
        ObjClass* klass = (ObjClass*)object;
        markObject((Obj*)klass->name);
        markTable(&klass->methods);
        break;
    }
    case OBJ_CLOSURE: {
        ObjClosure* closure = (ObjClosure*)object;
        markObject((Obj*)closure->function);
        for (int i = 0; i < closure->upvalueCount; i++) {
            markObject((Obj*)closure->upvalues[i]);
        }
        break;
    }
    case OBJ_FUNCTION: {
        ObjFunction* function = (ObjFunction*)object;
        markObject((Obj*)function->name);
        markArray(&function->chunk.constants);
        break;
    }
    case OBJ_INSTANCE: {
        ObjInstance* instance = (ObjInstance*)object;
        markObject((Obj*)instance->klass);
        markTable(&instance->fields);
        break;
    }
    case OBJ_UPVALUE:
        markValue(((ObjUpvalue*)object)->closed);
        break;
    case OBJ_NATIVE:
    case OBJ_STRING:
        break;
    }
}

// This enables us to handle freeing of different kind of objects and there
// unique implementations that may require memory allocation.
static void freeObject(Obj* object) {
    if (DEBUG_LOX & DF_LOG_GC) {
        printf("%p free type %d\n", (void*)object, object->type);
    }

    switch (object->type) {
    case OBJ_BOUND_METHOD:
        FREE(ObjBoundMethod, object);
        break;
    case OBJ_CLASS: {
        ObjClass* klass = (ObjClass*)object;
        freeTable(&klass->methods);
        FREE(ObjClass, object);
        break;
    } 
    case OBJ_CLOSURE: {
        ObjClosure* closure = (ObjClosure*)object;
        FREE_ARRAY(ObjUpvalue*, closure->upvalues, closure->upvalueCount);
        FREE(ObjClosure, object);
        break;
    }
    case OBJ_FUNCTION: {
        ObjFunction* function = (ObjFunction*)object;
        freeChunk(&function->chunk);
        FREE(ObjFunction, object);
        break;
    }
    case OBJ_INSTANCE: {
        ObjInstance* instance = (ObjInstance*)object;
        freeTable(&instance->fields);
        FREE(ObjInstance, object);
        break;
    }
    case OBJ_NATIVE:
        FREE(ObjNative, object);
        break;
    case OBJ_STRING: {
        ObjString* string = (ObjString*)object;
        FREE_ARRAY(char, string->chars, string->length + 1);
        FREE(ObjString, object);
        break;
    }
    case OBJ_UPVALUE:
        FREE(ObjUpvalue, object);
        break;
    }
}

static void markRoots() {
    // Traverse the stack and mark local variables and temporaries
    for (Value* slot = vm.stack; slot < vm.stackTop; slot++) {
        markValue(*slot);
    }

    // Traverse each call frame for constants and upvalues
    for (int i = 0; i < vm.frameCount; i++) {
        markObject((Obj*)vm.frames[i].closure);
    }

    // Traverse open upvalue list
    for (ObjUpvalue* upvalue = vm.openUpvalues; upvalue != NULL;
        upvalue = upvalue->next) {
        markObject((Obj*)upvalue);
    }

    // Traverse global variables
    markTable(&vm.globals);

    // Traverse roots the compiler uses
    markCompilerRoots();
}

static void traceReferences() {
    // Percolate through object graph, visiting references and adding new
    // objects to worklist for processing as needed
    while (vm.grayCount > 0) {
        Obj* object = vm.grayStack[--vm.grayCount];
        blackenObject(object);
    }
}

static void sweep() {
    Obj* previous = NULL;
    Obj* object = vm.objects;

    while (object != NULL) {
        if (object->isMarked) {
            // Reset marked status for next GC cycle
            object->isMarked = false;

            previous = object;
            object = object->next;
        } else {
            // Garbage collect this object
            Obj* unreached = object;
            object = object->next;

            // Unlink the object
            if (previous != NULL) {
                previous->next = object;
            } else {
                vm.objects = object;
            }

            // Release the object memory
            freeObject(unreached);
        }
    }
}

void freeObjects() {
    Obj* object = vm.objects;

    while (object != NULL) {
        Obj* next = object->next;
        freeObject(object);
        object = next;
    }

    free(vm.grayStack);
}

void collectGarbage() {
    size_t before;

    if (DEBUG_LOX & DF_LOG_GC) {
        printf("-- gc begin\n");
        before = vm.bytesAllocated;
    }

    // Traverse program object graph, markign each object and adding to worklist
    markRoots();

    // Process references from each object in worklist
    traceReferences();

    // Delete intern string entries in hash table before occupied strings are
    // collected by GC
    tableRemoveWhite(&vm.strings);

    // Collect
    sweep();

    // Adjust GC threshold
    vm.nextGC = vm.bytesAllocated * GC_HEAP_GROW_FACTOR;

    if (DEBUG_LOX & DF_LOG_GC) {
        printf("-- gc end\n");
        printf("   collected %zu bytes (from %zu to %zu) next at %zu\n",
            before - vm.bytesAllocated, before, vm.bytesAllocated, vm.nextGC);
    }
}