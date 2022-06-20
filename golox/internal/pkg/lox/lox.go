package lox

// https://yourbasic.org/golang/create-error
// https://www.digitalocean.com/community/tutorials/creating-custom-errors-in-go

import (
	"bufio"
	"fmt"
	"os"

	"github.com/mz1290/golox/internal/pkg/scanner"
)

type LoxError struct {
	Line    int
	Where   string
	Message string
}

func (le LoxError) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %q", le.Line, le.Where, le.Message)
}

type Lox struct {
	Err *LoxError
}

func New() *Lox {
	return &Lox{}
}

// Read and execute file
func (l *Lox) RunFile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		l.Err = &LoxError{0, "", err.Error()}
		return l.Err
	}

	return l.run(string(data))
}

// Start interactive golox prompt
func (l *Lox) RunPrompt() error {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			l.Err = &LoxError{0, "", err.Error()}
			return l.Err
		}

		// Check if user signaled end of session
		if line == "exit\n" {
			break
		}

		// Execute user lox statement or expression
		if err = l.run(line); err != nil {
			fmt.Println(err)
		}
	}

	return nil
}

func (l *Lox) run(source string) error {
	// create a new scanner instance
	s := scanner.New(source)
	tokens := s.ScanTokens()

	for _, token := range tokens {
		fmt.Println(token)
	}

	// Check if scanner encountered any errors
	for _, e := range s.Errors {
		fmt.Println(e)
	}

	return nil
}
