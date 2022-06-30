package lox

import (
	"time"
)

type Callable interface {
	Arity() int
	Call(*Interpreter, []interface{}) (interface{}, error)
}

func IsCallable(object interface{}) bool {
	_, ok := object.(Callable)
	return ok
}

type nativeFunctionClock struct{}

func (n nativeFunctionClock) Arity() int {
	return 0
}

func (n nativeFunctionClock) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().Unix()), nil
}

func (n nativeFunctionClock) String() string {
	return "<native fn>"
}
