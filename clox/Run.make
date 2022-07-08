# make -f Run.make run

all: build/clox

run: all
	./build/clox

leak: build/clox
	leaks --atExit -- ./build/clox