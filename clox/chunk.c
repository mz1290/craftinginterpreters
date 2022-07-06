#include <stdlib.h>

#include "chunk.h"


void initChunk(Chunk* chunk) {
    chunk->count = 0;
    chunk->capacity = 0;
    chunk->code = NULL;
    chunk->lines = NULL;
    initValueArray(&chunk->constants);
}

 void writeChunk(Chunk* chunk, uint8_t byte, int line) {
    // Verify that current array has enough room to store new byte
    if (chunk->capacity < chunk->count + 1) {
        // Grow the array to make more room
        int oldCapacity = chunk->capacity;
        chunk->capacity = GROW_CAPACITY(oldCapacity);
        chunk->code = GROW_ARRAY(uint8_t, chunk->code,
            oldCapacity, chunk->capacity);

        // Make corresponding change to lines array
        chunk->lines = GROW_ARRAY(int, chunk->lines,
            oldCapacity, chunk->capacity);
    }

    // Store provided byte in the chunk
    chunk->code[chunk->count] = byte;

    // Store line number in lines array for runtime error usage
    chunk->lines[chunk->count] = line;

    // Increment chunk count
    chunk->count++;
}

int addConstant(Chunk* chunk, Value value) {
    writeValueArray(&chunk->constants, value);
    return chunk->constants.count - 1;
}

void freeChunk(Chunk* chunk) {
    FREE_ARRAY(uint8_t, chunk->code, chunk->capacity);
    FREE_ARRAY(int, chunk->lines, chunk->capacity);
    freeValueArray(&chunk->constants);

    // Zero out memory
    initChunk(chunk);
}