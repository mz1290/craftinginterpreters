package lox

type Class struct {
	runtime    *Lox
	Name       string
	superclass *Class

	// Class is responsible for storing behavior.
	Methods map[string]*Function
}

func NewClass(l *Lox, name string, superclass *Class, methods map[string]*Function) *Class {
	return &Class{
		runtime:    l,
		superclass: superclass,
		Name:       name,
		Methods:    methods,
	}
}

func (c *Class) FindMethod(name string) *Function {
	if f, ok := c.Methods[name]; ok {
		return f
	}

	// If we don't find the method for the specific instance's class, we will
	// recurse up through the superclass chain and check therea.
	if c.superclass != nil {
		return c.superclass.FindMethod(name)
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

func IsClass(object interface{}) bool {
	_, ok := object.(*Class)
	return ok
}
