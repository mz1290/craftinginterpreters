# craftinginterpreters

### Overview
This repository is the product of my OSU CS 467 Capstone research project. For 
the project I read and completed Robert Nystrom's 
[Crafting Interpreters](https://craftinginterpreters.com/). This capstone 
gave me an opportunity to learn more about the concepts and implementations of 
programming language design.

Crafting Interpreters provides two implementations of a specified programming 
language called Lox. Each interpreter behaves the same and provides the same 
features. The first interpreter (`golox`) is considered 
more theory-based implementation and performance is considerably slower than the 
second interpreter (`clox`).

`golox` is implemented using Golang. This differs from the textbook version, 
`jlox`, which uses Java. The reasoning for picking Golang instead of Java was 
simply that I had more experience with Golang. As stated earlier, the author 
used this interpreter to focus on concepts related to program language 
design. `golox` uses an abstract-syntax tree to execute Lox code and leverages 
Golangs objects and garbage collection.

`clox` is implemented using C. This second implementation was focused on 
performance and understanding the "why" behind each design aspect. Unlike 
`golox`, `clox` is implemented as a virtual machine that executes compiled 
bytecode instructions. Every component of `clox` was implemented entirely by 
hand.

#### Lox Features
- A C-like scripting language
- Dynamically typed
- Memory management
- Data types
    - bool
    - numbers
    - strings
    - nil
- Expressions
    - arithmetic
    - comparisons
    - equality
    - logical
    - precedence
    - grouping
- Statements (things that produce effects rather than values)
- Variables
- Control flow
    - if-else
    - for loop
    - while loop
- Functions (with closure support)
- Classes (with inheritance)

### Required Dependencies
- Golang
- C compiler (clang or gcc)
- Make

### Usage
The repository contains both interpreters, a test suite, and a Makefile in the 
root directory containing targets for the users to easily interact with the 
project.

To get a compiled binary for each interpreter simply run `make` from the root 
directory. You will now have two compiled binaries that you can interact with.

You can run `make clean` to remove any compiled binaries and clean each 
interpreter directory accordingly.

Running `make test-[golox||clox]` will kick off the test suite for the 
specific interpreter. For example:
```bash
> make test-golox
```

To get a Read-Eval-Print Loop, or REPL, environment you can execute one of the 
compiled binaries or `make repl-[golox||clox]`. For example:
```bash
> make repl-clox
...
> print "hello clox";
hello clox
> // you can exit using CTRL-C or CTRL-D
> ^D
```

To execute a file you should run `make` to get the compiled binaries. Execute 
either binary followed by the name of the script you want to execute. For 
example:
```bash
> ./clox-1.0.0 test/benchmark/fib.lox 
true
3.53778
```

### Credits
Nystrom, R., 2015. Crafting interpreters.

