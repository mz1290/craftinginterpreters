package comments

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

func TestLineAtEOF(t *testing.T) {
	file := "line_at_eof.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "ok\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}

func TestOnlyLineCommentAndLine(t *testing.T) {
	file := "only_line_comment_and_line.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := ""

	if string(stdout) != expected {
		t.Fatalf("expected %q got %s", expected, string(stdout))
	}
}

func TestOnlyLineComment(t *testing.T) {
	file := "only_line_comment.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := ""

	if string(stdout) != expected {
		t.Fatalf("expected %q got %s", expected, string(stdout))
	}
}

func TestUnicode(t *testing.T) {
	file := "unicode.lox"
	cmd := exec.Command(interpreter, file)
	stdout, _ := cmd.Output()
	expected := "ok\n"

	if string(stdout) != expected {
		t.Fatalf("expected %s got %s", expected, string(stdout))
	}
}