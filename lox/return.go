package lox

// ReturnValue ...
type ReturnValue struct {
	ExprLiteral
}

// Error ...
func (rv ReturnValue) Error() string {
	return rv.ExprLiteral.Value.(string)
}
