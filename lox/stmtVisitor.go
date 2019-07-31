package lox

type StmtVisitor interface {
	VisitExpressionStmt(stmt *ExpressionStmt) interface{}
	VisitPrintStmt(stmt *PrintStmt) interface{}
	VisitVarStmt(stmt *VarStmt) interface{}
	VisitBlockStmt(stmt *BlockStmt) interface{}
	VisitIfStmt(stmt *IfStmt) interface{}
	VisitWhileStmt(stmt *WhileStmt) interface{}
	VisitFunctionStmt(stmt *FunctionStmt) interface{}
	VisitReturnStmt(stmt *ReturnStmt) interface{}
}
