package expression

import (
	"github.com/mohitranka/golox/token"
)

type Expr interface {
	Accept(v ExprVisitor) interface{}
}

type ExprAssign struct {
	Name  token.Token
	Value Expr
}

type ExprBinary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type ExprGrouping struct {
	Expr
}

type ExprLiteral struct {
	Value interface{}
}

type ExprUnary struct {
	Operator token.Token
	Right    Expr
}

type ExprVar struct {
	Name token.Token
}

func (e *ExprAssign) Accept(v ExprVisitor) interface{} { return v.VisitAssignExpr(e) }

func (e *ExprBinary) Accept(v ExprVisitor) interface{} { return v.VisitBinaryExpr(e) }

func (e *ExprGrouping) Accept(v ExprVisitor) interface{} { return v.VisitGroupingExpr(e) }

func (e *ExprLiteral) Accept(v ExprVisitor) interface{} { return v.VisitLiteralExpr(e) }

func (e *ExprUnary) Accept(v ExprVisitor) interface{} { return v.VisitUnaryExpr(e) }

func (e *ExprVar) Accept(v ExprVisitor) interface{} { return v.VisitVarExpr(e) }
