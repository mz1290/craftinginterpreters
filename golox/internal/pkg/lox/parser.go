package lox

import (
	"github.com/mz1290/golox/internal/pkg/ast"
	"github.com/mz1290/golox/internal/pkg/common"
	"github.com/mz1290/golox/internal/pkg/token"
)

type ParserError error

type Parser struct {
	runtime       *Lox
	tokens        []*token.Token
	current       int
	hadParseError bool
}

func NewParser(l *Lox, tokens []*token.Token) *Parser {
	return &Parser{
		runtime:       l,
		tokens:        tokens,
		current:       0,
		hadParseError: false,
	}
}

func (p *Parser) Parse() []ast.Stmt {
	var statements []ast.Stmt

	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	return statements
}

func (p *Parser) statement() ast.Stmt {
	if p.match(token.IF) {
		return p.ifStatement()
	}

	if p.match(token.PRINT) {
		return p.printStatement()
	}

	if p.match(token.LEFT_BRACE) {
		return ast.Block{Statements: p.block()}
	}

	return p.expressionStatement()
}

func (p *Parser) ifStatement() ast.Stmt {
	_, ok := p.consume(token.LEFT_PAREN)
	if !ok {
		p.NewParserError(p.peek(), "expected '(' after 'if'")
	}

	condition := p.expression()

	_, ok = p.consume(token.RIGHT_PAREN)
	if !ok {
		p.NewParserError(p.peek(), "expected ')' after if condition")
	}

	thenBranch := p.statement()
	var elseBranch ast.Stmt
	if p.match(token.ELSE) {
		elseBranch = p.statement()
	}

	return ast.If{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() ast.Stmt {
	value := p.expression()

	_, ok := p.consume(token.SEMICOLON)
	if !ok {
		p.NewParserError(p.peek(), "expected ';' after value")
	}

	return ast.Print{Expression: value}
}

func (p *Parser) varDeclaration() ast.Stmt {
	name, ok := p.consume(token.IDENTIFIER)
	if !ok {
		p.NewParserError(p.peek(), "expected variable name")
	}

	var initializer ast.Expr
	if p.match(token.EQUAL) {
		initializer = p.expression()
	}

	p.consume(token.SEMICOLON)
	if !ok {
		p.NewParserError(p.peek(), "expected ';' after variable declaration")
	}

	return ast.Var{Name: name, Initializer: initializer}
}

func (p *Parser) expressionStatement() ast.Stmt {
	expr := p.expression()

	_, ok := p.consume(token.SEMICOLON)
	if !ok {
		p.NewParserError(p.peek(), "expected ';' after expression")
	}

	return ast.Expression{Expression: expr}
}

func (p *Parser) block() []ast.Stmt {
	var statements []ast.Stmt

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}

	_, ok := p.consume(token.RIGHT_BRACE)
	if !ok {
		p.NewParserError(p.peek(), "expected '}' after block")
	}

	return statements
}

func (p *Parser) assignment() ast.Expr {
	expr := p.or()

	if p.match(token.EQUAL) {
		equals := p.previous()
		value := p.assignment()

		if common.IsVariableExpression(expr) {
			name := expr.(ast.Variable).Name
			return ast.Assign{Name: name, Value: value}
		}

		p.NewParserError(equals, "invalid assignment target")
	}

	return expr
}

func (p *Parser) or() ast.Expr {
	expr := p.and()

	for p.match(token.OR) {
		operator := p.previous()
		right := p.and()
		expr = ast.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() ast.Expr {
	expr := p.equality()

	for p.match(token.AND) {
		operator := p.previous()
		right := p.equality()
		expr = ast.Logical{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) expression() ast.Expr {
	return p.assignment()
}

func (p *Parser) declaration() ast.Stmt {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}

	res := p.statement()
	if p.hadParseError {
		p.synchronize()
		p.hadParseError = false
		return nil
	}

	return res
}

func (p *Parser) equality() ast.Expr {
	expr := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		operator := p.previous()
		right := p.comparison()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) comparison() ast.Expr {
	expr := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS,
		token.LESS_EQUAL) {
		operator := p.previous()
		right := p.term()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) term() ast.Expr {
	expr := p.factor()

	for p.match(token.MINUS, token.PLUS) {
		operator := p.previous()
		right := p.factor()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) factor() ast.Expr {
	expr := p.unary()

	for p.match(token.SLASH, token.STAR) {
		operator := p.previous()
		right := p.unary()
		expr = ast.Binary{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) unary() ast.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return ast.Unary{Operator: operator, Right: right}
	}

	return p.primary()
}

func (p *Parser) primary() ast.Expr {
	if p.match(token.FALSE) {
		return ast.Literal{Value: false}
	}

	if p.match(token.TRUE) {
		return ast.Literal{Value: true}
	}

	if p.match(token.NIL) {
		return ast.Literal{Value: nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return ast.Literal{Value: p.previous().Literal}
	}

	if p.match(token.IDENTIFIER) {
		return ast.Variable{Name: p.previous()}
	}

	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		_, ok := p.consume(token.RIGHT_PAREN)
		if !ok {
			p.NewParserError(p.peek(), "expected ')' after expression")
		}
		return ast.Grouping{Expression: expr}
	}

	// Token can't start an expression
	p.NewParserError(p.peek(), "expected expression")
	return nil
}

func (p *Parser) match(types ...token.Type) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}

	return false
}

func (p *Parser) consume(t token.Type) (*token.Token, bool) {
	if p.check(t) {
		return p.advance(), true
	}

	return nil, false
}

func (p Parser) check(t token.Type) bool {
	if p.isAtEnd() {
		return false
	}

	return p.peek().Type == t
}

// advance() method consumes the current token and returns it
func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.current++
	}

	return p.previous()
}

// isAtEnd() checks if we’ve run out of tokens to parse
func (p Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

// peek() returns the current token we have yet to consume
func (p Parser) peek() *token.Token {
	return p.tokens[p.current]
}

// previous() returns the most recently consumed token
func (p Parser) previous() *token.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) NewParserError(t *token.Token, message string) {
	p.runtime.ErrorTokenMessage(t, message)
	p.hadParseError = true
}

// synchronize() discards tokens until we're at the beginning of the next
// statement
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == token.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF,
			token.WHILE, token.PRINT, token.RETURN:
			return
		}

		p.advance()
	}
}
