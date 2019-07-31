package lox

import (
	"fmt"
	"io/ioutil"
	"os"

	//Third party
	"github.com/chzyer/readline"
)

type Lox struct {
	HadError        bool
	HadRuntimeError bool
	Interpreter     *Interpreter
}

func NewLox() *Lox {
	l := new(Lox)
	l.HadError = false
	l.HadRuntimeError = false
	l.Interpreter = NewInterpreter()
	return l
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: go[script]")
		os.Exit(64)
	} else {
		l := NewLox()
		if len(os.Args) == 2 {
			l.runFile(os.Args[1])
		} else {
			l.runPrompt()
		}
	}
}

func (l Lox) runFile(path string) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println(err)
		l.HadError = true
	}

	l.run(string(bytes))

	if l.HadError {
		os.Exit(65)
	}

	if l.HadRuntimeError {
		os.Exit(70)
	}
}

func (l Lox) runPrompt() {
	rl, e := readline.New("> ")
	if e != nil {
		panic(e)
	}
	defer rl.Close()

	for {
		line, e := rl.Readline()
		if e != nil {
			break
		}
		if line == "exit" {
			os.Exit(0)
		}
		l.run(line)
		l.HadError = false
	}
}

func (l Lox) run(source string) {
	s := NewScanner(source)
	tokens := s.ScanTokens()
	p := NewParser(tokens)
	stmts := p.Parse()
	if l.HadError {
		return
	}
	l.Interpreter.Interpret(stmts)
}
