#ifndef clox_object_h
#define clox_object_h

#include "common.h"
#include "chunk.h"
#include "value.h"

#define OBJ_TYPE(value)        (AS_OBJ(value)->type)

#define IS_FUNCTION(value)     isObjType(value, OBJ_FUNCTION)
#define IS_STRING(value)       isObjType(value, OBJ_STRING)

// Take a Value that is expected to contain a pointer to a valid ObjString on
// the heap. AS_STRING return the ObjString pointer and the AS_CSTRING returns
// the character array.
#define AS_FUNCTION(value)     ((ObjFunction*)AS_OBJ(value))
#define AS_STRING(value)       ((ObjString*)AS_OBJ(value))
#define AS_CSTRING(value)      (((ObjString*)AS_OBJ(value))->chars)

typedef enum {
    OBJ_FUNCTION,
    OBJ_STRING,
} ObjType;

struct Obj {
    ObjType     type;
    struct Obj* next;
};

typedef struct {
    Obj        obj;
    int        arity;
    Chunk      chunk;
    ObjString* name;
} ObjFunction;

struct ObjString {
    Obj      obj; // provides required state to be an "Obj"
    int      length;
    char*    chars;
    uint32_t hash;
};

ObjFunction* newFunction();
ObjString* takeString(char*, int);
ObjString* copyString(const char*, int);
void printObject(Value);

static inline bool isObjType(Value value, ObjType type) {
    return IS_OBJ(value) && AS_OBJ(value)->type == type;
}

#endif