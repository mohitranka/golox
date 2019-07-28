package expression

import (
	"github.com/mohitranka/golox/token"
)

type Expr interface {
	Expression() Expr
}

type ExprAssign struct {
	Name  token.Token
	Value Expr
}

type ExprBinary struct {
	Left     Expr
	Operator token.Token
	Right    Expr
}

type ExprGrouping struct {
	Expr
}

type ExprLiteral struct {
	Value interface{}
}

type ExprUnary struct {
	Operator token.Token
	Right    Expr
}

type ExprVar struct {
	Name token.Token
}

func (e *ExprAssign) Expression() Expr { return e }

func (e *ExprBinary) Expression() Expr { return e }

func (e *ExprGrouping) Expression() Expr { return e }

func (e *ExprLiteral) Expression() Expr { return e }

func (e *ExprUnary) Expression() Expr { return e }

func (e *ExprVar) Expression() Expr { return e }
