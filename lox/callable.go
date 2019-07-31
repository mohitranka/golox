package lox

type Callable interface {
	Call(i *Interpreter, args ...interface{}) interface{}
}
