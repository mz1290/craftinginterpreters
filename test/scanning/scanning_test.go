package scanning

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/mz1290/craftinginterpreters/test/common"
)

var interpreter = ""

func init() {
	// Tell lox to spit out scanning debug info
	os.Setenv("DEBUGLOX", "scanning")

	interpreter = common.GetInterpreter()
	fmt.Printf("USING: %s\n", interpreter)
}

func TestIdentifiers(t *testing.T) {
	file := "identifiers.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()

	expected := []common.TokenInfo{
		{"IDENTIFIER", "andy", "1"},
		{"IDENTIFIER", "formless", "1"},
		{"IDENTIFIER", "fo", "1"},
		{"IDENTIFIER", "_", "1"},
		{"IDENTIFIER", "_123", "1"},
		{"IDENTIFIER", "_abc", "1"},
		{"IDENTIFIER", "ab123", "1"},
		{"IDENTIFIER", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890_", "2"},
		{"EOF", "", "2"},
	}

	results := common.GetStdOutLines(stdout)
	for i, result := range results {
		token := common.GetTokenInfo(result)

		if token.Type == "EOF" {
			break
		}

		if token == expected[i] {
			continue
		}
		t.Fatalf("expected %v got %v", expected[i], token)
	}
}

func TestKeywords(t *testing.T) {
	file := "keywords.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()

	expected := []common.TokenInfo{
		{"AND", "and", "1"},
		{"CLASS", "class", "1"},
		{"ELSE", "else", "1"},
		{"FALSE", "false", "1"},
		{"FOR", "for", "1"},
		{"FUN", "fun", "1"},
		{"IF", "if", "1"},
		{"NIL", "nil", "1"},
		{"OR", "or", "1"},
		{"RETURN", "return", "1"},
		{"SUPER", "super", "1"},
		{"THIS", "this", "1"},
		{"TRUE", "true", "1"},
		{"VAR", "var", "1"},
		{"WHILE", "while", "1"},
		{"EOF", "", "1"},
	}

	results := common.GetStdOutLines(stdout)
	for i, result := range results {
		token := common.GetTokenInfo(result)

		if token.Type == "EOF" {
			break
		}

		if token == expected[i] {
			continue
		}
		t.Fatalf("expected %v got %v", expected[i], token)
	}
}

func TestNumbers(t *testing.T) {
	file := "numbers.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()

	expected := []common.TokenInfo{
		{"NUMBER", "123", "1"},
		{"NUMBER", "123.456", "2"},
		{"DOT", ".", "3"},
		{"NUMBER", "456", "3"},
		{"NUMBER", "123", "4"},
		{"DOT", ".", "4"},
		{"EOF", "", "4"},
	}

	results := common.GetStdOutLines(stdout)
	for i, result := range results {
		token := common.GetTokenInfo(result)

		if token.Type == "EOF" {
			break
		}

		if token == expected[i] {
			continue
		}
		t.Fatalf("expected %v got %v", expected[i], token)
	}
}

func TestPunctuators(t *testing.T) {
	file := "punctuators.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()

	expected := []common.TokenInfo{
		{"LEFT_PAREN", "(", "1"},
		{"RIGHT_PAREN", ")", "1"},
		{"LEFT_BRACE", "{", "1"},
		{"RIGHT_BRACE", "}", "1"},
		{"SEMICOLON", ";", "1"},
		{"COMMA", ",", "1"},
		{"PLUS", "+", "1"},
		{"MINUS", "-", "1"},
		{"STAR", "*", "1"},
		{"BANG_EQUAL", "!=", "1"},
		{"EQUAL_EQUAL", "==", "1"},
		{"LESS_EQUAL", "<=", "1"},
		{"GREATER_EQUAL", ">=", "1"},
		{"BANG_EQUAL", "!=", "1"},
		{"LESS", "<", "1"},
		{"GREATER", ">", "1"},
		{"SLASH", "/", "1"},
		{"DOT", ".", "1"},
		{"EOF", "", "1"},
	}

	results := common.GetStdOutLines(stdout)
	for i, result := range results {
		token := common.GetTokenInfo(result)

		if token.Type == "EOF" {
			break
		}

		if token == expected[i] {
			continue
		}
		t.Fatalf("expected %v got %v", expected[i], token)
	}
}

func TestStrings(t *testing.T) {
	file := "strings.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()

	expected := []common.TokenInfo{
		{"STRING", `""`, "1"},
		{"STRING", `"string"`, "2"},
		{"EOF", "", "2"},
	}

	results := common.GetStdOutLines(stdout)
	for i, result := range results {
		token := common.GetTokenInfo(result)

		if token.Type == "EOF" {
			break
		}

		if token == expected[i] {
			continue
		}
		t.Fatalf("expected %v got %v", expected[i], token)
	}
}

func TestWhitespace(t *testing.T) {
	file := "whitespace.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()

	expected := []common.TokenInfo{
		{"IDENTIFIER", "space", "2"},
		{"IDENTIFIER", "tabs", "2"},
		{"IDENTIFIER", "newlines", "2"},
		{"IDENTIFIER", "end", "7"},
		{"EOF", "", "8"},
	}

	results := common.GetStdOutLines(stdout)
	for i, result := range results {
		token := common.GetTokenInfo(result)

		if token.Type == "EOF" {
			break
		}

		if token == expected[i] {
			continue
		}
		t.Fatalf("expected %v got %v", expected[i], token)
	}
}
