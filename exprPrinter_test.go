package lox

import (
	"testing"
)

func TestExpressionPrinter(t *testing.T) {
	expr := ExprBinary{
		Left: &ExprUnary{
			Operator: lox.Token{lox.MINUS, "-", nil, 1},
			Right:    &ExprLiteral{123},
		},
		Operator: lox.Token{lox.STAR, "*", nil, 1},
		Right: &ExprGrouping{
			Expr: &ExprLiteral{45.67},
		},
	}
	expected := "(* (- 123) (group 45.67))"
	got := PrintExpr(&expr)

	if got != expected {
		t.Errorf("Evaluation was incorrect, got: %s, expected: %s", got, expected)
	}
}
