package environment

import (
	"github.com/mohitranka/golox/err"
	"github.com/mohitranka/golox/token"
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

func (e Environment) Get(name token.Token) interface{} {
	value, ok := e.Values[name.Lexeme]
	if !ok {
		panic(&err.VarError{Name: name.Lexeme, Msg: "Undefined variable '%s'"})
	}
	return value
}
