package common

import (
	"fmt"
	"reflect"

	"github.com/mz1290/golox/internal/pkg/ast"
	"github.com/mz1290/golox/internal/pkg/errors"
	"github.com/mz1290/golox/internal/pkg/token"
)

func IsDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func IsAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func IsAlphaNumeric(c byte) bool {
	return IsAlpha(c) || IsDigit(c)
}

func IsTruthy(object interface{}) bool {
	if object == nil {
		return false
	}

	// https://yourbasic.org/golang/find-type-of-object/
	switch v := object.(type) {
	case bool:
		return v
	default:
		return true
	}
}

func IsEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil {
		return false
	}

	return a == b
}

func isType(v interface{}, kind reflect.Kind) bool {
	return reflect.ValueOf(v).Kind() == kind
}

func IsString(object interface{}) bool {
	return isType(object, reflect.String)
}

func IsFloat64(object interface{}) bool {
	return isType(object, reflect.Float64)
}

func IsVariableExpression(object interface{}) bool {
	_, ok := object.(*ast.Variable)
	return ok
}

func IsGet(object interface{}) bool {
	_, ok := object.(*ast.Get)
	return ok
}

func CheckNumberOperand(operator *token.Token, operand interface{}) error {
	if IsFloat64(operand) {
		return nil
	}

	return errors.RuntimeError.New(operator, "operand must be a number")
}

func CheckNumberOperands(operator *token.Token, left, right interface{}) error {
	if IsFloat64(left) && IsFloat64(right) {
		return nil
	}

	return errors.RuntimeError.New(operator, "operands must be numbers")
}

func Stringfy(object interface{}) string {
	if object == nil {
		return "nil"
	}

	return fmt.Sprint(object)
}
