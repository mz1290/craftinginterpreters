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
		//interpreter = "golox"
		chk = "clox"
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
	tokenLine := strings.Fields(string(output))

	if strings.HasSuffix(interpreter, "golox") {
		return goloxTokenInfo(tokenLine)
	} else {
		return cloxTokenInfo(tokenLine)
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
		//lineRaw := tokenLine[5]
		//lineClean := lineRaw[:len(lineRaw)-1]
		//token.Line = lineClean

		token.Line = tokenLine[5][:len(tokenLine[5])-1]
	}

	return token
}

func cloxTokenInfo(tokenLine []string) TokenInfo {
	var token TokenInfo

	token.Line = tokenLine[0]

	val, _ := strconv.Atoi(tokenLine[1])
	token.Type = TokenType(val).String()

	token.Lexeme = tokenLine[2][1 : len(tokenLine[2])-1]

	return token
}
