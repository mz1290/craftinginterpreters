#ifndef clox_table_h
#define clox_table_h

#include "common.h"
#include "value.h"


typedef struct {
    ObjString* key;
    Value value;
} Entry;

typedef struct {
    int count;
    int capacity;
    Entry* entries;
} Table;

// Initialize hash table.
void initTable(Table*);

// Unallocate hash table.
void freeTable(Table*);

// Get the value based on provided key.
bool tableGet(Table*, ObjString*, Value*);

// Adds the given key/value pair to the given hash table.
bool tableSet(Table*, ObjString*, Value);

// Remove entry from hash table.
bool tableDelete(Table*, ObjString*);

// Copies all of the entries of one hash table into another.
void tableAddAll(Table*, Table*);

// Check if string exists in master set.
ObjString* tableFindString(Table*, const char*, int, uint32_t);

// Visit each entry in table and deletes entries that are not marked active.
// This prevents dangling pointers in the hash table once the occupied string
// has been garbage collected.
void tableRemoveWhite(Table* table);

// Garbage collector helper for marking global heap variables
void markTable(Table* table);

#endif