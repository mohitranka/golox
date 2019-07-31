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

func (expr *ExpressionStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitExpressionStmt(expr)
}

type PrintStmt struct {
	Expression Expr
}

func NewPrintStmt(expression Expr) Stmt {
	return &PrintStmt{Expression: expression}
}

func (expr *PrintStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitPrintStmt(expr)
}

type VarStmt struct {
	Name        Token
	Initializer Expr
}

func NewVarStmt(name Token, initializer Expr) Stmt {
	return &VarStmt{Name: name, Initializer: initializer}
}

func (expr *VarStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitVarStmt(expr)
}

type IfStmt struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func NewIfStmt(condition Expr, thenbranch Stmt, elsebranch Stmt) Stmt {
	return &IfStmt{Condition: condition, ThenBranch: thenbranch, ElseBranch: elsebranch}
}

func (expr *IfStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitIfStmt(expr)
}

type BlockStmt struct {
	Statements []Stmt
}

func NewBlockStmt(statements []Stmt) Stmt {
	return &BlockStmt{Statements: statements}
}

func (expr *BlockStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitBlockStmt(expr)
}

type WhileStmt struct {
	Condition Expr
	Body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) Stmt {
	return &WhileStmt{Condition: condition, Body: body}
}

func (expr *WhileStmt) Accept(v StmtVisitor) interface{} {
	return v.VisitWhileStmt(expr)
}
