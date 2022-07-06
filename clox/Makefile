CC = gcc
CFLAGS = -g -Wall

clox: main.c memory.o chunk.o debug.o value.o
	$(CC) $(CFLAGS) memory.o chunk.o debug.o value.o main.c -o clox

value.o: value.c
	$(CC) -c $(CFLAGS) value.c -o value.o

debug.o: debug.c
	$(CC) -c $(CFLAGS) debug.c -o debug.o

memory.o: memory.c
	$(CC) -c $(CFLAGS) memory.c -o memory.o

chunk.o: chunk.c
	$(CC) -c $(CFLAGS) chunk.c -o chunk.o

clean:
	rm -f *.o