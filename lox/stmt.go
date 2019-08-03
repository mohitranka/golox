package lox

// Stmt ...
type Stmt interface {
	Accept(v StmtVisitor) interface{}
}

// ExpressionStmt ...
type ExpressionStmt struct {
	Expression Expr
}

// NewExpressionStmt ...
func NewExpressionStmt(expression Expr) Stmt {
	return &ExpressionStmt{Expression: expression}
}

// Accept ...
func (stmt *ExpressionStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitExpressionStmt(stmt)
}

// PrintStmt ...
type PrintStmt struct {
	Expression Expr
}

// NewPrintStmt ...
func NewPrintStmt(expression Expr) Stmt {
	return &PrintStmt{Expression: expression}
}

// Accept ...
func (stmt *PrintStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(stmt)
}

// VarStmt ...
type VarStmt struct {
	Name        Token
	Initializer Expr
}

// NewVarStmt ...
func NewVarStmt(name Token, initializer Expr) Stmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

// Accept ...
func (stmt *VarStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(stmt)
}

// IfStmt ...
type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

// NewIfStmt ...
func NewIfStmt(condition Expr, thenbranch Stmt, elsebranch Stmt) Stmt {
	return &IfStmt{Condition: condition, ThenBranch: thenbranch, ElseBranch: elsebranch}
}

// Accept ...
func (stmt *IfStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitIfStmt(stmt)
}

// BlockStmt ...
type BlockStmt struct {
	Statements []Stmt
}

// NewBlockStmt ...
func NewBlockStmt(statements []Stmt) Stmt {
	return &BlockStmt{Statements: statements}
}

// Accept ...
func (stmt *BlockStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(stmt)
}

// WhileStmt ...
type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

// NewWhileStmt ...
func NewWhileStmt(condition Expr, body Stmt) Stmt {
	return &WhileStmt{Condition: condition, Body: body}
}

// Accept ...
func (stmt *WhileStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(stmt)
}

// FunctionStmt ...
type FunctionStmt struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

// NewFunctionStmt ...
func NewFunctionStmt(name Token, params []Token, body []Stmt) Stmt {
	return &FunctionStmt{Name: name, Params: params, Body: body}
}

// Accept ...
func (stmt *FunctionStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitFunctionStmt(stmt)
}

// ReturnStmt ...
type ReturnStmt struct {
	Keyword Token
	Value   Expr
}

// NewReturnStmt ...
func NewReturnStmt(keyword Token, value Expr) Stmt {
	return &ReturnStmt{Keyword: keyword, Value: value}
}

// Accept ...
func (stmt *ReturnStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitReturnStmt(stmt)
}
