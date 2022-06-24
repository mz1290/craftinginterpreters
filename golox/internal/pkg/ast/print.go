package ast

import (
	"fmt"
	"log"
	"strings"
)

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
	return p.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (p ASTPrinter) VisitGroupingExpr(expr Grouping) (interface{}, error) {
	return p.parenthesize("group", expr.expression)
}

func (p ASTPrinter) VisitLiteralExpr(expr Literal) (interface{}, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return fmt.Sprintf("%v", expr.value), nil
}

func (p ASTPrinter) VisitUnaryExpr(expr Unary) (interface{}, error) {
	return p.parenthesize(expr.operator.Lexeme, expr.right)
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
