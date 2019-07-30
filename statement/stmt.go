package statement

import (
	"github.com/mohitranka/golox/expression"
)

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type ExpressionStmt struct {
	Expression expression.Expr
}

func NewExpressionStmt(expression expression.Expr) Stmt {
	return &ExpressionStmt{Expression: expression}
}

func (expr *ExpressionStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitExpressionStmt(expr)
}

type PrintStmt struct {
	Expression expression.Expr
}

func NewPrintStmt(expression expression.Expr) Stmt {
	return &PrintStmt{Expression: expression}
}

func (expr *PrintStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(expr)
}
