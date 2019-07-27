package token

import (
	"fmt"
)

type Token struct {
	token_type TokenType
	lexeme     string
	literal    interface{}
	line       int
}

func NewToken(token_type TokenType, lexeme string, literal interface{}, line int) *Token {
	t := new(Token)
	t.token_type = token_type
	t.lexeme = lexeme
	t.literal = literal
	t.line = line
	return t
}

func (t Token) String() string {
	return fmt.Sprintf("TokenType:%s Lexeme:%s Literal:%v", t.token_type, t.lexeme, t.literal)
}
