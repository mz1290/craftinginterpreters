package expressions

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
	os.Setenv("DEBUGLOX", "")

	interpreter = common.GetInterpreter()
	fmt.Printf("USING: %s\n", interpreter)
}
func TestExpressions(t *testing.T) {
	file := "expressions.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "2\n"

	if string(stdout) != expected {
		t.Fatalf("expected %x got %x", expected, stdout)
	}
}
