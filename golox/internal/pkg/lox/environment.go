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

func (e *Environment) GetAt(distance int, name string) interface{} {
	return e.ancestor(distance).Values[name]
}

func (e *Environment) AssignAt(distance int, name *token.Token, value interface{}) {
	e.ancestor(distance).Values[name.Lexeme] = value
}

func (e *Environment) ancestor(distance int) *Environment {
	environment := e

	for i := 0; i < distance; i++ {
		environment = environment.Enclosing
	}

	return environment
}
