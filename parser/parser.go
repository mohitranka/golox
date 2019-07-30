package parser

import (
	"fmt"
	"github.com/mohitranka/golox/err"
	"github.com/mohitranka/golox/expression"
	"github.com/mohitranka/golox/statement"
	"github.com/mohitranka/golox/token"
)

var current int

type Parser struct {
	tokens []*token.Token
}

func NewParser(tokens []*token.Token) *Parser {
	np := new(Parser)
	np.tokens = tokens
	current = 0
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

func (p Parser) advance() *token.Token {
	if !p.isAtEnd() {
		current++
	}
	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p Parser) peek() token.Token {
	return *p.tokens[current]
}

func (p Parser) previous() *token.Token {
	return p.tokens[current-1]
}

func (p Parser) Parse() []statement.Stmt {
	statements := make([]statement.Stmt, 0)
	for {
		if p.isAtEnd() {
			break
		}
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p Parser) declaration() statement.Stmt {
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p Parser) varDeclaration() statement.Stmt {
	name := p.consume(token.IDENTIFIER, "Expect variable name")
	var initializer expression.Expr
	if p.match(token.EQUAL) {
		initializer, _ = p.expression()
	}
	p.consume(token.SEMICOLON, "Expect ';' after variable declaration")
	return statement.NewVarStmt(*name, initializer)
}

func (p Parser) expression() (expression.Expr, error) {
	return p.assignment()
}

func (p Parser) assignment() (expression.Expr, error) {
	expr, e := p.equality()
	if e != nil {
		return nil, e
	}

	if p.match(token.EQUAL) {
		equals := p.previous()
		value, e := p.assignment()
		if e != nil {
			return nil, e
		}
		token, ok := expr.(*expression.ExprVar)
		if ok {
			name := token.Name
			return &expression.ExprAssign{Name: name, Value: value}, nil
		} else {
			e := &err.VarError{Name: equals.Lexeme, Msg: "Invalid assignment target"}
			fmt.Println(e)
			return nil, e
		}
	}
	return expr, nil
}

func (p Parser) statement() statement.Stmt {
	if p.match(token.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

func (p Parser) printStatement() statement.Stmt {
	value, _ := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after value.")
	return statement.NewPrintStmt(value)
}

func (p Parser) expressionStatement() statement.Stmt {
	expr, _ := p.expression()
	p.consume(token.SEMICOLON, "Expect ';' after expression.")
	return statement.NewExpressionStmt(expr)
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
		expr = &expression.ExprBinary{Left: expr, Operator: *operator, Right: right}
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
		expr = &expression.ExprBinary{Left: expr, Operator: *operator, Right: right}
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
		expr = &expression.ExprBinary{Left: expr, Operator: *operator, Right: right}
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
		expr = &expression.ExprBinary{Left: expr, Operator: *operator, Right: right}
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
		return &expression.ExprUnary{Operator: *operator, Right: right}, nil
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
	if p.match(token.IDENTIFIER) {
		return &expression.ExprVar{*p.previous()}, nil
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

func (p Parser) consume(token_type token.TokenType, message string) *token.Token {
	if !p.check(token_type) {
		fmt.Println(p.parse_err(p.peek(), message))
		return nil
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
