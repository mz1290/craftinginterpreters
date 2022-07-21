package field

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

func TestCallFunctionField(t *testing.T) {
	file := "call_function_field.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "bar\n1\n2\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestCallNonFunctionField(t *testing.T) {
	file := "call_nonfunction_field.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 6] RuntimeError: can only call functions and classes`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestGetAndSetMethod(t *testing.T) {
	file := "get_and_set_method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "other\n1\nmethod\n2\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestGetOnBool(t *testing.T) {
	file := "get_on_bool.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have properties`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestGetOnClass(t *testing.T) {
	file := "get_on_class.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] RuntimeError: only instances have properties`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestGetOnFunction(t *testing.T) {
	file := "get_on_function.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] RuntimeError: only instances have properties`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestGetOnNil(t *testing.T) {
	file := "get_on_nil.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have properties`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestGetOnNum(t *testing.T) {
	file := "get_on_num.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have properties`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestGetOnString(t *testing.T) {
	file := "get_on_string.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have properties`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestMany(t *testing.T) {
	file := "many.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected :=
		"apple\n" +
			"apricot\n" +
			"avocado\n" +
			"banana\n" +
			"bilberry\n" +
			"blackberry\n" +
			"blackcurrant\n" +
			"blueberry\n" +
			"boysenberry\n" +
			"cantaloupe\n" +
			"cherimoya\n" +
			"cherry\n" +
			"clementine\n" +
			"cloudberry\n" +
			"coconut\n" +
			"cranberry\n" +
			"currant\n" +
			"damson\n" +
			"date\n" +
			"dragonfruit\n" +
			"durian\n" +
			"elderberry\n" +
			"feijoa\n" +
			"fig\n" +
			"gooseberry\n" +
			"grape\n" +
			"grapefruit\n" +
			"guava\n" +
			"honeydew\n" +
			"huckleberry\n" +
			"jabuticaba\n" +
			"jackfruit\n" +
			"jambul\n" +
			"jujube\n" +
			"juniper\n" +
			"kiwifruit\n" +
			"kumquat\n" +
			"lemon\n" +
			"lime\n" +
			"longan\n" +
			"loquat\n" +
			"lychee\n" +
			"mandarine\n" +
			"mango\n" +
			"marionberry\n" +
			"melon\n" +
			"miracle\n" +
			"mulberry\n" +
			"nance\n" +
			"nectarine\n" +
			"olive\n" +
			"orange\n" +
			"papaya\n" +
			"passionfruit\n" +
			"peach\n" +
			"pear\n" +
			"persimmon\n" +
			"physalis\n" +
			"pineapple\n" +
			"plantain\n" +
			"plum\n" +
			"plumcot\n" +
			"pomegranate\n" +
			"pomelo\n" +
			"quince\n" +
			"raisin\n" +
			"rambutan\n" +
			"raspberry\n" +
			"redcurrant\n" +
			"salak\n" +
			"salmonberry\n" +
			"satsuma\n" +
			"strawberry\n" +
			"tamarillo\n" +
			"tamarind\n" +
			"tangerine\n" +
			"tomato\n" +
			"watermelon\n" +
			"yuzu\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestMethodBindsThis(t *testing.T) {
	file := "method_binds_this.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "foo1\n1\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestMethod(t *testing.T) {
	file := "method.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "got method\narg\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestOnInstance(t *testing.T) {
	file := "on_instance.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "bar value\nbaz value\nbar value\nbaz value\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestSetEvaluationOrder(t *testing.T) {
	file := "set_evaluation_order.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: undefined variable "undefined1"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetOnBool(t *testing.T) {
	file := "set_on_bool.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have fields`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetOnClass(t *testing.T) {
	file := "set_on_class.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 2] RuntimeError: only instances have fields`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetOnFunction(t *testing.T) {
	file := "set_on_function.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 3] RuntimeError: only instances have fields`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetOnNil(t *testing.T) {
	file := "set_on_nil.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have fields`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetOnNum(t *testing.T) {
	file := "set_on_num.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have fields`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestSetOnString(t *testing.T) {
	file := "set_on_string.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 1] RuntimeError: only instances have fields`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}

func TestUndefined(t *testing.T) {
	file := "undefined.lox"
	cmd := exec.Command(interpreter, file)
	stderr, _ := cmd.StderrPipe()
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}

	expected := `[line 4] RuntimeError: undefined variable property "bar"`

	scanner := bufio.NewScanner(stderr)
	scanner.Scan()
	actualErr := scanner.Text()
	if actualErr != expected {
		t.Fatalf("expected error (%s) got %s", expected, actualErr)
	}
}
