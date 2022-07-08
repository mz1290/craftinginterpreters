#include <stdlib.h>

#include "memory.h"

void* reallocate(void* pointer, size_t oldSize, size_t newSize) {
    if (newSize == 0) {
        free(pointer);
        return NULL;
    }

    // https://man7.org/linux/man-pages/man3/realloc.3p.html
    void* result = realloc(pointer, newSize);

    // realloc() can fail in a scenario where there is not enough memory. We
    // must handle this scenario and exit accordingly.
    if (result == NULL) exit(1);

    return result;
}