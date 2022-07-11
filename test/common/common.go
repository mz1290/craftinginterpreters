package common

import (
	"bytes"
	"log"
	"os"
	"strconv"
	"strings"
)

// interpreter=golox go test -v -count=1 ./...
var interpreter = GetInterpreter()

func GetInterpreter() string {
	chk := os.Getenv("interpreter")

	if chk != "golox" && chk != "clox" {
		log.Println("ERROR: interpreter not specified")
		chk = "golox"
		//chk = "clox"
	}

	if chk == "golox" {
		return "../../golox/golox"
	} else {
		return "../../clox/build/clox"
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
	var token TokenInfo

	tokenLine := strings.Fields(string(output))

	if strings.HasSuffix(interpreter, "clox") {
		val, _ := strconv.Atoi(tokenLine[1])
		token.Type = TokenType(val).String()
	} else {
		token.Type = tokenLine[1]
	}

	if token.Type == "EOF" {
		token.Lexeme = ""
		token.Line = tokenLine[4]
	} else {
		token.Lexeme = tokenLine[3]
		token.Line = tokenLine[5]
	}

	return token
}
