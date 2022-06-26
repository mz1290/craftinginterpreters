package main

import (
	"fmt"
	"os"

	"github.com/mz1290/golox/internal/pkg/lox"
)

func main() {
	nArgs := len(os.Args)

	if nArgs > 2 {
		fmt.Println("Usage: golox [script]")
		os.Exit(64) // exit code standard https://www.freebsd.org/cgi/man.cgi?query=sysexits&apropos=0&sektion=0&manpath=FreeBSD+4.3-RELEASE&format=html
	} else {
		l := lox.New()

		if nArgs == 2 {
			l.RunFile(os.Args[1])
		} else {
			l.RunPrompt()
		}
	}
}
