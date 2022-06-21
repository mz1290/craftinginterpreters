package scanner

import (
	"fmt"
	"strconv"

	"github.com/mz1290/golox/internal/pkg/common"
	"github.com/mz1290/golox/internal/pkg/token"
)

type Scanner struct {
	source  string
	tokens  []*token.Token
	start   int
	current int
	line    int
	Errors  []*ScannerErr
}

type ScannerErr struct {
	Line    int
	Where   string
	Message string
}

func (se ScannerErr) Error() string {
	return fmt.Sprintf("[line %d] Error%s: %q", se.Line, se.Where, se.Message)
}

func New(source string) *Scanner {
	return &Scanner{
		source:  source,
		start:   0,
		current: 0,
		line:    1,
	}
}

func (s *Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		// We are at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.New(token.EOF, "", nil, s.line))
	return s.tokens
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(token.LEFT_PAREN, nil)
	case ')':
		s.addToken(token.RIGHT_PAREN, nil)
	case '{':
		s.addToken(token.LEFT_BRACE, nil)
	case '}':
		s.addToken(token.RIGHT_BRACE, nil)
	case ',':
		s.addToken(token.COMMA, nil)
	case '.':
		s.addToken(token.DOT, nil)
	case '-':
		s.addToken(token.MINUS, nil)
	case '+':
		s.addToken(token.PLUS, nil)
	case ';':
		s.addToken(token.SEMICOLON, nil)
	case '*':
		s.addToken(token.STAR, nil)
	case '!':
		var t token.Type
		if s.match('=') { // Maybe write a ternary-like token type func to avoid this
			t = token.BANG_EQUAL
		} else {
			t = token.BANG
		}
		s.addToken(t, nil)
	case '=':
		var t token.Type
		if s.match('=') {
			t = token.EQUAL_EQUAL
		} else {
			t = token.EQUAL
		}
		s.addToken(t, nil)
	case '<':
		var t token.Type
		if s.match('=') {
			t = token.LESS_EQUAL
		} else {
			t = token.LESS
		}
		s.addToken(t, nil)
	case '>':
		var t token.Type
		if s.match('=') {
			t = token.GREATER_EQUAL
		} else {
			t = token.GREATER
		}
		s.addToken(t, nil)
	case '/':
		if s.match('/') {
			// A copmment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH, nil)
		}
	case ' ', '\r', '\t':
		// ignore
	case '\n':
		s.line++
	case '"':
		s.string()
	default:
		if common.IsDigit(c) {
			s.number()
		} else if common.IsAlpha(c) {
			s.identifier()
		} else {
			s.Errors = append(s.Errors, &ScannerErr{s.line, " Scanner",
				fmt.Sprintf("Unexpected character %q", c)})
		}
	}
}

func (s *Scanner) identifier() {
	for common.IsAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	ty, ok := token.Keywords[text]

	// If text was not a reserver word, it is an identifier
	if !ok {
		ty = token.IDENTIFIER
	}

	s.addToken(ty, nil)
}

func (s *Scanner) number() {
	for common.IsDigit(s.peek()) {
		s.advance()
	}

	// Look for fractional part
	if s.peek() == '.' && common.IsDigit(s.peekNext()) {
		// Consumer the '.'
		s.advance()

		for common.IsDigit(s.peek()) {
			s.advance()
		}
	}

	// Convert string to Go's float64
	num, err := strconv.ParseFloat(s.source[s.start:s.current], 64)
	if err != nil {
		s.Errors = append(s.Errors, &ScannerErr{s.line, " Scanner",
			fmt.Sprintf("Failed to convert %q to Lox number",
				s.source[s.start:s.current])})
		return
	}

	s.addToken(token.NUMBER, num)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.Errors = append(s.Errors,
			&ScannerErr{s.line, " Scanner", "Unterminated string."})
		return
	}

	// The closing "
	s.advance()

	// Trim the surrounding quotes
	value := s.source[s.start+1 : s.current-1]
	s.addToken(token.STRING, value)
}

// lookahead, we do not consume and advance
func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}

	return s.source[s.current]
}

func (s Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}

	s.current++
	return true
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	ch := s.source[s.current]
	s.current++
	return ch
}

func (s *Scanner) addToken(t token.Type, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, token.New(t, text, literal, s.line))
}
