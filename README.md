# craftinginterpreters

## Overview
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

![image](https://user-images.githubusercontent.com/29135072/183303132-bfcf7a27-046a-45fa-9751-0733347017ba.png)

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

## Required Dependencies
- Golang
- C compiler (clang or gcc)
- Make

## Usage
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

## Performance
After completing both interpreters I wanted to spend some time to get an idea on how exactly each interpreter compares to the other but also how it stacks up against a real world language. Benchmarking programming languages and comparing performance is a much more controversial task than benchmarking a real-world running application since the algorithm, hardware, and compiler can have a huge impact in the results. Robert Nystrom provides a handful of benchmarks in his [repository](https://github.com/munificent/craftinginterpreters), one specifically called `zoo_batch.lox`. The script simply create an instance of an object and runs the objects methods in a 10 second loop. The results of `zoo_batch.lox` is the count of batches completed. I decided to use this same approach for my comparison so I recreated the same code in Golang and C. The results of tests are below:

![image](https://user-images.githubusercontent.com/29135072/183302721-85ec665b-4b3c-41a1-b730-95c78ff4fe58.png)

![image](https://user-images.githubusercontent.com/29135072/183302735-7558e3a9-c948-4036-a9dc-c1fff970b13a.png)

![image](https://user-images.githubusercontent.com/29135072/183302728-c8bcf4de-bdb3-4f59-86c1-d7f4c70cd58a.png)

**Note:** I disabled all optimizations for the C and Golang compiled binaries.

## Credits
Nystrom, R., 2015. Crafting interpreters.

