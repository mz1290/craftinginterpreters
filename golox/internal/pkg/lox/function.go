package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/ast"
)

// weird error i was encountering, due to passing in a literal instance of Function
// and not the pointer version.
// https://stackoverflow.com/questions/36527261/interface-conversion-panic-when-method-is-not-actually-missing

type Function struct {
	Closure       *Environment
	Declaration   *ast.Function
	isInitializer bool
}

func NewFunction(declaration *ast.Function, closure *Environment, isInitializer bool) *Function {
	return &Function{
		Closure:       closure,
		Declaration:   declaration,
		isInitializer: isInitializer,
	}
}

func (f *Function) Bind(instance *Instance) *Function {
	// Add new environment within method's original closure
	environment := NewLocalEnvironment(f.Closure)

	// Declare "this" in new env and bind to the instance that this method is
	// being accessed from.
	environment.Define("this", instance)

	// Return function that contains instance is bound as "this"
	return NewFunction(f.Declaration, environment, f.isInitializer)
}

func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	environment := NewLocalEnvironment(f.Closure)

	for i := 0; i < len(f.Declaration.Params); i++ {
		environment.Define(f.Declaration.Params[i].Lexeme, arguments[i])
	}

	_, err := interpreter.executeBlock(f.Declaration.Body, environment)
	if err != nil {
		if IsReturnable(err) {
			if f.isInitializer {
				return f.Closure.GetAt(0, "this"), nil
			}

			return err.(*Return).Value, nil
		}
	}

	if f.isInitializer {
		return f.Closure.GetAt(0, "this"), nil
	}

	return nil, err
}

func (f *Function) Arity() int {
	return len(f.Declaration.Params)
}

func (f Function) String() string {
	return fmt.Sprintf("<fn %s>", f.Declaration.Name.Lexeme)
}
