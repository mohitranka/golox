package environment

import (
	"fmt"
	"github.com/mohitranka/golox/err"
)

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	ne := new(Environment)
	ne.Values = make(map[string]interface{})
	ne.Enclosing = enclosing
	return ne
}

func (e Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e Environment) Get(name string) interface{} {
	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}
	value, ok := e.Values[name]
	if !ok {
		fmt.Println(&err.VarError{Name: name, Msg: "Undefined variable"})
		return nil
	}
	return value
}

func (e Environment) Assign(name string, value interface{}) {
	if e.Enclosing != nil {
		e.Enclosing.Assign(name, value)
		return
	}
	_, ok := e.Values[name]
	if !ok {
		fmt.Println(&err.VarError{Name: name, Msg: "Undefined variable"})
		return
	}
	e.Values[name] = value

}
