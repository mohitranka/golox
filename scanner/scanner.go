package scanner

import (
	"github.com/mohitranka/golox/err"
	"github.com/mohitranka/golox/token"
	"strconv"
	"unicode"
)

var start int
var current int
var line int
var tokens []*token.Token
var keywords map[string]token.TokenType

type Scanner struct {
	source string
}

func NewScanner(source string) *Scanner {
	s := new(Scanner)
	s.source = source

	tokens = make([]*token.Token, 0)
	start = 0
	current = 0
	line = 1

	keywords = map[string]token.TokenType{
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
	for {
		if s.isAtEnd() {
			break
		}
		start = current
		s.scanToken()
	}

	tokens = append(tokens, token.NewToken(token.EOF, "", nil, line))
	return tokens
}

func (s Scanner) isAtEnd() bool {
	return current >= len(s.source)
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
			for {
				if s.isAtEnd() || s.peek() == '\n' {
					break
				}
				s.advance()
			}
		} else {
			s.addToken(token.SLASH)
		}
	case ' ', '\r', '\t':
		break
	case '\n':
		line += 1
	case '"':
		s.stringTokenizer()
	default:
		if s.isDigit(c) {
			s.numberTokenizer()
		} else if s.isAlpha(c) {
			s.identifierTokenizer()
		} else {
			err.Error(line, "Unexpected character.")
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
	token_type, ok := keywords[s.source[start:current]]
	if !ok {
		token_type = token.IDENTIFIER
	}
	s.addToken(token_type)
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

	value, e := strconv.ParseFloat(s.source[start:current], 64)
	if e != nil {
		err.Error(line, "Invalid number")
		return e
	}
	s.addTokenWithLiteral(token.NUMBER, value)
	return nil
}

func (s Scanner) peekNext() byte {
	if current+1 >= len(s.source) {
		return 0
	}

	return s.source[current+1]
}

func (s Scanner) stringTokenizer() {
	for {
		if s.isAtEnd() || s.peek() == '"' {
			break
		}
		if s.peek() == '\n' {
			line += 1
		}
		s.advance()
	}

	if s.isAtEnd() {
		err.Error(line, "Unterminated String")
		return
	}

	s.advance()
	text := s.source[start+1 : current]
	s.addTokenWithLiteral(token.STRING, text)
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	}
	return s.source[current]
}

func (s Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	}

	if s.source[current] != expected {
		return false
	}
	current += 1
	return true
}

func (s Scanner) advance() byte {
	current += 1
	return s.source[current-1]
}

func (s Scanner) addToken(token_type token.TokenType) {
	s.addTokenWithLiteral(token_type, nil)
}

func (s Scanner) addTokenWithLiteral(token_type token.TokenType, literal interface{}) {
	text := s.source[start:current]
	tokens = append(tokens, token.NewToken(token_type, text, literal, line))
}
