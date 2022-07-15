package logicaloperator

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

func TestAndTruth(t *testing.T) {
	file := "and_truth.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "false\nnil\nok\nok\nok\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestAnd(t *testing.T) {
	file := "and.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "false\n1\nfalse\ntrue\n3\ntrue\nfalse\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestOrTruth(t *testing.T) {
	file := "or_truth.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "ok\nok\ntrue\n0\ns\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestOr(t *testing.T) {
	file := "or.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "1\n1\ntrue\nfalse\nfalse\nfalse\ntrue\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}
