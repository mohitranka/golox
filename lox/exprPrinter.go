package lox

import (
	"bytes"
	"fmt"
)

// PrintExpr ...
func PrintExpr(e Expr) string {
	switch t := e.(type) {
	case *ExprBinary:
		return parenthesize(t.Operator.Lexeme, t.Left, t.Right)
	case *ExprLiteral:
		if t.Value == nil {
			return "nil"
		}
		return fmt.Sprint(t.Value)
	case *ExprGrouping:
		return parenthesize("group", t.Expr)
	case *ExprUnary:
		return parenthesize(t.Operator.Lexeme, t.Right)
	default:
		return ""
	}
}

func parenthesize(name string, exprs ...Expr) string {
	buf := bytes.Buffer{}
	buf.WriteRune('(')
	buf.WriteString(name)
	for _, expr := range exprs {
		buf.WriteRune(' ')
		buf.WriteString(PrintExpr(expr))
	}
	buf.WriteRune(')')
	return buf.String()
}
