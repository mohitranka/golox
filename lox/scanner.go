package lox

import (
	"fmt"
	"strconv"
	"unicode"
)

var start int
var currentScannerPointer int
var line int
var tokens []*Token
var keywords map[string]TokenType

// Scanner ...
type Scanner struct {
	source string
}

// NewScanner ...
func NewScanner(source string) *Scanner {
	s := new(Scanner)
	s.source = source

	tokens = make([]*Token, 0)
	start = 0
	currentScannerPointer = 0
	line = 1

	keywords = map[string]TokenType{
		"and":    TokenTypeAnd,
		"class":  TokenTypeClass,
		"else":   TokenTypeElse,
		"false":  TokenTypeFalse,
		"for":    TokenTypeFor,
		"fun":    TokenTypeFun,
		"if":     TokenTypeIf,
		"nil":    TokenTypeNil,
		"or":     TokenTypeOr,
		"print":  TokenTypePrint,
		"return": TokenTypeReturn,
		"super":  TokenTypeSuper,
		"this":   TokenTypeThis,
		"true":   TokenTypeTrue,
		"var":    TokenTypeVar,
		"while":  TokenTypeWhile,
	}
	return s
}

// ScanTokens ...
func (s Scanner) ScanTokens() []*Token {
	for {
		if s.isAtEnd() {
			break
		}
		start = currentScannerPointer
		s.scanToken()
	}

	tokens = append(tokens, NewToken(TokenTypeEOF, "", nil, line))
	return tokens
}

func (s Scanner) isAtEnd() bool {
	return currentScannerPointer >= len(s.source)
}

func (s Scanner) scanToken() {
	switch c := s.advance(); c {
	case '(':
		s.addToken(TokenTypeLeftParen)
	case ')':
		s.addToken(TokenTypeRightParen)
	case '{':
		s.addToken(TokenTypeLeftBrace)
	case '}':
		s.addToken(TokenTypeRightBrace)
	case ',':
		s.addToken(TokenTypeComma)
	case '.':
		s.addToken(TokenTypeDot)
	case '-':
		s.addToken(TokenTypeMinus)
	case '+':
		s.addToken(TokenTypePlus)
	case ';':
		s.addToken(TokenTypeSemiColon)
	case '*':
		s.addToken(TokenTypeStar)
	case '!':
		if s.match('=') {
			s.addToken(TokenTypeBangEqual)
		} else {
			s.addToken(TokenTypeBang)
		}
	case '=':
		if s.match('=') {
			s.addToken(TokenTypeEqualEqual)
		} else {
			s.addToken(TokenTypeEqual)
		}
	case '<':
		if s.match('=') {
			s.addToken(TokenTypeLessEqual)
		} else {
			s.addToken(TokenTypeLess)
		}
	case '>':
		if s.match('=') {
			s.addToken(TokenTypeGreaterEqual)
		} else {
			s.addToken(TokenTypeGreater)
		}
	case '/':
		if s.match('/') {
			for {
				if s.isAtEnd() || s.peek() == '\n' {
					break
				}
				s.advance()
			}
		} else {
			s.addToken(TokenTypeSlash)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		line++
	case '"':
		s.stringTokenizer()
	default:
		if s.isDigit(c) {
			s.numberTokenizer()
		} else if s.isAlpha(c) {
			s.identifierTokenizer()
		} else {
			fmt.Println(&RuntimeError{line, "Unexpected character."})
			return
		}
	}
}

func (s Scanner) isAlpha(c byte) bool {
	return unicode.IsLetter(rune(c)) || c == '_'
	//	return !((c < 'a' || c > 'z') && (c < 'A' && c > 'Z') && c != '_')
}

func (s Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s Scanner) identifierTokenizer() {
	for {
		if !s.isAlphaNumeric(s.peek()) {
			break
		}
		s.advance()
	}
	tokenType, ok := keywords[s.source[start:currentScannerPointer]]
	if !ok {
		tokenType = TokenTypeIdentifier
	}
	s.addToken(tokenType)
}

func (s Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s Scanner) numberTokenizer() error {
	for {
		if !s.isDigit(s.peek()) {
			break
		}
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for {
			if !s.isDigit(s.peek()) {
				break
			}
			s.advance()
		}
	}

	value, e := strconv.ParseFloat(s.source[start:currentScannerPointer], 64)
	if e != nil {
		return &RuntimeError{line, e.Error()}
	}
	s.addTokenWithLiteral(TokenTypeNumber, value)
	return nil
}

func (s Scanner) peekNext() byte {
	if currentScannerPointer+1 >= len(s.source) {
		return 0
	}

	return s.source[currentScannerPointer+1]
}

func (s Scanner) stringTokenizer() error {
	for {
		if s.isAtEnd() || s.peek() == '"' {
			break
		}
		if s.peek() == '\n' {
			line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		return &RuntimeError{line, "Unterminated String"}
	}

	s.advance()
	text := s.source[start+1 : currentScannerPointer-1]
	s.addTokenWithLiteral(TokenTypeString, text)
	return nil
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[currentScannerPointer]
}

func (s Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[currentScannerPointer] != expected {
		return false
	}
	currentScannerPointer++
	return true
}

func (s Scanner) advance() byte {
	currentScannerPointer++
	return s.source[currentScannerPointer-1]
}

func (s Scanner) addToken(tokenType TokenType) {
	s.addTokenWithLiteral(tokenType, nil)
}

func (s Scanner) addTokenWithLiteral(tokenType TokenType, literal interface{}) {
	text := s.source[start:currentScannerPointer]
	tokens = append(tokens, NewToken(tokenType, text, literal, line))
}
