#ifndef clox_compiler_h
#define clox_compiler_h

#include <stdio.h>
#include <stdlib.h>

#include "common.h"
#include "scanner.h"
#include "vm.h"

//#ifdef DEBUG_PRINT_CODE
#include "debug.h"
//#endif

bool compile(const char*, Chunk*);

#endif