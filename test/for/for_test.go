package forloop

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

func TestClassInBody(t *testing.T) {
	file := "class_in_body.lox"
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

func TestClosureInBody(t *testing.T) {
	file := "closure_in_body.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "4\n1\n4\n2\n4\n3\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestFuncInBody(t *testing.T) {
	file := "fun_in_body.lox"
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

func TestReturnClosure(t *testing.T) {
	file := "return_closure.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "i\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestReturnInside(t *testing.T) {
	file := "return_inside.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "i\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestScope(t *testing.T) {
	file := "scope.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "0\n-1\nafter\n0\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestStatementCondition(t *testing.T) {
	file := "statement_condition.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] Error at "{": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestStatementIncrement(t *testing.T) {
	file := "statement_increment.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] Error at "{": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestStatementInitializer(t *testing.T) {
	file := "statement_initializer.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] Error at "{": "expected expression"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSyntax(t *testing.T) {
	file := "syntax.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "1\n2\n3\n0\n1\n2\ndone\n0\n1\n0\n1\n2\n0\n1\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestVarInBody(t *testing.T) {
	file := "var_in_body.lox"
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
