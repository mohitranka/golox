package lox

import (
	"fmt"
)

func report(line int, where string, message string) string {
	return fmt.Sprintf("[line %d ] Error %s: %s\n", line, where, message)
}

// RuntimeError ...
type RuntimeError struct {
	Line int
	Msg  string
}

// ParseError ...
type ParseError struct {
	Token Token
	Msg   string
}

// VarError ...
type VarError struct {
	Name string
	Msg  string
}

// Error ...
func (re *RuntimeError) Error() string {
	return report(re.Line, "", re.Msg)
}

// Error ...
func (pe *ParseError) Error() string {
	if pe.Token.Type == TokenTypeEOF {
		return report(pe.Token.Line, " at end", pe.Msg)
	}
	return report(pe.Token.Line, fmt.Sprintf(" at '%s'", pe.Token.Lexeme), pe.Msg)
}

// Error ...
func (ve *VarError) Error() string {
	return fmt.Sprintf("Error:: Variable name: %s, Message: %s", ve.Name, ve.Msg)
}
