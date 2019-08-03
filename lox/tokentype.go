package lox

// TokenType ...
type TokenType int

// const ...
const (

	// 1 character token

	LEFT_PAREN TokenType = iota + 1
	RIGHT_PAREN
	LEFT_BRACE
	RIGHT_BRACE
	COMMA
	DOT
	MINUS
	PLUS
	SEMICOLON
	SLASH
	STAR

	//1 or 2 character token

	BANG
	BANG_EQUAL
	EQUAL
	EQUAL_EQUAL
	GREATER
	GREATER_EQUAL
	LESS
	LESS_EQUAL

	//Literals

	IDENTIFIER
	STRING
	NUMBER

	// KEYWORDS

	AND
	CLASS
	ELSE
	FALSE
	FUN
	FOR
	IF
	NIL
	OR
	PRINT
	RETURN
	SUPER
	THIS
	VAR
	WHILE
	TRUE
	EOF
)
