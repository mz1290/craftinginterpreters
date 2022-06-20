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
		var err error
		l := lox.New()

		if nArgs == 2 {
			err = l.RunFile(os.Args[1])
		} else {
			err = l.RunPrompt()
		}

		// Eventually some function here that does a type check on error and
		// returns proper status code https://yourbasic.org/golang/create-error/
		if err != nil {
			fmt.Println(err)
			os.Exit(65)
		}
	}

	os.Exit(0)
}
