package expression

type ExprVisitor interface {
	VisitAssignExpr(ea *ExprAssign) interface{}
	VisitBinaryExpr(eb *ExprBinary) interface{}
	VisitGroupingExpr(eg *ExprGrouping) interface{}
	VisitLiteralExpr(el *ExprLiteral) interface{}
	VisitUnaryExpr(eu *ExprUnary) interface{}
	VisitVarExpr(ev *ExprVar) interface{}
	VisitLogicalExpr(eb *ExprLogical) interface{}
}
