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

func (p Parser) expression() expression.Expr {
	return p.equality()
}

func (p Parser) equality() expression.Expr {
	expr := p.comparison()
	for {
		if !p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
			break
		}
		operator := p.previous()
		right := p.comparison()
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
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

func (p Parser) comparison() expression.Expr {
	expr := p.addition()
	for {
		if !p.match(token.GREATER, token.GREATER_EQUAL, token.LESS, token.LESS_EQUAL) {
			break
		}
		operator := p.previous()
		right := p.addition()
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p Parser) addition() expression.Expr {
	expr := p.multiplication()
	for {
		if !p.match(token.MINUS, token.PLUS) {
			break
		}
		operator := p.previous()
		right := p.multiplication()
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p Parser) multiplication() expression.Expr {
	expr := p.unary()
	for {
		if !p.match(token.SLASH, token.STAR) {
			break
		}

		operator := p.previous()
		right := p.unary()
		expr = &expression.ExprBinary{Left: expr, Operator: operator, Right: right}
	}
	return expr
}

func (p Parser) unary() expression.Expr {
	if p.match(token.BANG, token.MINUS) {
		operator := p.previous()
		right := p.unary()
		return &expression.ExprUnary{Operator: operator, Right: right}
	}
	return p.primary()
}

func (p Parser) primary() expression.Expr {
	if p.match(token.FALSE) {
		return &expression.ExprLiteral{false}
	}

	if p.match(token.TRUE) {
		return &expression.ExprLiteral{true}
	}

	if p.match(token.NIL) {
		return &expression.ExprLiteral{nil}
	}

	if p.match(token.NUMBER, token.STRING) {
		return &expression.ExprLiteral{p.previous().Literal}
	}
	if p.match(token.LEFT_PAREN) {
		expr := p.expression()
		p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		return &expression.ExprGrouping{expr}
	}
	return nil
}

func (p Parser) consume(token_type token.TokenType, message string) token.Token {
	if !p.check(token_type) {
		p.parse_err(p.peek(), message)
	}
	return p.advance()
}

func (p Parser) parse_err(t token.Token, message string) {
	panic(&err.ParseError{Token: t, Msg: message})
}
