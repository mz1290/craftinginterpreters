CC = gcc
CFLAGS = -g -Wall

all: clox

leak: clox
	leaks --atExit -- ./clox

clox: main.c memory.o chunk.o debug.o value.o vm.o compiler.o scanner.o
	$(CC) $(CFLAGS) memory.o chunk.o debug.o value.o vm.o compiler.o scanner.o main.c -o clox

scanner.o: scanner.c
	$(CC) -c $(CFLAGS) scanner.c -o scanner.o

compiler.o: compiler.c
	$(CC) -c $(CFLAGS) compiler.c -o compiler.o

vm.o: vm.c
	$(CC) -c $(CFLAGS) vm.c -o vm.o

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