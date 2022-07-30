#include <stdlib.h>
#include <string.h>

#include "memory.h"
#include "object.h"
#include "table.h"
#include "value.h"

#define TABLE_MAX_LOAD 0.75


void initTable(Table* table) {
    table->count = 0;
    table->capacity = 0;
    table->entries = NULL;
}

void freeTable(Table* table) {
    FREE_ARRAY(Entry, table->entries, table->capacity);
    initTable(table);
}

// Takes a key and an array of buckets, and figures out which bucket the entry
// belongs in. Can be used to look up existing entries in the hash table and to
// decide where to insert new ones.
static Entry* findEntry(Entry* entries, int capacity, ObjString* key) {
    // Find index with hash key
    uint32_t index = key->hash % capacity;

    // The idea is to use this variable to track when we pass tombstones. In
    // scenarios where we are inserting a new key we will eventually hit an
    // empty bucket when probing. If we had passed a tombstone before the empty
    // bucket we want to use that instead.
    Entry* tombstone = NULL;

    for (;;) {
        // Get the current entry at probe
        Entry* entry = &entries[index];

        // Return match or empty entry so caller can handle accordingly
        if (entry->key == NULL) {
            if (IS_NIL(entry->value)) {
                // Empty entry. Did we pass tombstone?
                return tombstone != NULL ? tombstone : entry;
            } else {
                // Found a tombstone. Update our tracker variable.
                if (tombstone == NULL) tombstone = entry;
            }
        } else if (entry->key == key) {
            // Found the key
            return entry;
        }

        // The current bucket entry has an existing key that is different. To
        // resolve collisions we use this to being linear probing with the for
        // loop. Note that infinite loops should not be possible due to the load
        // factor used in clox.
        index = (index + 1) % capacity;
    }
}

bool tableGet(Table* table, ObjString* key, Value* value) {
    // If the table is empty, return early
    if (table->count == 0) return false;

    // Get the result from table
    Entry* entry = findEntry(table->entries, table->capacity, key);

    // If empty, retun false
    if (entry->key == NULL) return false;

    // If not empty, copy the value to output parameter so caller can use it
    *value = entry->value;

    return true;
}

static void adjustCapacity(Table* table, int capacity) {
    // Allocate a new array
    Entry* entries = ALLOCATE(Entry, capacity);

    // Initialize each element in new array
    for (int i = 0; i < capacity; i++) {
        entries[i].key = NULL;
        entries[i].value = NIL_VAL;
    }

    // Prevent collisions in new table by re-inserting existing entries
    table->count = 0;
    for (int i = 0; i < table->capacity; i++) {
        Entry* entry = &table->entries[i];
        if (entry->key == NULL) continue;

        Entry* dest = findEntry(entries, capacity, entry->key);
        dest->key = entry->key;
        dest->value = entry->value;

        // Count real entries and ignore previous tombstones
        table->count++;
    }

    // Release memory from old table
    FREE_ARRAY(Entry, table->entries, table->capacity);

    // Store the new array in the 'hash' struct to represent new storage
    table->entries = entries;
    table->capacity = capacity;
}

bool tableSet(Table* table, ObjString* key, Value value) {
    // If we don’t have enough capacity to insert an item, we reallocate and
    // grow the array.
    if (table->count + 1 > table->capacity * TABLE_MAX_LOAD) {
        int capacity = GROW_CAPACITY(table->capacity);
        adjustCapacity(table, capacity);
    }

    Entry* entry = findEntry(table->entries, table->capacity, key);
    bool isNewKey = entry->key == NULL;

    // If the value is not NULL then it was a tombstone. Counting in that case
    // would be a double count so we can ignore. Important for maintaining load
    // factor.
    if (isNewKey && IS_NIL(entry->value)) table->count++;

    entry->key = key;
    entry->value = value;
    return isNewKey;
}

bool tableDelete(Table* table, ObjString* key) {
    if (table->count == 0) return false;

    // Find the entry
    Entry* entry = findEntry(table->entries, table->capacity, key);
    if (entry->key == NULL) return false;

    // Place a tombstone in the entry
    entry->key = NULL;
    entry->value = BOOL_VAL(true);
    return true;
}

void tableAddAll(Table* from, Table* to) {
    for (int i = 0; i < from->capacity; i++) {
        Entry* entry = &from->entries[i];

        if (entry->key != NULL) {
            tableSet(to, entry->key, entry->value);
        }
    }
}

ObjString* tableFindString(Table* table, const char* chars, int length,
    uint32_t hash) {
    // Faile early if table is empty
    if (table->count == 0) return NULL;

    // Find the bucket index for given hash
    uint32_t index = hash % table->capacity;

    for (;;) {
        // Get current bucket entry
        Entry* entry = &table->entries[index];

        if (entry->key == NULL) {
            // Stop if we find an empty non-tombstone entry.
            if (IS_NIL(entry->value)) return NULL;
        } else if (entry->key->length == length &&
            entry->key->hash == hash &&
            memcmp(entry->key->chars, chars, length) == 0) {
                // Found entry after deep comparison
                return entry->key;
            }

        index = (index + 1) % table->capacity;
    }
}

void tableRemoveWhite(Table* table) {
    for (int i = 0; i < table->capacity; i++) {
        Entry* entry = &table->entries[i];

        if (entry->key != NULL && !entry->key->obj.isMarked) {
            tableDelete(table, entry->key);
        }
    }
}

void markTable(Table* table) {
    // walk the array and mark each value AND key string
    for (int i = 0; i < table->capacity; i++) {
        Entry* entry = &table->entries[i];
        markObject((Obj*)entry->key);
        markValue(entry->value);
    }
}