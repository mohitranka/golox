package lox

// Callable ...
type Callable interface {
	Call(i *Interpreter, args []interface{}) interface{}
	Arity() int
}
