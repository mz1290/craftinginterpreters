#ifndef clox_debug_h
#define clox_debug_h

#include "chunk.h"
#include "value.h"

// An offset is simply the number of bytes from the beginning of the chunk.

// Disassemble all of the instructions in the entire chunk.
void disassembleChunk(Chunk*, const char*);

// Disassemble single instruction, return offset of next instruction.
int disassembleInstruction(Chunk*, int);

#endif