package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/errors"
	"github.com/mz1290/golox/internal/pkg/token"
)

type Environment struct {
	runtime *Lox
	Values  map[string]interface{}
}

func NewEnvironment(l *Lox) *Environment {
	return &Environment{
		runtime: l,
		Values:  make(map[string]interface{}),
	}
}

func (e Environment) Get(name *token.Token) interface{} {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val
	}

	e.runtime.RuntimeError(errors.RuntimeError.New(name,
		fmt.Sprintf("undefined variable %s", name.Lexeme)))
	return nil
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}
