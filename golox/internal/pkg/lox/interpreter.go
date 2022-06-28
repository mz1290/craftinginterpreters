package lox

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/ast"
	"github.com/mz1290/golox/internal/pkg/common"
	"github.com/mz1290/golox/internal/pkg/errors"
	"github.com/mz1290/golox/internal/pkg/token"
)

type Interpreter struct {
	runtime     *Lox
	environment *Environment
}

func NewInterpreter(runtime *Lox) *Interpreter {
	return &Interpreter{
		runtime:     runtime,
		environment: NewEnvironment(runtime),
	}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) {
	for _, stmt := range statements {
		_, err := i.execute(stmt)
		if err != nil {
			i.runtime.RuntimeError(err)
			return
		}
	}
}

func (i *Interpreter) VisitLiteralExpr(expr ast.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr ast.Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr ast.Unary) (interface{}, error) {
	// Evaluate subexpression first
	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.BANG:
		return !common.IsTruthy(right), nil
	case token.MINUS:
		err := common.CheckNumberOperand(expr.Operator, right)
		if err != nil {
			return nil, err
		}
		return -1 * right.(float64), nil
	}

	// Unreachable
	return nil, errors.RuntimeError.New(nil, "unreachable")
}

func (i *Interpreter) VisitVariableExpr(expr ast.Variable) (interface{}, error) {
	return i.environment.Get(expr.Name), nil
}

// evaluate sends the expression back into the interpreter's visitor implementation
func (i *Interpreter) evaluate(expr ast.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) VisitBinaryExpr(expr ast.Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.Right)
	if err != nil {
		return nil, err
	}

	switch expr.Operator.Type {
	case token.GREATER:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case token.GREATER_EQUAL:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case token.LESS:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case token.LESS_EQUAL:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case token.MINUS:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case token.PLUS:
		// Check if expression is arithmetic
		if common.IsFloat64(left) && common.IsFloat64(right) {
			return left.(float64) + right.(float64), nil
		}

		// Check if expression is concatenaation
		if common.IsString(left) && common.IsString(right) {
			return left.(string) + right.(string), nil
		}

		return nil, errors.RuntimeError.New(expr.Operator,
			"operands must be two numbers or two strings")
	case token.SLASH:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) / right.(float64), nil
	case token.STAR:
		err := common.CheckNumberOperands(expr.Operator, left, right)
		if err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case token.BANG_EQUAL:
		return !common.IsEqual(left, right), nil
	case token.EQUAL_EQUAL:
		return common.IsEqual(left, right), nil
	}

	// Unreachable
	return nil, errors.RuntimeError.New(nil, "unreachable")
}

func (i *Interpreter) VisitExpressionStmt(stmt ast.Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitPrintStmt(stmt ast.Print) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(common.Stringfy(value))
	return nil, nil
}

func (i *Interpreter) VisitVarStmt(stmt ast.Var) (interface{}, error) {
	var value interface{}
	var err error

	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}

	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}
