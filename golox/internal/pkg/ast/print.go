package ast

import (
	"fmt"
	"log"
	"strings"
)

// In this example, the printer itself is our "visitor".
// Print accepts an expression
// Calls that expression's "Accept" with our visior
// The expression then calls our visitors implementation of it's required
// visitor function so we get the desired behavior
type ASTPrinter struct{}

func (p ASTPrinter) Print(expr Expr) string {
	ret, err := expr.Accept(p)
	if err != nil {
		log.Print(err)
		return ""
	}

	str, ok := ret.(string)
	if !ok {
		log.Printf("not a string: %v", ret)
		return ""
	}

	return str
}

func (p ASTPrinter) VisitBinaryExpr(expr Binary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p ASTPrinter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p ASTPrinter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.Value), nil
}

func (p ASTPrinter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p ASTPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	w := &strings.Builder{}

	w.WriteString("(" + name)
	for _, exp := range exprs {
		w.WriteString(" ")

		ret, err := exp.Accept(p)
		if err != nil {
			return "", err
		}

		s, ok := ret.(string)
		if !ok {
			return "", fmt.Errorf("not a string: %v", p)
		}

		w.WriteString(s)
	}
	w.WriteString(")")

	return w.String(), nil
}
