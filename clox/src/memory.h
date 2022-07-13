#ifndef clox_memory_h
#define clox_memory_h

#include "common.h"
#include "object.h"

#define GROW_DEFAULT    8
#define GROW_FACTOR     2

// Allocate a new array on the heap, big enough for the string’s characters and
// the trailing terminator.
#define ALLOCATE(type, count) \
    (type*)reallocate(NULL, 0, sizeof(type) * (count))

// Calculate a new capacity based on input capacity. Scales current capacity by
// a GROW_FACTOR. If 0 elements, initialize at GROW_DEFAULT elements.
#define GROW_CAPACITY(capacity) \
    ((capacity) < GROW_DEFAULT ? GROW_DEFAULT : (capacity) * GROW_FACTOR)

// Grow array to specified size. This macro takes care of getting the size of
// the array's element type and casting the resulting void* back to a pointer
// of the right type. Wrapper of call to our reallocate().
#define GROW_ARRAY(type, pointer, oldCount, newCount) \
    (type*)reallocate(pointer, sizeof(type) * (oldCount), \
        sizeof(type) * (newCount))

// Wrapper to reallocate(). Frees memory by passing '0' as the new size arg.
#define FREE_ARRAY(type, pointer, oldCount) \
    reallocate(pointer, sizeof(type) * (oldCount), 0)

// Primary function used for memory management. Provides: allocation, freeing,
// and changing sizes of existing allocations.
// The two size arguments (left=old, right=new) dictate which operation is
// performed:
//    operation         old         new
// 1. allocate new      0           non-zero
// 2. free allocation   non-zero    0
// 3. shrink existing   non-zero    < old
// 4. grow existing     non-zero    > old
void* reallocate(void*, size_t, size_t);

#endif