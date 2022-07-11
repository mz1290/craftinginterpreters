package bool

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

func TestEquality(t *testing.T) {
	file := "equality.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "true\nfalse\nfalse\ntrue\n" +
		"false\nfalse\nfalse\nfalse\nfalse\n" +
		"false\ntrue\ntrue\nfalse\n" +
		"true\ntrue\ntrue\ntrue\ntrue\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestNot(t *testing.T) {
	file := "not.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "false\ntrue\ntrue\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}
