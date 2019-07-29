package parser

import (
	"github.com/mohitranka/golox/err"
	"github.com/mohitranka/golox/expression"
	"github.com/mohitranka/golox/token"
)

type Parser struct {
	tokens  []token.Token
	current int
}

func NewParser(tokens []token.Token, current int) *Parser {
	np := new(Parser)
	np.tokens = tokens
	np.current = current
	return np
}

func (p Parser) match(types ...token.TokenType) bool {
	for _, token_type := range types {
		if p.check(token_type) {
			p.advance()
			return true
		}
	}
	return false
}

func (p Parser) check(token_type token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == token_type
}

func (p Parser) advance() token.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p Parser) peek() token.Token {
	return p.tokens[p.current]
}

func (p Parser) previous() token.Token {
	return p.tokens[p.current-1]
}

func (p Parser) parse() expression.Expr {
	expr, e := p.expression()
	if e != nil {
		return nil
	}
	return expr
}

func (p Parser) expression() (expression.Expr, error) {
	return p.equality()
}

func (p Parser) equality() (expression.Expr, error) {
	expr, e := p.comparison()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
			break
		}
		operator := p.previous()
		right, e := p.comparison()
		if e != nil {
			return nil, e
		}
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p Parser) comparison() (expression.Expr, error) {
	expr, e := p.addition()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
			break
		}
		operator := p.previous()
		right, e := p.addition()
		if e != nil {
			return nil, e
		}
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p Parser) addition() (expression.Expr, error) {
	expr, e := p.multiplication()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(token.MINUS, token.PLUS) {
			break
		}
		operator := p.previous()
		right, e := p.multiplication()
		if e != nil {
			return nil, e
		}
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p Parser) multiplication() (expression.Expr, error) {
	expr, e := p.unary()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(token.SLASH, token.STAR) {
			break
		}

		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			return nil, e
		}
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

func (p Parser) unary() (expression.Expr, error) {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			return nil, e
		}
		return &expression.ExprUnary{Operator: operator, Right: right}, nil
	}
	return p.primary()
}

func (p Parser) primary() (expression.Expr, error) {
	if p.match(token.FALSE) {
		return &expression.ExprLiteral{false}, nil
	}

	if p.match(token.TRUE) {
		return &expression.ExprLiteral{true}, nil
	}

	if p.match(token.NIL) {
		return &expression.ExprLiteral{nil}, nil
	}

	if p.match(token.NUMBER, token.STRING) {
		return &expression.ExprLiteral{p.previous().Literal}, nil
	}
	if p.match(token.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &expression.ExprGrouping{expr}, nil
	}
	return nil, p.parse_err(p.peek(), "Expect expression.")
}

func (p Parser) consume(token_type token.TokenType, message string) token.Token {
	if !p.check(token_type) {
		panic(p.parse_err(p.peek(), message))
	}
	return p.advance()
}

func (p Parser) parse_err(t token.Token, message string) error {
	return &err.ParseError{Token: t, Msg: message}
}

// Synchronize

func (p Parser) synchronize() {
	p.advance()
	for {
		if p.isAtEnd() {
			break
		}
		if p.previous().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.CLASS, token.FUN, token.VAR, token.FOR, token.IF, token.WHILE, token.PRINT, token.RETURN:
			return
		}
		p.advance()
	}
}
