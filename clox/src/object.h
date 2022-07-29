#ifndef clox_object_h
#define clox_object_h

#include "common.h"
#include "chunk.h"
#include "value.h"

#define OBJ_TYPE(value)        (AS_OBJ(value)->type)

#define IS_CLOSURE(value)      isObjType(value, OBJ_CLOSURE)
#define IS_FUNCTION(value)     isObjType(value, OBJ_FUNCTION)
#define IS_NATIVE(value)       isObjType(value, OBJ_NATIVE)
#define IS_STRING(value)       isObjType(value, OBJ_STRING)

#define AS_CLOSURE(value)      ((ObjClosure*)AS_OBJ(value))
// Take a Value that is expected to contain a pointer to a valid ObjString on
// the heap. AS_STRING return the ObjString pointer and the AS_CSTRING returns
// the character array.
#define AS_FUNCTION(value)     ((ObjFunction*)AS_OBJ(value))
#define AS_NATIVE(value)       (((ObjNative*)AS_OBJ(value))->function)
#define AS_STRING(value)       ((ObjString*)AS_OBJ(value))
#define AS_CSTRING(value)      (((ObjString*)AS_OBJ(value))->chars)

typedef enum {
    OBJ_CLOSURE,
    OBJ_FUNCTION,
    OBJ_NATIVE,
    OBJ_STRING,
    OBJ_UPVALUE
} ObjType;

struct Obj {
    ObjType     type;
    bool        isMarked;
    struct Obj* next;
};

typedef struct {
    Obj        obj;
    int        arity;
    int        upvalueCount;
    Chunk      chunk;
    ObjString* name;
} ObjFunction;

typedef Value (*NativeFn)(int argCount, Value* args);

typedef struct {
    Obj      obj;
    NativeFn function;
} ObjNative;


struct ObjString {
    Obj      obj; // provides required state to be an "Obj"
    int      length;
    char*    chars;
    uint32_t hash;
};

typedef struct ObjUpvalue {
    Obj                obj;
    Value*             location;
    Value              closed;
    struct ObjUpvalue* next;
}   ObjUpvalue;

typedef struct {
    Obj          obj;
    ObjFunction* function;
    ObjUpvalue** upvalues;
    int upvalueCount;
} ObjClosure;

ObjClosure* newClosure(ObjFunction*);
ObjFunction* newFunction();
ObjNative* newNative(NativeFn);
ObjString* takeString(char*, int);
ObjString* copyString(const char*, int);
ObjUpvalue* newUpvalue(Value*);
void printObject(Value);

static inline bool isObjType(Value value, ObjType type) {
    return IS_OBJ(value) && AS_OBJ(value)->type == type;
}

#endif