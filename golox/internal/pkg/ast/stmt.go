// Code generated by generateAST; DO NOT EDIT.
package ast

import (
	"github.com/mz1290/golox/internal/pkg/token"
)

type Stmt interface {
	StmtAcceptor
}

type StmtVisitor interface {
	VisitBlockStmt(stmt Block) (interface{}, error)
	VisitClassStmt(stmt Class) (interface{}, error)
	VisitExpressionStmt(stmt Expression) (interface{}, error)
	VisitFunctionStmt(stmt Function) (interface{}, error)
	VisitIfStmt(stmt If) (interface{}, error)
	VisitPrintStmt(stmt Print) (interface{}, error)
	VisitReturnStmt(stmt Return) (interface{}, error)
	VisitVarStmt(stmt Var) (interface{}, error)
	VisitWhileStmt(stmt While) (interface{}, error)
}

type StmtAcceptor interface {
	Accept(v StmtVisitor) (interface{}, error)
}

type Block struct {
	Statements []Stmt
}

func (x Block) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBlockStmt(x)
}

type Class struct {
	Name *token.Token
	Methods []Function
}

func (x Class) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitClassStmt(x)
}

type Expression struct {
	Expression Expr
}

func (x Expression) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitExpressionStmt(x)
}

type Function struct {
	Name *token.Token
	Params []*token.Token
	Body []Stmt
}

func (x Function) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitFunctionStmt(x)
}

type If struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (x If) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitIfStmt(x)
}

type Print struct {
	Expression Expr
}

func (x Print) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintStmt(x)
}

type Return struct {
	Keyword *token.Token
	Value Expr
}

func (x Return) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitReturnStmt(x)
}

type Var struct {
	Name *token.Token
	Initializer Expr
}

func (x Var) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitVarStmt(x)
}

type While struct {
	Condition Expr
	Body Stmt
}

func (x While) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitWhileStmt(x)
}

