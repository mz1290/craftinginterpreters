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
	globals     *Environment
	environment *Environment
}

func NewInterpreter(runtime *Lox) *Interpreter {

	i := &Interpreter{
		runtime:     runtime,
		globals:     NewEnvironment(runtime),
		environment: NewEnvironment(runtime),
	}

	i.globals.Define("clock", nativeFunctionClock{})
	return i
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

func (i *Interpreter) VisitLogicalExpr(expr ast.Logical) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	if err != nil {
		return nil, err
	}

	if expr.Operator.Type == token.OR {
		if common.IsTruthy(left) {
			return left, nil
		} else {
			if !common.IsTruthy(left) {
				return left, nil
			}
		}
	}

	return i.evaluate(expr.Right)
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
	return nil, errors.RuntimeError.New(expr.Operator, "unreachable")
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

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) {
	previous := i.environment

	i.environment = environment
	for _, stmt := range statements {
		i.execute(stmt)
	}

	i.environment = previous
}

func (i *Interpreter) VisitBlockStmt(stmt ast.Block) (interface{}, error) {
	i.executeBlock(stmt.Statements, NewLocalEnvironment(i.environment))
	return nil, nil
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
	return nil, errors.RuntimeError.New(expr.Operator, "unreachable")
}

func (i *Interpreter) VisitCallExpr(expr ast.Call) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}

	var arguments []interface{}
	for _, arg := range expr.Arguments {
		argRes, err := i.evaluate(arg)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, argRes)
	}

	// Confirm the object is indeed callable
	if !IsCallable(callee) {
		return nil, errors.RuntimeError.New(expr.Paren, "can only call functions "+
			"and classes")
	}

	function := callee.(Callable)
	if len(arguments) != function.Arity() {
		return nil, errors.RuntimeError.New(expr.Paren, fmt.Sprintf("expected %d "+
			"arguments but got %d", function.Arity(), len(arguments)))
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) VisitExpressionStmt(stmt ast.Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunctionStmt(stmt ast.Function) (interface{}, error) {
	function := NewFunction(stmt)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt ast.If) (interface{}, error) {
	condition, _ := i.evaluate(stmt.Condition)

	if common.IsTruthy(condition) {
		i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		i.execute(stmt.ElseBranch)
	}

	return nil, nil
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

func (i *Interpreter) VisitWhileStmt(stmt ast.While) (interface{}, error) {
	for {
		condition, err := i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}

		if common.IsTruthy(condition) {
			i.execute(stmt.Body)
		} else {
			break
		}
	}

	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr ast.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, err
	}

	i.environment.Assign(expr.Name, value)
	return value, nil
}
