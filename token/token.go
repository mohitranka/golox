package token

import (
	"fmt"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(token_type TokenType, lexeme string, literal interface{}, line int) *Token {
	t := new(Token)
	t.Type = token_type
	t.Lexeme = lexeme
	t.Literal = literal
	t.Line = line
	return t
}

func (t Token) String() string {
	return fmt.Sprintf("TokenType:%s Lexeme:%s Literal:%v", t.Type, t.Lexeme, t.Literal)
}
