package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/errors"
	"github.com/mz1290/golox/internal/pkg/token"
)

// Every instance in Lox is similiar to Python. Every instance is an open
// collection of named values. Methods on the instance’s class can access and
// modify properties, but so can outside code.
type Instance struct {
	runtime *Lox
	Klass   *Class

	// Instance is responsible for storing state.
	Fields map[string]interface{}
}

func NewInstance(l *Lox, klass *Class) *Instance {
	return &Instance{
		runtime: l,
		Klass:   klass,
		Fields:  make(map[string]interface{}),
	}
}

func (i *Instance) String() string {
	return fmt.Sprintf("%s instance", i.Klass)
}

func (i *Instance) Get(name *token.Token) interface{} {
	if val, ok := i.Fields[name.Lexeme]; ok {
		return val
	}

	// If we did not find a matching field, check the Class's methods
	method := i.Klass.FindMethod(name.Lexeme)
	if method != nil {
		return method.Bind(i)
	}

	i.runtime.RuntimeError(errors.RuntimeError.New(name,
		fmt.Sprintf("undefined variable property %s", name.Lexeme)))

	return nil
}

func (i *Instance) Set(name *token.Token, value interface{}) {
	i.Fields[name.Lexeme] = value
}

func IsInstance(object interface{}) bool {
	_, ok := object.(*Instance)
	return ok
}
