#ifndef clox_object_h
#define clox_object_h

#include "common.h"
#include "chunk.h"
#include "table.h"
#include "value.h"

#define OBJ_TYPE(value)        (AS_OBJ(value)->type)

#define IS_CLASS(value)        isObjType(value, OBJ_CLASS)
#define IS_CLOSURE(value)      isObjType(value, OBJ_CLOSURE)
#define IS_FUNCTION(value)     isObjType(value, OBJ_FUNCTION)
#define IS_INSTANCE(value)     isObjType(value, OBJ_INSTANCE)
#define IS_NATIVE(value)       isObjType(value, OBJ_NATIVE)
#define IS_STRING(value)       isObjType(value, OBJ_STRING)

#define AS_CLASS(value)        ((ObjClass*)AS_OBJ(value))
#define AS_CLOSURE(value)      ((ObjClosure*)AS_OBJ(value))
#define AS_FUNCTION(value)     ((ObjFunction*)AS_OBJ(value))
#define AS_INSTANCE(value)     ((ObjInstance*)AS_OBJ(value))
#define AS_NATIVE(value)       (((ObjNative*)AS_OBJ(value))->function)
// Take a Value that is expected to contain a pointer to a valid ObjString on
// the heap. AS_STRING return the ObjString pointer and the AS_CSTRING returns
// the character array.
#define AS_STRING(value)       ((ObjString*)AS_OBJ(value))
#define AS_CSTRING(value)      (((ObjString*)AS_OBJ(value))->chars)

typedef enum {
    OBJ_CLASS,
    OBJ_CLOSURE,
    OBJ_FUNCTION,
    OBJ_INSTANCE,
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
    int          upvalueCount;
} ObjClosure;

typedef struct {
    Obj        obj;
    ObjString* name;
} ObjClass;

typedef struct {
    Obj       obj;
    ObjClass* klass;
    Table     fields;
} ObjInstance;


ObjClass* newClass(ObjString*);
ObjClosure* newClosure(ObjFunction*);
ObjFunction* newFunction();
ObjInstance* newInstance(ObjClass*);
ObjNative* newNative(NativeFn);
ObjString* takeString(char*, int);
ObjString* copyString(const char*, int);
ObjUpvalue* newUpvalue(Value*);
void printObject(Value);

static inline bool isObjType(Value value, ObjType type) {
    return IS_OBJ(value) && AS_OBJ(value)->type == type;
}

#endif