package expression

type exprVisitor interface {
	visitBinary(eb *ExprBinary) interface{}
	visitGrouping(eg *ExprGrouping) interface{}
	visitLiteral(el *ExprLiteral) interface{}
	visitUnary(eb *ExprUnary) interface{}
}
