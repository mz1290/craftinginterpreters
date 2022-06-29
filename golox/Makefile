AST_PATH := "/Users/mzimmerman/Engineering/github/craftinginterpreters/golox/internal/pkg/ast"

all: build generate

run: build generate
	./golox

test: build generate test/myscript.lox
	./golox test/myscript.lox

generate: tools/generateAST
	go run tools/generateAST/generateAST.go $(AST_PATH)

build:
	go build -gcflags=all="-N -l" -o golox main.go

clean:
	rm -f golox