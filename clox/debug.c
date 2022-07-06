#include <stdio.h>

#include "debug.h"

void disassembleChunk(Chunk* chunk, const char* name) {
    printf("== %s ==\n", name);

    for (int offset = 0; offset < chunk->count;) {
        offset = disassembleInstruction(chunk, offset);
    }
}

static int constantInstruction(const char* name, Chunk* chunk, int offset) {
    // Extract the constant index
    uint8_t constant = chunk->code[offset + 1];

    // Print the name of the opcode
    printf("%-16s %4d '", name, constant);

    // Use the index to lookup the actual constant value and print
    printValue(chunk->constants.values[constant]);
    printf("'\n");

    // A constant instruction consists of two bytes (opcode, operand). The +2
    // is required to return the correct offset of the next instruction.
    return offset + 2;
}

// Print name of opcode and return next byte offset.
static int simpleInstruction(const char* name, int offset) {
    printf("%s\n", name);
    return offset + 1;
}

int disassembleInstruction(Chunk* chunk, int offset) {
    // Print the offset of the instruction
    printf("%04d ", offset);

    // Print the source line of the instruction
    if (offset > 0 &&
        chunk->lines[offset] == chunk->lines[offset - 1]) {
        // A single line of source code can compile to a large sequence of
        // instructions. As a result, we print '|' when printing instructions
        // coming from the same source lines as the preceeding one.
        printf("   | ");
    } else {
        printf("%4d ", chunk->lines[offset]);
    }

    // Read opcode
    uint8_t instruction = chunk->code[offset];
    switch (instruction) {
    case OP_CONSTANT:
      return constantInstruction("OP_CONSTANT", chunk, offset);
    case OP_RETURN:
        return simpleInstruction("OP_RETURN", offset);
    default:
        printf("Unknown opcode %d\n", instruction);
        return offset + 1;
    }
}

