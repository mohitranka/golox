package lox

// TokenType ...
type TokenType int

// const ...
const (

	// 1 character token

	TokenTypeLeftParen TokenType = iota + 1
	TokenTypeRightParen
	TokenTypeLeftBrace
	TokenTypeRightBrace
	TokenTypeComma
	TokenTypeDot
	TokenTypeMinus
	TokenTypePlus
	TokenTypeSemiColon
	TokenTypeSlash
	TokenTypeStar

	//1 or 2 character token

	TokenTypeBang
	TokenTypeBangEqual
	TokenTypeEqual
	TokenTypeEqualEqual
	TokenTypeGreater
	TokenTypeGreaterEqual
	TokenTypeLess
	TokenTypeLessEqual

	//Literals

	TokenTypeIdentifier
	TokenTypeString
	TokenTypeNumber

	// KEYWORDS

	TokenTypeAnd
	TokenTypeClass
	TokenTypeElse
	TokenTypeFalse
	TokenTypeFun
	TokenTypeFor
	TokenTypeIf
	TokenTypeNil
	TokenTypeOr
	TokenTypePrint
	TokenTypeReturn
	TokenTypeSuper
	TokenTypeThis
	TokenTypeVar
	TokenTypeWhile
	TokenTypeTrue
	TokenTypeEOF
)
