package lox

type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

type ExpressionStmt struct {
	Expression Expr
}

func NewExpressionStmt(expression Expr) Stmt {
	return &ExpressionStmt{Expression: expression}
}

func (stmt *ExpressionStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitExpressionStmt(stmt)
}

type PrintStmt struct {
	Expression Expr
}

func NewPrintStmt(expression Expr) Stmt {
	return &PrintStmt{Expression: expression}
}

func (stmt *PrintStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(stmt)
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func NewVarStmt(name Token, initializer Expr) Stmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

func (stmt *VarStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(stmt)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, thenbranch Stmt, elsebranch Stmt) Stmt {
	return &IfStmt{Condition: condition, ThenBranch: thenbranch, ElseBranch: elsebranch}
}

func (stmt *IfStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitIfStmt(stmt)
}

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) Stmt {
	return &BlockStmt{Statements: statements}
}

func (stmt *BlockStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(stmt)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) Stmt {
	return &WhileStmt{Condition: condition, Body: body}
}

func (stmt *WhileStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(stmt)
}

type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func NewFunctionStmt(name Token, params []Token, body []Stmt) Stmt {
	return &FunctionStmt{Name: name, Params: params, Body: body}
}

func (stmt *FunctionStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitFunctionStmt(stmt)
}

type ReturnStmt struct {
	Keyword Token
	Value   Expr
}

func NewReturnStmt(keyword Token, value Expr) Stmt {
	return &ReturnStmt{Keyword: keyword, Value: value}
}

func (stmt *ReturnStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitReturnStmt(stmt)
}
