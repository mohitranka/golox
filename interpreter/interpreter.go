package interpreter

import (
	"fmt"
	"github.com/mohitranka/golox/expression"
	"github.com/mohitranka/golox/token"
	"reflect"
)

type Interpreter struct {
}

func (i Interpreter) VisitLiteralExpr(expr expression.ExprLiteral) interface{} {
	return expr.Value
}

func (i Interpreter) VisitGroupingExpr(expr expression.ExprGrouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i Interpreter) Interpret(expr expression.Expr) {
	obj := i.evaluate(expr)
	str := fmt.Sprintf("%v", obj)
	fmt.Println(str)
}

func (i Interpreter) evaluate(expr expression.Expr) interface{} {
	return expression.PrintExpr(expr) //FIXME: This returns string, requires work
}

func (i Interpreter) VisitUnaryExpr(expr expression.ExprUnary) interface{} {
	right := i.evaluate(&expr)
	switch expr.Operator.Type {
	case token.BANG:
		return !i.isTruthy(right.(float64))
	case token.MINUS:
		return -right.(float64)
	}
	return nil
}

func (i Interpreter) isTruthy(obj float64) bool {
	if obj == float64(0) {
		return false
	}

	return true
}

func (i Interpreter) VisitBinaryExpr(expr expression.ExprBinary) interface{} {
	left := i.evaluate(expr.Left)
	right := i.evaluate(expr.Right)

	switch expr.Operator.Type {
	case token.GREATER:
		return left.(float64) > right.(float64)
	case token.GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case token.LESS:
		return left.(float64) < right.(float64)
	case token.LESS_EQUAL:
		return left.(float64) <= right.(float64)
	case token.BANG_EQUAL:
		return !(left.(float64) == right.(float64))
	case token.EQUAL_EQUAL:
		return left.(float64) == right.(float64)
	case token.MINUS:
		return left.(float64) - right.(float64)
	case token.PLUS:
		typeLeft := reflect.TypeOf(left).String()
		typeRight := reflect.TypeOf(right).String()
		if (typeLeft == "float64" || typeLeft == "float32" || typeLeft == "int") && (typeRight == "float64" || typeRight == "float32" || typeRight == "int") {
			return left.(float64) + right.(float64)
		} else if typeLeft == "string" && typeRight == "string" {
			return left.(string) + right.(string)
		}
	case token.SLASH:
		return left.(float64) / right.(float64)
	case token.STAR:
		return left.(float64) * right.(float64)
	}
	return nil
}
