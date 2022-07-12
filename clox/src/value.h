#ifndef clox_value_h
#define clox_value_h

#include "common.h"

typedef enum {
    VAL_BOOL,
    VAL_NIL, 
    VAL_NUMBER,
} ValueType;

// Representation of clox value
typedef struct {
    ValueType type;
    union {
      bool boolean;
      double number;
    } as; 
} Value;

// Check a clox Value's type
#define IS_BOOL(value)    ((value).type == VAL_BOOL)
#define IS_NIL(value)     ((value).type == VAL_NIL)
#define IS_NUMBER(value)  ((value).type == VAL_NUMBER)

// Unwrap clox Value and return corresponding raw C value
#define AS_BOOL(value)    ((value).as.boolean)
#define AS_NUMBER(value)  ((value).as.number)

// Promote native C value to clox Value
#define BOOL_VAL(value)   ((Value){VAL_BOOL, {.boolean = value}})
#define NIL_VAL           ((Value){VAL_NIL, {.number = 0}})
#define NUMBER_VAL(value) ((Value){VAL_NUMBER, {.number = value}})

// Dynamic array for value arrays
typedef struct {
  int capacity;
  int count;
  Value* values;
} ValueArray;

bool valuesEqual(Value a, Value b);
void initValueArray(ValueArray*);
void writeValueArray(ValueArray*, Value);
void freeValueArray(ValueArray*);
void printValue(Value);

#endif