package lox

type Class struct {
	runtime *Lox
	Name    string

	// Class is responsible for storing behavior.
	Methods map[string]*Function
}

func NewClass(l *Lox, name string, methods map[string]*Function) *Class {
	return &Class{
		runtime: l,
		Name:    name,
		Methods: methods,
	}
}

func (c *Class) FindMethod(name string) *Function {
	if f, ok := c.Methods[name]; ok {
		return f
	}

	return nil
}

func (c *Class) String() string {
	return c.Name
}

func (c *Class) Call(i *Interpreter, arguments []interface{}) (interface{}, error) {
	return NewInstance(c.runtime, c), nil
}

func (c *Class) Arity() int {
	return 0
}
