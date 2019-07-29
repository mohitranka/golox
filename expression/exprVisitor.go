package expression

type ExprVisitor interface {
	VisitAssignExpr(eb *ExprAssign) interface{}
	VisitBinaryExpr(eb *ExprBinary) interface{}
	VisitGroupingExpr(eg *ExprGrouping) interface{}
	VisitLiteralExpr(el *ExprLiteral) interface{}
	VisitUnaryExpr(eb *ExprUnary) interface{}
	VisitVarExpr(eb *ExprVar) interface{}
}
