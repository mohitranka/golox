package err

import (
	"fmt"
	"github.com/mohitranka/golox/token"
)

func report(line int, where string, message string) string {
	return fmt.Sprintf("[line %d ] Error %s: %s\n", line, where, message)
}

type RuntimeError struct {
	Line int
	Msg  string
}

type ParseError struct {
	Token token.Token
	Msg   string
}

func (re *RuntimeError) Error() string {
	return report(re.Line, "", re.Msg)
}

func (pe *ParseError) Error() string {
	if pe.Token.Ttype == token.EOF {
		return report(pe.Token.Line, " at end", pe.Msg)
	} else {
		return report(pe.Token.Line, fmt.Sprintf(" at '%s'", pe.Token.Lexeme), pe.Msg)
	}
}
