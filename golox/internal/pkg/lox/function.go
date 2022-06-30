package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/ast"
)

// weird error i was encountering, due to passing in a literal instance of Function
// and not the pointer version.
// https://stackoverflow.com/questions/36527261/interface-conversion-panic-when-method-is-not-actually-missing

type Function struct {
	Declaration ast.Function
}

func NewFunction(declaration ast.Function) *Function {
	return &Function{
		Declaration: declaration,
	}
}

func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewLocalEnvironment(interpreter.globals)

	for i := 0; i < len(f.Declaration.Params); i++ {
		environment.Define(f.Declaration.Params[i].Lexeme, arguments[i])
	}

	interpreter.executeBlock(f.Declaration.Body, environment)

	return nil, nil
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f Function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}
