package interpreter

import (
	"fmt"
	"github.com/mohitranka/golox/expression"
	"github.com/mohitranka/golox/statement"
	"github.com/mohitranka/golox/token"
	"reflect"
)

type Interpreter struct {
}

func (i Interpreter) VisitLiteralExpr(expr *expression.ExprLiteral) interface{} {
	return expr.Value
}

func (i Interpreter) VisitGroupingExpr(expr *expression.ExprGrouping) interface{} {
	return i.evaluate(expr.Expr)
}

func (i Interpreter) VisitVarExpr(expr *expression.ExprVar) interface{} {
	return nil //TODO: To be implemented. Placeholder implementation so that Interpreter is an interface of ExprVisitor
}

func (i Interpreter) Interpret(statements []statement.Stmt) {
	for _, statement := range statements {
		i.execute(statement)
	}
}

func (i Interpreter) execute(stmt statement.Stmt) interface{} {
	return stmt.Accept(i)
}

func (i Interpreter) evaluate(expr expression.Expr) interface{} {
	return expr.Accept(i)
}

func (i Interpreter) VisitAssignExpr(expr *expression.ExprAssign) interface{} {
	return nil //TODO: To be implemented. Placeholder implementation so that Interpreter is an interface of ExprVisitor
}

func (i Interpreter) VisitUnaryExpr(expr *expression.ExprUnary) interface{} {
	right := i.evaluate(expr.Right)
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

func (i Interpreter) VisitBinaryExpr(expr *expression.ExprBinary) interface{} {
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

func (i Interpreter) VisitExpressionStmt(stmt *statement.ExpressionStmt) interface{} {
	return i.evaluate(stmt.Expression)
}

func (i Interpreter) VisitPrintStmt(stmt *statement.PrintStmt) interface{} {
	value := i.evaluate(stmt.Expression)
	fmt.Printf("%v\n", value)
	return nil
}
