package ast

import (
	"fmt"
	"testing"

	"github.com/mz1290/golox/internal/pkg/token"
)

func TestPrint(t *testing.T) {
	var expression Expr = Binary{
		left: Unary{
			operator: *token.New(token.MINUS, "-", nil, 1),
			right: Literal{
				value: 123,
			},
		},
		operator: *token.New(token.STAR, "*", nil, 1),
		right: Grouping{
			expression: Literal{
				value: 45.67,
			},
		},
	}

	p := ASTPrinter{}
	fmt.Println(p.Print(expression))
}
