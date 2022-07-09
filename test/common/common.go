package common

import (
	"bytes"
	"os"
	"strings"
)

var interpreter = os.Getenv("interpreter")

func init() {
	if interpreter != "golox" && interpreter != "clox" {
		//interpreter = "golox"
		interpreter = "clox"
	}

	if interpreter == "golox" {
		interpreter = "../../golox/golox"
	} else {
		interpreter = "../../clox/build/clox"
	}
}

func GetStdOutLines(output []byte) [][]byte {
	return bytes.Split(output, []byte{'\n'})
}

type TokenInfo struct {
	Type   string
	Lexeme string
	Line   string
}

func GetTokenInfo(output []byte) TokenInfo {
	tokenLine := strings.Fields(string(output[:len(output)-1]))

	if strings.HasSuffix(interpreter, "golox") {
		return goloxTokenInfo(tokenLine)
	} else {
		return TokenInfo{}
	}
}

func goloxTokenInfo(tokenLine []string) TokenInfo {
	var token TokenInfo

	token.Type = tokenLine[1]
	if token.Type == "EOF" {
		token.Lexeme = ""
		token.Line = tokenLine[4]
	} else {
		token.Lexeme = tokenLine[3]
		token.Line = tokenLine[5]
	}

	return token
}
