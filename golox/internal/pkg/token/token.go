package token

type Token struct {
	Type    Type
	lexeme  string
	literal interface{}
	line    int
}

func New(t Type, lexeme string, literal interface{}, line int) *Token {
	return &Token{t, lexeme, literal, line}
}
