package lox

type ReturnValue struct {
	Value interface{}
}

func (rv ReturnValue) Error() string {
	return "No one cares!"
}
