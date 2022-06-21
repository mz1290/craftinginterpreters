package token

import "fmt"

type Token struct {
	Type    Type
	lexeme  string
	literal interface{}
	line    int
}

func New(t Type, lexeme string, literal interface{}, line int) *Token {
	return &Token{t, lexeme, literal, line}
}

func (t Token) String() string {
	return fmt.Sprintf(
		"{type: %-13s lexeme: %-15s line: %d}", t.Type, t.lexeme, t.line)
}
