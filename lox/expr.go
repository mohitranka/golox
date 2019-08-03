package lox

// Expr ...
type Expr interface {
	Accept(v ExprVisitor) interface{}
}

// ExprAssign ...
type ExprAssign struct {
	Name  Token
	Value Expr
}

// ExprBinary ...
type ExprBinary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// ExprGrouping ...
type ExprGrouping struct {
	Expr
}

// ExprLiteral ...
type ExprLiteral struct {
	Value interface{}
}

// ExprUnary ...
type ExprUnary struct {
	Operator Token
	Right    Expr
}

// ExprVar ...
type ExprVar struct {
	Name Token
}

// ExprLogical ...
type ExprLogical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

// ExprCall ..
type ExprCall struct {
	Callee    Expr
	Paren     Token
	Arguments []*Expr
}

// Accept ...
func (e *ExprAssign) Accept(v ExprVisitor) interface{} { return v.VisitAssignExpr(e) }

// Accept ...
func (e *ExprBinary) Accept(v ExprVisitor) interface{} { return v.VisitBinaryExpr(e) }

// Accept ...
func (e *ExprGrouping) Accept(v ExprVisitor) interface{} { return v.VisitGroupingExpr(e) }

// Accept ...
func (e *ExprLiteral) Accept(v ExprVisitor) interface{} { return v.VisitLiteralExpr(e) }

// Accept ...
func (e *ExprUnary) Accept(v ExprVisitor) interface{} { return v.VisitUnaryExpr(e) }

// Accept ...
func (e *ExprVar) Accept(v ExprVisitor) interface{} { return v.VisitVarExpr(e) }

// Accept ...
func (e *ExprLogical) Accept(v ExprVisitor) interface{} { return v.VisitLogicalExpr(e) }

// Accept ...
func (e *ExprCall) Accept(v ExprVisitor) interface{} { return v.VisitCallExpr(e) }
