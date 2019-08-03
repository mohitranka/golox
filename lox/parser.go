package lox

import (
	"fmt"
)

var currentParserPointer int

// Parser ...
type Parser struct {
	tokens []*Token
}

// NewParser ...
func NewParser(tokens []*Token) *Parser {
	np := new(Parser)
	np.tokens = tokens
	currentParserPointer = 0
	return np
}

func (p Parser) match(types ...TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

func (p Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tokenType
}

func (p Parser) advance() *Token {
	if !p.isAtEnd() {
		currentParserPointer++
	}
	return p.previous()
}

func (p Parser) isAtEnd() bool {
	return p.peek().Type == EOF
}

func (p Parser) peek() Token {
	return *p.tokens[currentParserPointer]
}

func (p Parser) previous() *Token {
	return p.tokens[currentParserPointer-1]
}

// Parse ...
func (p Parser) Parse() []Stmt {
	statements := make([]Stmt, 0)
	for {
		if p.isAtEnd() {
			break
		}
		statements = append(statements, p.declaration())
	}
	return statements
}

func (p Parser) declaration() Stmt {
	if p.match(FUN) {
		return p.function("function")
	}
	if p.match(VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p Parser) function(kind string) Stmt {
	name := p.consume(IDENTIFIER, fmt.Sprintf("Expect %s name.", kind))
	p.consume(LEFT_PAREN, fmt.Sprintf("Expect '(' after %s name.", kind))
	params := make([]Token, 0)
	if !p.check(RIGHT_PAREN) {
		for {
			params = append(params, *p.consume(IDENTIFIER, "Expect ')' after parameters."))
			if !p.match(COMMA) {
				break
			}
		}
	}
	p.consume(RIGHT_PAREN, "Expect ')' after parameters")
	p.consume(LEFT_BRACE, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return NewFunctionStmt(*name, params, body)
}

func (p Parser) varDeclaration() Stmt {
	name := p.consume(IDENTIFIER, "Expect variable name")
	var initializer Expr
	if p.match(EQUAL) {
		initializer, _ = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after variable declaration")
	return NewVarStmt(*name, initializer)
}

func (p Parser) expression() (Expr, error) {
	return p.assignment()
}

func (p Parser) assignment() (Expr, error) {
	expr, e := p.or()
	if e != nil {
		return nil, e
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, e := p.assignment()
		if e != nil {
			return nil, e
		}
		token, ok := expr.(*ExprVar)
		if ok {
			name := token.Name
			return &ExprAssign{Name: name, Value: value}, nil
		}
		e = &VarError{Name: equals.Lexeme, Msg: "Invalid assignment target"}
		fmt.Println(e)
		return nil, e
	}
	return expr, nil
}

func (p Parser) or() (Expr, error) {
	expr, e := p.and()
	if e != nil {
		return nil, e
	}

	for {
		if !p.match(OR) {
			break
		}

		operator := p.previous()
		right, e := p.and()
		if e != nil {
			return nil, e
		}
		expr = &ExprLogical{Left: expr, Operator: *operator, Right: right}
	}
	return expr, nil
}

func (p Parser) and() (Expr, error) {
	expr, e := p.equality()
	if e != nil {
		return nil, e
	}

	for {
		if !p.match(AND) {
			break
		}
		operator := p.previous()
		right, e := p.equality()
		if e != nil {
			return nil, e
		}
		expr = &ExprLogical{Left: expr, Operator: *operator, Right: right}
	}
	return expr, nil
}

func (p Parser) statement() Stmt {
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(RETURN) {
		return p.returnStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(LEFT_BRACE) {
		return NewBlockStmt(p.block())
	}
	return p.expressionStatement()
}

func (p Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	var e error
	if !p.check(SEMICOLON) {
		value, e = p.expression()
		if e != nil {
			fmt.Println("Error while getting expression for the return statement.")
			return nil
		}
	}
	p.consume(SEMICOLON, "Expect ';' after return value.")
	return NewReturnStmt(*keyword, value)
}

func (p Parser) block() []Stmt {
	statements := make([]Stmt, 0)
	for {
		if p.check(RIGHT_BRACE) || p.isAtEnd() {
			break
		}
		statements = append(statements, p.declaration())
	}
	p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return statements
}

func (p Parser) ifStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	condition, _ := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after if condition")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch = p.statement()
	}
	return NewIfStmt(condition, thenBranch, elseBranch)
}

func (p Parser) printStatement() Stmt {
	value, _ := p.expression()
	p.consume(SEMICOLON, "Expect ';' after value.")
	return NewPrintStmt(value)
}

func (p Parser) forStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after for")
	var initializer Stmt
	if p.match(SEMICOLON) {
		initializer = nil
	} else if p.match(VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(SEMICOLON) {
		condition, _ = p.expression()
	}
	p.consume(SEMICOLON, "Expect ';' after loop condition")

	var increment Expr

	if !p.check(RIGHT_PAREN) {
		increment, _ = p.expression()
	}
	p.consume(RIGHT_PAREN, "Expect ')' after for clauses")
	body := p.statement()

	if increment != nil {
		statements := make([]Stmt, 0)
		statements = append(statements, body, NewExpressionStmt(increment))
		body = NewBlockStmt(statements)
	}

	if condition == nil {
		condition = &ExprLiteral{true}
	}

	body = NewWhileStmt(condition, body)

	if initializer != nil {
		statements := make([]Stmt, 0)
		statements = append(statements, initializer, body)
		body = NewBlockStmt(statements)
	}
	return body
}

func (p Parser) whileStatement() Stmt {
	p.consume(LEFT_PAREN, "Expect '(' after while.")
	condition, _ := p.expression()
	p.consume(RIGHT_PAREN, "Expect ')' after condition")
	body := p.statement()
	return NewWhileStmt(condition, body)
}

func (p Parser) expressionStatement() Stmt {
	expr, _ := p.expression()
	p.consume(SEMICOLON, "Expect ';' after expression.")
	return NewExpressionStmt(expr)
}

func (p Parser) equality() (Expr, error) {
	expr, e := p.comparison()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(BANG_EQUAL, EQUAL_EQUAL) {
			break
		}
		operator := p.previous()
		right, e := p.comparison()
		if e != nil {
			return nil, e
		}
		expr = &ExprBinary{Left: expr, Operator: *operator, Right: right}
	}
	return expr, nil
}

func (p Parser) comparison() (Expr, error) {
	expr, e := p.addition()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
			break
		}
		operator := p.previous()
		right, e := p.addition()
		if e != nil {
			return nil, e
		}
		expr = &ExprBinary{Left: expr, Operator: *operator, Right: right}
	}
	return expr, nil
}

