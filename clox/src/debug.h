#ifndef clox_debug_h
#define clox_debug_h

#include <ctype.h>
#include <string.h>
#include <stdio.h>

#include "chunk.h"
#include "value.h"

int DEBUG_LOX;

// https://stackoverflow.com/questions/1631266/flags-enum-c
typedef enum {
    DF_SCANNING = 1 << 0,
    DF_CODE     = 1 << 1,
    DF_TRACE    = 1 << 2,
} DebugFlag;

void SetDebug(char*);

// An offset is simply the number of bytes from the beginning of the chunk.

// Disassemble all of the instructions in the entire chunk.
void disassembleChunk(Chunk*, const char*);

// Disassemble single instruction, return offset of next instruction.
int disassembleInstruction(Chunk*, int);

#endif