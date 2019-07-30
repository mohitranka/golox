package environment

import (
	"github.com/mohitranka/golox/err"
)

type Environment struct {
	Values map[string]interface{}
}

func NewEnvironment() *Environment {
	ne := new(Environment)
	ne.Values = make(map[string]interface{})
	return ne
}

func (e Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e Environment) Get(name string) interface{} {
	value, ok := e.Values[name]
	if !ok {
		panic(&err.VarError{Name: name, Msg: "Undefined variable '%s'"})
	}
	return value
}

func (e Environment) Assign(name string, value interface{}) {
	_, ok := e.Values[name]
	if ok {
		e.Values[name] = value
	}
	panic(&err.VarError{Name: name, Msg: "Undefined variable '%s'"})
}
