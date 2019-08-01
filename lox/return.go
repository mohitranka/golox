package lox

type ReturnValue struct {
	ExprLiteral
}

func (rv ReturnValue) Error() string {
	return rv.ExprLiteral.Value.(string)
}
