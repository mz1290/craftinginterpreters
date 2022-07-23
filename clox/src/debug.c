#include <stdio.h>
#include <ctype.h>
#include <string.h>

#include "debug.h"
#include "object.h"
#include "value.h"

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

// Local variable names never get stored in the chunk, therefore, we can't show
// the variable name. The best we can do is show the slot number.
static int byteInstruction(const char* name, Chunk* chunk, int offset) {
    uint8_t slot = chunk->code[offset + 1];
    printf("%-16s %4d\n", name, slot);
    return offset + 2; 
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
    case OP_NIL:
        return simpleInstruction("OP_NIL", offset);
    case OP_TRUE:
        return simpleInstruction("OP_TRUE", offset);
    case OP_FALSE:
        return simpleInstruction("OP_FALSE", offset);
    case OP_POP:
        return simpleInstruction("OP_POP", offset);
    case OP_GET_LOCAL:
        return byteInstruction("OP_GET_LOCAL", chunk, offset);
    case OP_SET_LOCAL:
        return byteInstruction("OP_SET_LOCAL", chunk, offset);
    case OP_GET_GLOBAL:
        return constantInstruction("OP_GET_GLOBAL", chunk, offset);
    case OP_DEFINE_GLOBAL:
        return constantInstruction("OP_DEFINE_GLOBAL", chunk, offset);
    case OP_SET_GLOBAL:
        return constantInstruction("OP_SET_GLOBAL", chunk, offset);
    case OP_EQUAL:
        return simpleInstruction("OP_EQUAL", offset);
    case OP_GREATER:
        return simpleInstruction("OP_GREATER", offset);
    case OP_LESS:
        return simpleInstruction("OP_LESS", offset);
    case OP_ADD:
        return simpleInstruction("OP_ADD", offset);
    case OP_SUBTRACT:
        return simpleInstruction("OP_SUBTRACT", offset);
    case OP_MULTIPLY:
        return simpleInstruction("OP_MULTIPLY", offset);
    case OP_DIVIDE:
        return simpleInstruction("OP_DIVIDE", offset);
    case OP_NOT:
        return simpleInstruction("OP_NOT", offset);
    case OP_NEGATE:
        return simpleInstruction("OP_NEGATE", offset);
    case OP_PRINT:
        return simpleInstruction("OP_PRINT", offset);
    case OP_RETURN:
        return simpleInstruction("OP_RETURN", offset);
    default:
        printf("Unknown opcode %d\n", instruction);
        return offset + 1;
    }
}

char* debugFlagOf(DebugFlag df) {
    switch (df) {
    case DF_SCANNING:   return "scanning"; break;
    case DF_CODE:       return "code"; break;
    case DF_TRACE:      return "trace"; break;
    default:            return "";
    }
}

void SetDebug(char* settings) {
    // Get the first debug setting
    char* setting = strtok(settings, ",");

    while (setting != NULL) {
        // Lowercase the string for comparison
        for(int i = 0; setting[i]; i++){
            setting[i] = tolower(setting[i]);
        }

        if (strcmp(setting, debugFlagOf(DF_SCANNING)) == 0) {
            //printf("enabling debug scanning output\n");
            DEBUG_LOX |= DF_SCANNING;
        } else if (strcmp(setting, debugFlagOf(DF_CODE)) == 0) {
            //printf("enabling debug code output\n");
            DEBUG_LOX |= DF_CODE;
        } else if (strcmp(setting, debugFlagOf(DF_TRACE)) == 0) {
            //printf("enabling debug trace output\n");
            DEBUG_LOX |= DF_TRACE;
        }
    
        // Advance to next setting
        setting = strtok(NULL, ",");
    }
}