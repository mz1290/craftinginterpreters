package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/ast"
)

// weird error i was encountering, due to passing in a literal instance of Function
// and not the pointer version.
// https://stackoverflow.com/questions/36527261/interface-conversion-panic-when-method-is-not-actually-missing

type Function struct {
	Closure     *Environment
	Declaration ast.Function
}

func NewFunction(declaration ast.Function, closure *Environment) *Function {
	return &Function{
		Closure:     closure,
		Declaration: declaration,
	}
}

func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewLocalEnvironment(f.Closure)

	for i := 0; i < len(f.Declaration.Params); i++ {
		environment.Define(f.Declaration.Params[i].Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(f.Declaration.Body, environment)
	if err != nil {
		if IsReturnable(err) {
			return err.(*Return).Value, nil
		}
	}

	return nil, nil
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f Function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}
