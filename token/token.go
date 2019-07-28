package token

import (
	"fmt"
)

type Token struct {
	Ttype   TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(token_type TokenType, lexeme string, literal interface{}, line int) *Token {
	t := new(Token)
	t.Ttype = token_type
	t.Lexeme = lexeme
	t.Literal = literal
	t.Line = line
	return t
}

func (t Token) String() string {
	return fmt.Sprintf("TokenType:%s Lexeme:%s Literal:%v", t.Ttype, t.Lexeme, t.Literal)
}
