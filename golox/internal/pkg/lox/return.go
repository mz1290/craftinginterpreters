package lox

type Return struct {
	Value interface{}
}

func NewReturn(object interface{}) *Return {
	return &Return{object}
}

func IsReturnable(err error) bool {
	switch err.(type) {
	case *Return:
		return true
	default:
		return false
	}
}

func (r *Return) Error() string {
	return ""
}
