package nillox

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/mz1290/craftinginterpreters/test/common"
)

var interpreter = ""

func init() {
	interpreter = common.GetInterpreter()
	fmt.Printf("USING: %s\n", interpreter)
}

func TestNil(t *testing.T) {
	file := "literal.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "nil\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}
