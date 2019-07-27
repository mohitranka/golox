package scanner

import (
	"github.com/mohitranka/golox/err"
	"github.com/mohitranka/golox/token"
	"strconv"
)

type Scanner struct {
	source   string
	tokens   []*token.Token
	start    int
	current  int
	line     int
	keywords map[string]token.TokenType
}

func NewScanner(source string) *Scanner {
	s := new(Scanner)
	s.source = source
	s.tokens = make([]*token.Token, 0)
	s.start = 0
	s.current = 0
	s.line = 1
	s.keywords = map[string]token.TokenType{
		"and":    token.AND,
		"class":  token.CLASS,
		"else":   token.ELSE,
		"false":  token.FALSE,
		"for":    token.FOR,
		"fun":    token.FUN,
		"if":     token.IF,
		"nil":    token.NIL,
		"or":     token.OR,
		"print":  token.PRINT,
		"return": token.RETURN,
		"super":  token.SUPER,
		"this":   token.THIS,
		"true":   token.TRUE,
		"var":    token.VAR,
		"while":  token.WHILE,
	}
	return s
}

func (s Scanner) ScanTokens() []*token.Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}

	s.tokens = append(s.tokens, token.NewToken(token.EOF, "", struct{}{}, s.line))
	return s.tokens
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s Scanner) scanToken() {
	switch c := s.advance(); c {
	case '(':
		s.addToken(token.LEFT_PAREN)
	case ')':
		s.addToken(token.RIGHT_PAREN)
	case '{':
		s.addToken(token.LEFT_BRACE)
	case '}':
		s.addToken(token.RIGHT_BRACE)
	case ',':
		s.addToken(token.COMMA)
	case '.':
		s.addToken(token.DOT)
	case '-':
		s.addToken(token.MINUS)
	case '+':
		s.addToken(token.PLUS)
	case ';':
		s.addToken(token.SEMICOLON)
	case '*':
		s.addToken(token.STAR)
	case '!':
		if s.match('=') {
			s.addToken(token.BANG_EQUAL)
		} else {
			s.addToken(token.BANG)
		}
	case '=':
		if s.match('=') {
			s.addToken(token.EQUAL_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(token.LESS_EQUAL)
		} else {
			s.addToken(token.EQUAL)
		}
	case '>':
		if s.match('=') {
			s.addToken(token.GREATER_EQUAL)
		} else {
			s.addToken(token.GREATER)
		}
	case '/':
		if s.match('/') {
			for !s.isAtEnd() && s.peek() != '\n' {
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		s.line++
	case '"':
		s.stringTokenizer()
	default:
		if s.isDigit(c) {
			s.numberTokenizer()
		} else if s.isAlpha(c) {
			s.identifierTokenizer()
		} else {
			err.Error(s.line, "Unexpected character.")
		}
	}
}

func (s Scanner) isAlpha(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' || c <= 'Z') || c == '_'
}

func (s Scanner) isAlphaNumeric(c byte) bool {
	return s.isAlpha(c) || s.isDigit(c)
}

func (s Scanner) identifierTokenizer() {
	for s.isAlphaNumeric(s.peek()) {
		s.advance()
	}
	token_type, ok := s.keywords[s.source[s.start:s.current+1]]
	if !ok {
		token_type = token.IDENTIFIER
	}
	s.addToken(token_type)
}

func (s Scanner) isDigit(c byte) bool {
	return c >= '0' && c <= '9'
}

func (s Scanner) numberTokenizer() error {
	for s.isDigit(s.peek()) {
		s.advance()
	}

	if s.peek() == '.' && s.isDigit(s.peekNext()) {
		s.advance()
		for s.isDigit(s.peek()) {
			s.advance()
		}
	}

	value, e := strconv.ParseFloat(s.source[s.start:s.current+1], 64)
	if e != nil {
		err.Error(s.line, "Invalid number")
		return e
	}
	s.addTokenWithLiteral(token.NUMBER, value)
	return nil
}

func (s Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	}

	return s.source[s.current+1]
}

func (s Scanner) stringTokenizer() {
	for !s.isAtEnd() && s.peek() != '"' {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		err.Error(s.line, "Unterminated String")
	}

	s.advance()

	text := s.source[s.start+1 : s.current]
	s.addTokenWithLiteral(token.STRING, text)
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s Scanner) advance() byte {
	s.current++
	return s.source[s.current-1]
}

func (s Scanner) addToken(token_type token.TokenType) {
	s.addTokenWithLiteral(token_type, nil)
}

func (s Scanner) addTokenWithLiteral(token_type token.TokenType, literal interface{}) {
	text := s.source[s.start : s.current+1]
	s.tokens = append(s.tokens, token.NewToken(token_type, text, literal, s.line))
}