func (p Parser) addition() (Expr, error) {
	expr, e := p.multiplication()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(MINUS, PLUS) {
			break
		}
		operator := p.previous()
		right, e := p.multiplication()
		if e != nil {
			return nil, e
		}
		expr = &ExprBinary{Left: expr, Operator: *operator, Right: right}
	}
	return expr, nil
}

func (p Parser) multiplication() (Expr, error) {
	expr, e := p.unary()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(SLASH, STAR) {
			break
		}

		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			return nil, e
		}
		expr = &ExprBinary{Left: expr, Operator: *operator, Right: right}
	}
	return expr, nil
}

func (p Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		operator := p.previous()
		right, e := p.unary()
		if e != nil {
			return nil, e
		}
		return &ExprUnary{Operator: *operator, Right: right}, nil
	}
	return p.call()
}

func (p Parser) call() (Expr, error) {
	expr, e := p.primary()
	if e != nil {
		return nil, e
	}
	for {
		if p.match(LEFT_PAREN) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr, nil
}

func (p Parser) finishCall(callee Expr) Expr {
	arguments := make([]*Expr, 0)
	if !p.check(RIGHT_PAREN) {
		for {
			thisExpr, e := p.expression()
			if e != nil {
				fmt.Println(e)
				return nil
			}
			arguments = append(arguments, &thisExpr)
			if !p.match(COMMA) {
				break
			}
		}
	}
	paren := p.consume(RIGHT_PAREN, "Expect ')' after arguments")
	return &ExprCall{Callee: callee, Paren: *paren, Arguments: arguments}
}

func (p Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return &ExprLiteral{false}, nil
	}

	if p.match(TRUE) {
		return &ExprLiteral{true}, nil
	}

	if p.match(NIL) {
		return &ExprLiteral{nil}, nil
	}

	if p.match(NUMBER, STRING) {
		return &ExprLiteral{p.previous().Literal}, nil
	}
	if p.match(IDENTIFIER) {
		return &ExprVar{*p.previous()}, nil
	}
	if p.match(LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return &ExprGrouping{expr}, nil
	}
	return nil, p.parseErr(p.peek(), "Expect expression.")
}

func (p Parser) consume(tokenType TokenType, message string) *Token {
	if !p.check(tokenType) {
		fmt.Println(p.parseErr(p.peek(), message))
		return nil
	}
	return p.advance()
}

func (p Parser) parseErr(t Token, message string) error {
	return &ParseError{Token: t, Msg: message}
}

// Synchronize

func (p Parser) synchronize() {
	p.advance()
	for {
		if p.isAtEnd() {
			break
		}
		if p.previous().Type == SEMICOLON {
			return
		}
		switch p.peek().Type {
		case CLASS, FUN, VAR, FOR, IF, WHILE, PRINT, RETURN:
			return
		}
		p.advance()
	}
}
