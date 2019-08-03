package lox

// ExprVisitor ...
type ExprVisitor interface {
	VisitAssignExpr(ea *ExprAssign) interface{}
	VisitBinaryExpr(eb *ExprBinary) interface{}
	VisitGroupingExpr(eg *ExprGrouping) interface{}
	VisitLiteralExpr(el *ExprLiteral) interface{}
	VisitUnaryExpr(eu *ExprUnary) interface{}
	VisitVarExpr(ev *ExprVar) interface{}
	VisitLogicalExpr(eb *ExprLogical) interface{}
	VisitCallExpr(ec *ExprCall) interface{}
}
