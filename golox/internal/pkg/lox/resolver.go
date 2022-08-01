package lox

import (
	"github.com/mz1290/golox/internal/pkg/ast"
	"github.com/mz1290/golox/internal/pkg/common"
	"github.com/mz1290/golox/internal/pkg/token"
)

type FunctionType byte

const (
	FT_NONE FunctionType = iota
	FT_FUNCTION
	FT_INITIALIZER
	FT_METHOD
)

type ClassType byte

const (
	CT_NONE ClassType = iota
	CT_CLASS
	CT_SUBCLASS
)

// The resolver drives our semantic analysis to figure out variable bindings.
// The parser creates a syntax tree that this resolver crawls. After resolving,
// the interpreter takes the same tree, crawls, and processes each node.
//
// Where the parser's job is to determine if a program is 'grammatically'
// correct (syntatic analysis), it is the resolvers job to figure out what
// pieces of the program actually mean (semantic analysis).
//
// Our resolver only cares about the following kinds of nodes:
// - block statement: introduces a new scope for the statements it conatins.
// - function declaration: introduces a new scope for its body and binds its
//   parameters in that scope.
// - variable declaration: adds a new variable to the current scope.
// - any variable and assignment expressions: must have their variables resolved
//
// Each time Resolver visits a variable, it tells the interpreter how many
// scopes there are between the current scope and the scope where the variable
// is defined. At runtime, this number represents the number of environments
// between the current one and the enclosing one where the interpreter can find
// the variable's value. This is done using Interpreter.Resolve().
type Resolver struct {
	runtime         *Lox
	interpreter     *Interpreter
	currentFunction FunctionType
	currentClass    ClassType

	// scopes stack is only used for local block scopes. Global variables are
	// not tracked by the Resolver. If a variable cannot be fonud in scopes,
	// then we assume it is global.
	scopes *common.Stack
}

type Scope map[string]bool

func NewResolver(l *Lox, i *Interpreter) *Resolver {
	return &Resolver{
		runtime:         l,
		interpreter:     i,
		currentFunction: FT_NONE,
		currentClass:    CT_NONE,
		scopes:          common.NewStack(),
	}
}

func (r *Resolver) Resolve(statements []ast.Stmt) {
	r.resolveStatements(statements)
}

func (r *Resolver) resolveStatements(statements []ast.Stmt) {
	for _, stmt := range statements {
		r.resolveStatement(stmt)
	}
}

func (r *Resolver) resolveStatement(statement ast.Stmt) {
	statement.Accept(r)
}

func (r *Resolver) resolveExpression(expression ast.Expr) {
	expression.Accept(r)
}

func (r *Resolver) resolveFunction(function ast.Function, ftype FunctionType) {
	enclosingFunction := r.currentFunction
	r.currentFunction = ftype

	// Create new scope for function body (static analysis)
	r.beginScope()

	for _, param := range function.Params {
		r.declare(param)
		r.define(param)
	}

	r.resolveStatements(function.Body)

	r.endScope()
	r.currentFunction = enclosingFunction
}

func (r *Resolver) beginScope() {
	scope := make(Scope)
	r.scopes.Push(scope)
}

func (r *Resolver) endScope() {
	r.scopes.Pop()
}

func (r *Resolver) declare(name *token.Token) {
	if r.scopes.Len() == 0 {
		return
	}

	// Get the scope from stack
	scope := r.scopes.Peek().(Scope)

	// If variable already declared in local scope, collision. Report error.
	if _, ok := scope[name.Lexeme]; ok {
		r.runtime.ErrorTokenMessage(name, "already a variable with this name "+
			"in this scope")
		return
	}

	// false shows that we have not finished resolving the variables initializer
	scope[name.Lexeme] = false
}

