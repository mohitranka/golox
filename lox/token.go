package lox

import (
	"fmt"
)

// Token ...
type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

// NewToken ...
func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	t := new(Token)
	t.Type = tokenType
	t.Lexeme = lexeme
	t.Literal = literal
	t.Line = line
	return t
}

// String ...
func (t Token) String() string {
	return fmt.Sprintf("TokenType:%s Lexeme:%s Literal:%v", t.Type, t.Lexeme, t.Literal)
}
