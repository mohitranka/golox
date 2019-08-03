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
	return p.peek().Type == TokenTypeEOF
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
	if p.match(TokenTypeFun) {
		return p.function("function")
	}
	if p.match(TokenTypeVar) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p Parser) function(kind string) Stmt {
	name := p.consume(TokenTypeIdentifier, fmt.Sprintf("Expect %s name.", kind))
	p.consume(TokenTypeLeftParen, fmt.Sprintf("Expect '(' after %s name.", kind))
	params := make([]Token, 0)
	if !p.check(TokenTypeRightParen) {
		for {
			params = append(params, *p.consume(TokenTypeIdentifier, "Expect ')' after parameters."))
			if !p.match(TokenTypeComma) {
				break
			}
		}
	}
	p.consume(TokenTypeRightParen, "Expect ')' after parameters")
	p.consume(TokenTypeLeftBrace, fmt.Sprintf("Expect '{' before %s body.", kind))
	body := p.block()
	return NewFunctionStmt(*name, params, body)
}

func (p Parser) varDeclaration() Stmt {
	name := p.consume(TokenTypeIdentifier, "Expect variable name")
	var initializer Expr
	if p.match(TokenTypeEqual) {
		initializer, _ = p.expression()
	}
	p.consume(TokenTypeSemiColon, "Expect ';' after variable declaration")
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

	if p.match(TokenTypeEqual) {
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
		if !p.match(TokenTypeOr) {
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
		if !p.match(TokenTypeAnd) {
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
	if p.match(TokenTypeFor) {
		return p.forStatement()
	}
	if p.match(TokenTypeIf) {
		return p.ifStatement()
	}
	if p.match(TokenTypePrint) {
		return p.printStatement()
	}
	if p.match(TokenTypeReturn) {
		return p.returnStatement()
	}
	if p.match(TokenTypeWhile) {
		return p.whileStatement()
	}
	if p.match(TokenTypeLeftBrace) {
		return NewBlockStmt(p.block())
	}
	return p.expressionStatement()
}

func (p Parser) returnStatement() Stmt {
	keyword := p.previous()
	var value Expr
	var e error
	if !p.check(TokenTypeSemiColon) {
		value, e = p.expression()
		if e != nil {
			fmt.Println("Error while getting expression for the return statement.")
			return nil
		}
	}
	p.consume(TokenTypeSemiColon, "Expect ';' after return value.")
	return NewReturnStmt(*keyword, value)
}

func (p Parser) block() []Stmt {
	statements := make([]Stmt, 0)
	for {
		if p.check(TokenTypeRightBrace) || p.isAtEnd() {
			break
		}
		statements = append(statements, p.declaration())
	}
	p.consume(TokenTypeRightBrace, "Expect '}' after block.")
	return statements
}

func (p Parser) ifStatement() Stmt {
	p.consume(TokenTypeLeftParen, "Expect '(' after 'if'.")
	condition, _ := p.expression()
	p.consume(TokenTypeRightParen, "Expect ')' after if condition")
	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(TokenTypeElse) {
		elseBranch = p.statement()
	}
	return NewIfStmt(condition, thenBranch, elseBranch)
}

func (p Parser) printStatement() Stmt {
	value, _ := p.expression()
	p.consume(TokenTypeSemiColon, "Expect ';' after value.")
	return NewPrintStmt(value)
}

func (p Parser) forStatement() Stmt {
	p.consume(TokenTypeLeftParen, "Expect '(' after for")
	var initializer Stmt
	if p.match(TokenTypeSemiColon) {
		initializer = nil
	} else if p.match(TokenTypeVar) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(TokenTypeSemiColon) {
		condition, _ = p.expression()
	}
	p.consume(TokenTypeSemiColon, "Expect ';' after loop condition")

	var increment Expr

	if !p.check(TokenTypeRightParen) {
		increment, _ = p.expression()
	}
	p.consume(TokenTypeRightParen, "Expect ')' after for clauses")
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
	p.consume(TokenTypeLeftParen, "Expect '(' after while.")
	condition, _ := p.expression()
	p.consume(TokenTypeRightParen, "Expect ')' after condition")
	body := p.statement()
	return NewWhileStmt(condition, body)
}

func (p Parser) expressionStatement() Stmt {
	expr, _ := p.expression()
	p.consume(TokenTypeSemiColon, "Expect ';' after expression.")
	return NewExpressionStmt(expr)
}

func (p Parser) equality() (Expr, error) {
	expr, e := p.comparison()
	if e != nil {
		return nil, e
	}
	for {
		if !p.match(TokenTypeBangEqual, TokenTypeEqualEqual) {
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
		if !p.match(TokenTypeGreater, TokenTypeGreaterEqual, TokenTypeLess, TokenTypeLessEqual) {
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
		if !p.match(TokenTypeMinus, TokenTypePlus) {
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
		if !p.match(TokenTypeSlash, TokenTypeStar) {
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
	if p.match(TokenTypeBang, TokenTypeMinus) {
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
		if p.match(TokenTypeLeftParen) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}
	return expr, nil
}

func (p Parser) finishCall(callee Expr) Expr {
	arguments := make([]*Expr, 0)
	if !p.check(TokenTypeRightParen) {
		for {
			thisExpr, e := p.expression()
			if e != nil {
				fmt.Println(e)
				return nil
			}
			arguments = append(arguments, &thisExpr)
			if !p.match(TokenTypeComma) {
				break
			}
		}
	}
	paren := p.consume(TokenTypeRightParen, "Expect ')' after arguments")
	return &ExprCall{Callee: callee, Paren: *paren, Arguments: arguments}
}

func (p Parser) primary() (Expr, error) {
	if p.match(TokenTypeFalse) {
		return &ExprLiteral{false}, nil
	}

	if p.match(TokenTypeTrue) {
		return &ExprLiteral{true}, nil
	}

	if p.match(TokenTypeNil) {
		return &ExprLiteral{nil}, nil
	}

	if p.match(TokenTypeNumber, TokenTypeString) {
		return &ExprLiteral{p.previous().Literal}, nil
	}
	if p.match(TokenTypeIdentifier) {
		return &ExprVar{*p.previous()}, nil
	}
	if p.match(TokenTypeLeftParen) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		p.consume(TokenTypeRightParen, "Expect ')' after expression.")
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
		if p.previous().Type == TokenTypeSemiColon {
			return
		}
		switch p.peek().Type {
		case TokenTypeClass, TokenTypeFun, TokenTypeVar, TokenTypeFor, TokenTypeIf, TokenTypeWhile, TokenTypePrint, TokenTypeReturn:
			return
		}
		p.advance()
	}
}
