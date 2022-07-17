package precedence

import (
	//"bufio"
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

func TestPrecedence(t *testing.T) {
	file := "precedence.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected :=
		`14
8
4
0
true
true
true
true
0
0
0
0
4
`

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}
