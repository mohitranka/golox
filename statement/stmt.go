package statement

import (
	"github.com/mohitranka/golox/expression"
	"github.com/mohitranka/golox/token"
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

type VarStmt struct {
	Name        token.Token
	Initializer expression.Expr
}

func NewVarStmt(name token.Token, initializer expression.Expr) Stmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

func (expr *VarStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(expr)
}
