package ast

import (
	"fmt"
	"testing"

	"github.com/mz1290/golox/internal/pkg/token"
)

func TestPrint(t *testing.T) {
	var expression Expr = Binary{
		Left: Unary{
			Operator: token.New(token.MINUS, "-", nil, 1),
			Right: Literal{
				Value: 123,
			},
		},
		Operator: token.New(token.STAR, "*", nil, 1),
		Right: Grouping{
			Expression: Literal{
				Value: 45.67,
			},
		},
	}

	p := ASTPrinter{}
	fmt.Println(p.Print(expression))
}