func (r *Resolver) define(name *token.Token) {
	if r.scopes.Len() == 0 {
		return
	}

	// Get the scope from stack
	scope := r.scopes.Peek().(Scope)

	// true shows that the variable is fully initialized
	scope[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr ast.Expr, name *token.Token) {
	// Iterate linked list from back-to-front (top of stack is last)
	for n, distance := r.scopes.Top(), 0; distance < r.scopes.Len(); n, distance = n.Next(), distance+1 {
		// Cast element value as Scope
		scope := n.Value().(Scope)

		if _, ok := scope[name.Lexeme]; ok {
			r.interpreter.Resolve(expr, distance)
			return
		}
	}
}

func (r *Resolver) VisitBlockStmt(stmt ast.Block) (interface{}, error) {
	// Begins a new scope
	r.beginScope()

	// Traverse the statements within the block
	r.resolveStatements(stmt.Statements)

	// Discard the scope
	r.endScope()

	return nil, nil
}

func (r *Resolver) VisitClassStmt(stmt ast.Class) (interface{}, error) {
	enclosingClass := r.currentClass
	r.currentClass = CT_CLASS

	r.declare(stmt.Name)
	r.define(stmt.Name)

	// Superclass is not pointer to struct, this checks if struct is not empty
	if stmt.Superclass != (ast.Variable{}) {
		// Confirm that there is no cycle in the inheritance chain
		if stmt.Name.Lexeme == stmt.Superclass.Name.Lexeme {
			r.runtime.ErrorTokenMessage(stmt.Superclass.Name,
				"a class can't inherit from itself")
			return nil, nil
		}

		r.currentClass = CT_SUBCLASS
		r.resolveExpression(stmt.Superclass)
	}

	if stmt.Superclass != (ast.Variable{}) {
		// Begin a new scope for defining super
		r.beginScope()
		// Get the scope from stack and add "super"
		scope := r.scopes.Peek().(Scope)
		scope["super"] = true
	}

	// Begin a new scope for defning "this"
	r.beginScope()

	// Get the scope from stack and add "this"
	scope := r.scopes.Peek().(Scope)
	scope["this"] = true

	for _, method := range stmt.Methods {
		declaration := FT_METHOD
		if method.Name.Lexeme == "init" {
			declaration = FT_INITIALIZER
		}

		r.resolveFunction(method, declaration)
	}

	// Discard "this" scope
	r.endScope()

	// If superclass existed, we need to remove the super scope
	if stmt.Superclass != (ast.Variable{}) {
		r.endScope()
	}

	// Rever class type back to previous state
	r.currentClass = enclosingClass

	return nil, nil
}

func (r *Resolver) VisitFunctionStmt(stmt ast.Function) (interface{}, error) {
	// Declare and define the function name in current scope
	r.declare(stmt.Name)
	r.define(stmt.Name)

	// Resolve function body. Declaring and defining before this allows the
	// function to recursively refer to istelf inside its own body.
	r.resolveFunction(stmt, FT_FUNCTION)

	return nil, nil
}

func (r *Resolver) VisitVarStmt(stmt ast.Var) (interface{}, error) {
	r.declare(stmt.Name)

	if stmt.Initializer != nil {
		r.resolveExpression(stmt.Initializer)
	}

	r.define(stmt.Name)

	return nil, nil
}

func (r *Resolver) VisitAssignExpr(expr ast.Assign) (interface{}, error) {
	r.resolveExpression(expr.Value)
	r.resolveLocal(expr, expr.Name)
	return nil, nil
}

func (r *Resolver) VisitVariableExpr(expr ast.Variable) (interface{}, error) {
	if r.scopes.Len() != 0 {
		// Get the scope from stack
		scope := r.scopes.Peek().(Scope)

		// Get the variable initialized status from the scope
		if varStatus, ok := scope[expr.Name.Lexeme]; ok {
			// If we have declared but not yet initialized a variable, report error
			if varStatus == false {
				r.runtime.ErrorTokenMessage(expr.Name, "can't read local variable in "+
					"its own initializer")
				return nil, nil
			}
		}
	}

	r.resolveLocal(expr, expr.Name)

	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(stmt ast.Expression) (interface{}, error) {
	r.resolveExpression(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitIfStmt(stmt ast.If) (interface{}, error) {
	r.resolveExpression(stmt.Condition)

	r.resolveStatement(stmt.ThenBranch)

	if stmt.ElseBranch != nil {
		r.resolveStatement(stmt.ElseBranch)
	}

	return nil, nil
}

func (r *Resolver) VisitPrintStmt(stmt ast.Print) (interface{}, error) {
	r.resolveExpression(stmt.Expression)
	return nil, nil
}

func (r *Resolver) VisitReturnStmt(stmt ast.Return) (interface{}, error) {
	// Check if we are inside a function
	if r.currentFunction == FT_NONE {
		r.runtime.ErrorTokenMessage(stmt.Keyword, "can't return from "+
			"top-level code")
		return nil, nil
	}

	if stmt.Value != nil {
		if r.currentFunction == FT_INITIALIZER {
			r.runtime.ErrorTokenMessage(stmt.Keyword, "can't return a value "+
				"from an initializer")
		}

		r.resolveExpression(stmt.Value)
	}

	return nil, nil
}

func (r *Resolver) VisitWhileStmt(stmt ast.While) (interface{}, error) {
	r.resolveExpression(stmt.Condition)
	r.resolveStatement(stmt.Body)
	return nil, nil
}

func (r *Resolver) VisitBinaryExpr(expr ast.Binary) (interface{}, error) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitCallExpr(expr ast.Call) (interface{}, error) {
	r.resolveExpression(expr.Callee)

	for _, arg := range expr.Arguments {
		r.resolveExpression(arg)
	}

	return nil, nil
}

func (r *Resolver) VisitGetExpr(expr ast.Get) (interface{}, error) {
	r.resolveExpression(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitGroupingExpr(expr ast.Grouping) (interface{}, error) {
	r.resolveExpression(expr.Expression)
	return nil, nil
}

func (r *Resolver) VisitLiteralExpr(expr ast.Literal) (interface{}, error) {
	return nil, nil
}

func (r *Resolver) VisitLogicalExpr(expr ast.Logical) (interface{}, error) {
	r.resolveExpression(expr.Left)
	r.resolveExpression(expr.Right)
	return nil, nil
}

func (r *Resolver) VisitSetExpr(expr ast.Set) (interface{}, error) {
	r.resolveExpression(expr.Value)
	r.resolveExpression(expr.Object)
	return nil, nil
}

func (r *Resolver) VisitSuperExpr(expr ast.Super) (interface{}, error) {
	if r.currentClass == CT_NONE {
		r.runtime.ErrorTokenMessage(expr.Keyword, "can't use 'super' outside "+
			"of a class")
	} else if r.currentClass != CT_SUBCLASS {
		r.runtime.ErrorTokenMessage(expr.Keyword, "can't use 'super' in a "+
			"class with no superclass")
	}

	// Treat 'super' token as a variable. Resolve and store the number of hops
	// through the environment chain the interpreter needs to take to find the
	// environment that contains the superclass.
	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitThisExpr(expr ast.This) (interface{}, error) {
	if r.currentClass == CT_NONE {
		r.runtime.ErrorTokenMessage(expr.Keyword, "can't use \"this\" outside "+
			"of a class")
	}

	r.resolveLocal(expr, expr.Keyword)
	return nil, nil
}

func (r *Resolver) VisitUnaryExpr(expr ast.Unary) (interface{}, error) {
	r.resolveExpression(expr.Right)
	return nil, nil
}
