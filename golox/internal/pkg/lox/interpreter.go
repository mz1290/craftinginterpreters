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
	locals      map[ast.Expr]int
}

func NewInterpreter(runtime *Lox) *Interpreter {

	i := &Interpreter{
		runtime: runtime,
		globals: NewEnvironment(runtime),
		locals:  make(map[ast.Expr]int),
	}
	i.environment = i.globals

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
		}
	} else {
		if !common.IsTruthy(left) {
			return left, nil
		}
	}

	return i.evaluate(expr.Right)
}

func (i *Interpreter) VisitSetExpr(expr ast.Set) (interface{}, error) {
	// Evaluate object whose property is being set
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	// If evaluated object is no an Instance, invalid
	if !IsInstance(object) {
		return nil, errors.RuntimeError.New(expr.Name, "only instances have "+
			"fields")
	}

	// Evaluate the value being set
	value, err := i.evaluate(expr.Value)
	if err != nil {
		return nil, errors.RuntimeError.New(expr.Name,
			fmt.Sprintf("invalid assignment on %s instance", object.(*Instance)))
	}

	// Store evaluated value in instance
	object.(*Instance).Set(expr.Name, value)

	// This is a setter so don't need to return any value
	return value, nil
}

func (i *Interpreter) VisitSuperExpr(expr ast.Super) (interface{}, error) {
	// Get number of hops to superclass env
	distance := i.locals[expr]

	// Get the superclass
	superclass := i.environment.GetAt(distance, "super").(*Class)

	// A hacky way of arriving at the proper "this" for the superclass. The env
	// where "this" is bound is always +1 from the env that stores "super".
	object := i.environment.GetAt(distance-1, "this").(*Instance)

	// Look up method
	method := superclass.FindMethod(expr.Method.Lexeme)
	if method == nil {
		return nil, errors.RuntimeError.New(expr.Method,
			fmt.Sprintf("undefined property %s", expr.Method.Lexeme))
	}

	return method.Bind(object), nil
}

func (i *Interpreter) VisitThisExpr(expr ast.This) (interface{}, error) {
	return i.lookUpVariable(expr.Keyword, expr), nil
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
	return i.lookUpVariable(expr.Name, expr), nil
}

func (i *Interpreter) lookUpVariable(name *token.Token, expr ast.Expr) interface{} {
	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	}

	return i.globals.Get(name)
}

// evaluate sends the expression back into the interpreter's visitor implementation
func (i *Interpreter) evaluate(expr ast.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *Environment) (interface{}, error) {
	var val interface{}
	previous := i.environment
	defer func() { i.environment = previous }()

	i.environment = environment
	for _, stmt := range statements {
		var err error

		val, err = i.execute(stmt)
		if err != nil {
			return nil, err
		}
	}

	return val, nil
}

func (i *Interpreter) VisitBlockStmt(stmt ast.Block) (interface{}, error) {
	return i.executeBlock(stmt.Statements, NewLocalEnvironment(i.environment))
}

func (i *Interpreter) VisitClassStmt(stmt ast.Class) (interface{}, error) {
	// This two-stage variable binding process allows references to the class
	// inside its own methods.

	var superclass interface{}
	if stmt.Superclass != (ast.Variable{}) {
		var err error

		superclass, err = i.evaluate(stmt.Superclass)
		if err != nil {
			return nil, err
		}

		// Confirm that superclass expression evaluated to a class
		if !IsClass(superclass) {
			return nil, errors.RuntimeError.New(stmt.Superclass.Name,
				"superclass must be a class")
		}
	}

	// Declare class name in current env
	i.environment.Define(stmt.Name.Lexeme, nil)

	// If we created a runtime object for superclass, we need to update the
	// interpreter's current environment with a new one that stores the
	// superclass.
	if stmt.Superclass != (ast.Variable{}) {
		i.environment = NewLocalEnvironment(i.environment)
		i.environment.Define("super", superclass)
	}

	// Convert each class method into a runtime represenation (Function)
	methods := make(map[string]*Function)
	for _, method := range stmt.Methods {
		function := NewFunction(method, i.environment, method.Name.Lexeme == "init")
		methods[method.Name.Lexeme] = function
	}

	// Convert the class *syntax node* into a runtime representation of Class
	var sc *Class = nil
	if superclass != nil {
		sc = superclass.(*Class)
	}
	klass := NewClass(i.runtime, stmt.Name.Lexeme, sc, methods)

	// If we updated our superclass environment, we need to revert back to
	// previous environment.
	if superclass != nil {
		i.environment = i.environment.Enclosing
	}

	// Store the runtime oobject with previously declared env variable
	i.environment.Assign(stmt.Name, klass)

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

func (i *Interpreter) VisitGetExpr(expr ast.Get) (interface{}, error) {
	object, err := i.evaluate(expr.Object)
	if err != nil {
		return nil, err
	}

	if IsInstance(object) {
		return object.(*Instance).Get(expr.Name), nil
	}

	return nil, errors.RuntimeError.New(expr.Name, "only instances have "+
		"properties")
}

func (i *Interpreter) VisitExpressionStmt(stmt ast.Expression) (interface{}, error) {
	return i.evaluate(stmt.Expression)
}

func (i *Interpreter) VisitFunctionStmt(stmt ast.Function) (interface{}, error) {
	function := NewFunction(stmt, i.environment, false)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitIfStmt(stmt ast.If) (interface{}, error) {
	condition, _ := i.evaluate(stmt.Condition)

	if common.IsTruthy(condition) {
		return i.execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.execute(stmt.ElseBranch)
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

func (i *Interpreter) VisitReturnStmt(stmt ast.Return) (interface{}, error) {
	var value interface{}

	if stmt.Value != nil {
		var err error

		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}

	return nil, NewReturn(value)
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
			_, err := i.execute(stmt.Body)
			if err != nil {
				return nil, err
			}
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

	if distance, ok := i.locals[expr]; ok {
		i.environment.AssignAt(distance, expr.Name, value)
	} else {
		i.globals.Assign(expr.Name, value)
	}

	return value, nil
}
