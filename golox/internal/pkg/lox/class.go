package lox

import "fmt"

type Class struct {
	Name string
}

func NewClass(name string) *Class {
	return &Class{
		Name: name,
	}
}

func (c *Class) String() string {
	return fmt.Sprintf("%s", c.Name)
}
