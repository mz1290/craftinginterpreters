package iflox

import (
	"bufio"
	"fmt"
	"github.com/mz1290/craftinginterpreters/test/common"
	"os/exec"
	"testing"
)

var interpreter = ""

func init() {
	interpreter = common.GetInterpreter()
	fmt.Printf("USING: %s\n", interpreter)
}

func TestClassInElse(t *testing.T) {
	file := "class_in_else.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "class": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestClassInThen(t *testing.T) {
	file := "class_in_then.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "class": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestDanglingElse(t *testing.T) {
	file := "dangling_else.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "good\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestElse(t *testing.T) {
	file := "else.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "good\ngood\nblock\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestFunInElse(t *testing.T) {
	file := "fun_in_else.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "fun": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestFunInThen(t *testing.T) {
	file := "fun_in_then.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "fun": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestIf(t *testing.T) {
	file := "if.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "good\nblock\ntrue\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestTruth(t *testing.T) {
	file := "truth.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "false\nnil\ntrue\n0\nempty\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestVarInElse(t *testing.T) {
	file := "var_in_else.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "var": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestVarInThen(t *testing.T) {
	file := "var_in_then.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "var": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}
