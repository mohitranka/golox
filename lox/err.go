package lox

import (
	"fmt"
)

func report(line int, where string, message string) string {
	return fmt.Sprintf("[line %d ] Error %s: %s\n", line, where, message)
}

type RuntimeError struct {
	Line int
	Msg  string
}

type ParseError struct {
	Token Token
	Msg   string
}

type VarError struct {
	Name string
	Msg  string
}

func (re *RuntimeError) Error() string {
	return report(re.Line, "", re.Msg)
}

func (pe *ParseError) Error() string {
	if pe.Token.Type == EOF {
		return report(pe.Token.Line, " at end", pe.Msg)
	} else {
		return report(pe.Token.Line, fmt.Sprintf(" at '%s'", pe.Token.Lexeme), pe.Msg)
	}
}

func (ve *VarError) Error() string {
	return fmt.Sprintf("Error:: Variable name: %s, Message: %s", ve.Name, ve.Msg)
}