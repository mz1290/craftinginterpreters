package lox

// https://yourbasic.org/golang/create-error
// https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mz1290/golox/internal/pkg/errors"
	"github.com/mz1290/golox/internal/pkg/token"
)

type Lox struct {
	HadError        bool // represents syntax/static errors
	HadRuntimeError bool // errors during execution
	Interpreter     *Interpreter
}

func New() *Lox {
	l := &Lox{
		HadError:        false,
		HadRuntimeError: false,
	}

	l.Interpreter = NewInterpreter(l)
	return l
}

// Read and execute file
func (l *Lox) RunFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		os.Exit(65)
	}

	l.run(string(data))
	if l.HadError {
		os.Exit(65)
	} else if l.HadRuntimeError {
		os.Exit(70)
	}
}

// Start interactive golox prompt
func (l *Lox) RunPrompt() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			os.Exit(65)
		}

		// Check if user signaled end of session
		if line == "exit\n" {
			break
		}

		// Execute user lox statement or expression
		l.run(line)
		l.HadError = false
		l.HadRuntimeError = false
	}
}

func (l *Lox) run(source string) {
	// create a new scanner instance
	s := NewScanner(l, source)
	tokens := s.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}

	// create new parser instance
	parser := NewParser(l, tokens)
	statements := parser.Parse()

	// Stop if there was a syntax error
	if l.HadError {
		return
	}

	// Run the resolver to find variable bindings
	resolver := NewResolver(l, l.Interpreter)
	resolver.Resolve(statements)

	// Stop if there was a semantic error
	if l.HadError {
		return
	}

	// Execute/evaluate expression
	l.Interpreter.Interpret(statements)
}

//ErrorMessage prints error message as stderr
func (l *Lox) ErrorMessage(line int, message string) {
	l.report(line, "", message)
}

func (l *Lox) ErrorTokenMessage(t *token.Token, message string) {
	if t.Type == token.EOF {
		l.report(t.Line, " at end", message)
	} else {
		l.report(t.Line, fmt.Sprintf(" at %q", t.Lexeme), message)
	}
}

func (l *Lox) report(line int, where string, message string) {
	fmt.Fprintf(os.Stderr, "[line %d] Error%s: %q\n", line, where, message)
	l.HadError = true
}

func (l *Lox) RuntimeError(err error) {
	e := err.(*errors.CustomErr)
	fmt.Fprintf(os.Stderr, "[line %d] %s\n", e.Token.Line, err)
	l.HadRuntimeError = true
}
