package statement

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitVarStmt(stmt *VarStmt) interface{}
}
