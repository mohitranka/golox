package environment

import (
	"fmt"
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
		fmt.Println(&err.VarError{Name: name, Msg: "Undefined variable"})
		return nil
	}
	return value
}

func (e Environment) Assign(name string, value interface{}) {
	_, ok := e.Values[name]
	if !ok {
		fmt.Println(&err.VarError{Name: name, Msg: "Undefined variable"})
		return
	}
	e.Values[name] = value

}
