package lox

import (
	"github.com/mz1290/golox/internal/pkg/ast"
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

func (p *Parser) Parse() ast.Expr {
	expr := p.expression()

	if p.hadParseError {
		return nil
	} else {
		return expr
	}
}

func (p *Parser) expression() ast.Expr {
	return p.equality()
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
