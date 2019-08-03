package lox

import (
	"fmt"
)

// Environment ...
type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

// NewEnvironment ...
func NewEnvironment(enclosing *Environment) *Environment {
	ne := new(Environment)
	ne.Values = make(map[string]interface{})
	ne.Enclosing = enclosing
	return ne
}

// Define ...
func (e Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

// Get ...
func (e Environment) Get(name string) interface{} {
	if value, ok := e.Values[name]; ok {
		return value
	}
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}
	fmt.Println(&VarError{Name: name, Msg: "Undefined variable"})
	return nil
}

// Assign ...
func (e Environment) Assign(name string, value interface{}) {
	if _, ok := e.Values[name]; ok {
		e.Values[name] = value
		return
	}
	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}

	fmt.Println(&VarError{Name: name, Msg: "Undefined variable"})
}
