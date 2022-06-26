package errors

import (
	"fmt"

	"github.com/mz1290/golox/internal/pkg/token"
)

var RunetimeError = NewCustomErr("RuntimeError")

type CustomErr struct {
	Token   *token.Token
	Message string
	Type    string
}

func (e *CustomErr) Error() string {
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

func (e *CustomErr) New(t *token.Token, message string) error {
	return &CustomErr{
		Token:   t,
		Message: message,
		Type:    e.Type,
	}
}

func NewCustomErr(ty string) *CustomErr {
	return &CustomErr{Type: ty}
}
