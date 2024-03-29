#ifndef clox_debug_h
#define clox_debug_h

#include "chunk.h"


// https://stackoverflow.com/questions/1631266/flags-enum-c
typedef enum {
    DF_SCANNING  = 1 << 0,
    DF_CODE      = 1 << 1,
    DF_TRACE     = 1 << 2,
    DF_STRESS_GC = 1 << 3,
    DF_LOG_GC    = 1 << 4,
} DebugFlag;

void SetDebug(char*);
int GetDebug();

// An offset is simply the number of bytes from the beginning of the chunk.

// Disassemble all of the instructions in the entire chunk.
void disassembleChunk(Chunk*, const char*);

// Disassemble single instruction, return offset of next instruction.
int disassembleInstruction(Chunk*, int);

#endif