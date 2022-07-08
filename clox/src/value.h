#ifndef clox_value_h
#define clox_value_h

#include "common.h"

// Representation of clox value
typedef double Value;

// Dynamic array for value arrays
typedef struct {
  int capacity;
  int count;
  Value* values;
} ValueArray;

void initValueArray(ValueArray*);
void writeValueArray(ValueArray*, Value);
void freeValueArray(ValueArray*);
void printValue(Value);

#endif