package expression

import (
	"github.com/mohitranka/golox/token"
	"testing"
)

func TestExpressionPrinter(t *testing.T) {
	expr := ExprBinary{
		Left: &ExprUnary{
			Operator: token.Token{token.MINUS, "-", nil, 1},
			Right:    &ExprLiteral{123},
		},
		Operator: token.Token{token.STAR, "*", nil, 1},
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
