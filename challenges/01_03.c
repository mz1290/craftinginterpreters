#include <assert.h>
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

#ifndef TYPE
#define TYPE char*
#endif

struct Link
{
	TYPE value;
	struct Link*	next;
	struct Link*	prev;
};

struct LinkedList
{
	struct Link*	frontSentinel;
	struct Link*	backSentinel;
	int size;
};

void init(struct LinkedList*);
struct LinkedList* linkedListCreate();
void linkedListDestroy(struct LinkedList*);
int linkedListIsEmpty(struct LinkedList*);
void linkedListRemoveFront(struct LinkedList*);
void removeLink(struct LinkedList*, struct Link*);
void linkedListAdd(struct LinkedList*, TYPE);
void addLinkBefore(struct LinkedList*, struct Link*, TYPE);
int linkedListContains(struct LinkedList*, TYPE);

// Write functions to insert, find, and delete items from it. Test them.
void
init(struct LinkedList* list)
{
	assert(list != NULL);

	list->frontSentinel = malloc(sizeof (*list->frontSentinel));
	assert(list->frontSentinel != 0);

	list->backSentinel = malloc(sizeof (*list->backSentinel));
	assert(list->backSentinel != 0);

	list->frontSentinel->next = list->backSentinel;
	list->frontSentinel->prev = 0;

	list->backSentinel->next = 0;
	list->backSentinel->prev = list->frontSentinel;

	list->size = 0;
}

struct LinkedList*
linkedListCreate()
{
	struct LinkedList* list = malloc(sizeof (*list));
	assert(list != NULL);
	init(list);
	return (list);
}

void
linkedListDestroy(struct LinkedList* list)
{
	assert(list != NULL);
	while (!linkedListIsEmpty(list)) {
		linkedListRemoveFront(list);
	}
	free(list->frontSentinel);
	free(list->backSentinel);
	free(list);
	list = NULL;
}

int
linkedListIsEmpty(struct LinkedList* list)
{
	assert(list != NULL);

	if (list->size == 0) {
		return (1);
	} else {
		return (0);
	}
}

void
linkedListRemoveFront(struct LinkedList* list)
{
	assert(list != NULL);
	assert(!linkedListIsEmpty(list));
	struct Link *del = list->frontSentinel->next;
	removeLink(list, del);
}

void
removeLink(struct LinkedList* list, struct Link* link)
{
	assert(list != NULL);
	assert(link != NULL);

	link->prev->next = link->next;
	link->next->prev = link->prev;
	free(link->value);
	free(link);
	link = NULL;

	list->size -= 1;
}

void
linkedListAdd(struct LinkedList* list, TYPE value)
{
	assert(list != NULL);
	struct Link *curr = list->backSentinel;
	addLinkBefore(list, curr, value);
}

void
addLinkBefore(struct LinkedList* list, struct Link* link, TYPE value)
{
	assert(list != NULL);
	assert(link != NULL);

	struct Link *newLink = malloc(sizeof (*newLink));
	assert(newLink != NULL);

	int len = strlen(value);
	newLink->value = malloc(sizeof (*newLink->value) * (len + 1));
	assert(newLink->value != NULL);
	strcpy(newLink->value, value);
	newLink->next = link;
	newLink->prev = link->prev;

	link->prev->next = newLink;
	link->prev = newLink;

	list->size += 1;
}

int
linkedListContains(struct LinkedList* list, TYPE value)
{
	assert(list != 0);
	assert(!linkedListIsEmpty(list));

	struct Link *curr = list->frontSentinel->next;

	while (curr != list->backSentinel)
	{
		if (!strcmp(curr->value, value)) {
			return (1);
		}
		curr = curr->next;
	}
	return (0);
}

int
main(void)
{
	struct LinkedList* plist = linkedListCreate();
	linkedListAdd(plist, "Test string");
	printf("Contains = %d\n", linkedListContains(plist, "Test string"));
	linkedListDestroy(plist);
	return(0);
}
