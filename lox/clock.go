package lox

import "time"

// Clock ...
type Clock struct{}

// Arity ...
func (c Clock) Arity() int {
	return 0
}

// Call ...
func (c Clock) Call(i *Interpreter, args []interface{}) interface{} {
	return time.Now().UnixNano()
}
