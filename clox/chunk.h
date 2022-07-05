// A "chunk" refers to a sequence of bytecode.
// "bytecode" is a linear sequence of binary instructions.
#ifndef clox_chunk_h
#define clox_chunk_h

#include "common.h"

// Each instruction has a one-byte operation code ("opcode"). This number
// controls what kind of instruction we're dealing with.
typedef enum {
    OP_RETURN,
} OpCode;

// Chunk is the data structure used to store data with instructions.
typedef struct {
    int count;
    int capacity;
    uint8_t* code;
} Chunk;

// Initialize a new chunk.
void initChunk(Chunk*);

// Appends a byte to the end of the chunk.
void writeChunk(Chunk*, uint8_t);

// Deallocates and zeroes out contents of memory.
void freeChunk(Chunk*);

#endif