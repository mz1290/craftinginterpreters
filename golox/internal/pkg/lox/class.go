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
	instance := NewInstance(c.runtime, c)

	initializer := c.FindMethod("init")
	if initializer != nil {
		initializer.Bind(instance).Call(i, arguments)
	}

	return instance, nil
}

func (c *Class) Arity() int {
	initializer := c.FindMethod("init")
	if initializer == nil {
		return 0
	}

	return initializer.Arity()
}
