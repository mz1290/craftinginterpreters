package super

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

func TestBoundMethod(t *testing.T) {
	file := "bound_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "A.method(arg)\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestCallOtherMethod(t *testing.T) {
	file := "call_other_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "Derived.bar()\nBase.foo()\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestCallSameMethod(t *testing.T) {
	file := "call_same_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "Derived.foo()\nBase.foo()\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestClosure(t *testing.T) {
	file := "closure.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "Base\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestConstructor(t *testing.T) {
	file := "constructor.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "Derived.init()\nBase.init(a, b)\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestExtraArguments(t *testing.T) {
	file := "extra_arguments.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: expected 2 arguments but got 4`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestIndirectlyInherited(t *testing.T) {
	file := "indirectly_inherited.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "C.foo()\nA.foo()\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestMissingArguments(t *testing.T) {
	file := "missing_arguments.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: expected 2 arguments but got 1`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestNoSuperclassBind(t *testing.T) {
	file := "no_superclass_bind.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] Error at "super": "can't use 'super' in a class with no superclass"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestNoSuperclassCall(t *testing.T) {
	file := "no_superclass_call.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] Error at "super": "can't use 'super' in a class with no superclass"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestNoSuperclassMethod(t *testing.T) {
	file := "no_superclass_method.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 5] undefined property doesNotExist`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestParenthesized(t *testing.T) {
	file := "parenthesized.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 8] Error at ")": "expected '.' after 'super'"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestReassignSuperclass(t *testing.T) {
	file := "reassign_superclass.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "Base.method()\nBase.method()\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestSuperAtTopLevel(t *testing.T) {
	file := "super_at_top_level.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] Error at "super": "can't use 'super' outside of a class"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}

	expected = `[line 2] Error at "super": "can't use 'super' outside of a class"`

	scanner.Scan()
	actualErr = scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSuperInClosureInInheritedMethod(t *testing.T) {
	file := "super_in_closure_in_inherited_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "A\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestSuperInInheritedMethod(t *testing.T) {
	file := "super_in_inherited_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "A\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestSuperInTopLevelFunction(t *testing.T) {
	file := "super_in_top_level_function.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] Error at "super": "can't use 'super' outside of a class"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSuperWithoutDot(t *testing.T) {
	file := "super_without_dot.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 6] Error at ";": "expected '.' after 'super'"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSuperWithouName(t *testing.T) {
	file := "super_without_name.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 5] Error at ";": "expected superclass method name"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestThisInSuperclassMethod(t *testing.T) {
	file := "this_in_superclass_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "a\nb\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}
