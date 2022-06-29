package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/errors"
	"github.com/mz1290/golox/internal/pkg/token"
)

type Environment struct {
	runtime   *Lox
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnvironment(l *Lox) *Environment {
	return &Environment{
		runtime:   l,
		Enclosing: nil,
		Values:    make(map[string]interface{}),
	}
}

func NewLocalEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		runtime:   enclosing.runtime,
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}

func (e Environment) Get(name *token.Token) interface{} {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	e.runtime.RuntimeError(errors.RuntimeError.New(name,
		fmt.Sprintf("undefined variable %s", name.Lexeme)))
	return nil
}

func (e *Environment) Assign(name *token.Token, value interface{}) {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return
	}

	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	e.runtime.RuntimeError(errors.RuntimeError.New(name,
		fmt.Sprintf("undefined variable %s", name.Lexeme)))
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}
