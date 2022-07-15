package inheritance

import (
	"bufio"
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

func TestConstructor(t *testing.T) {
	file := "constructor.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "value\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestInheritFromFunction(t *testing.T) {
	file := "inherit_from_function.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] RuntimeError: superclass must be a class`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestInheritFromNil(t *testing.T) {
	file := "inherit_from_nil.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] RuntimeError: superclass must be a class`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestInheritFromNumber(t *testing.T) {
	file := "inherit_from_number.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] RuntimeError: superclass must be a class`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestInheritMethods(t *testing.T) {
	file := "inherit_methods.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "foo\nbar\nbar\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestParenthesizedSuperclas(t *testing.T) {
	file := "parenthesized_superclass.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 4] Error at "(": "expected superclass name"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetFieldsFromBaseClass(t *testing.T) {
	file := "set_fields_from_base_class.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "foo 1\nfoo 2\nbar 1\nbar 2\nbar 1\nbar 2\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}
