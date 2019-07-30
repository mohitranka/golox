package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	//Local packages
	"github.com/mohitranka/golox/interpreter"
	"github.com/mohitranka/golox/parser"
	"github.com/mohitranka/golox/scanner"
)

type Lox struct {
	HadError        bool
	HadRuntimeError bool
}

func NewLox() *Lox {
	l := new(Lox)
	l.HadError = false
	l.HadRuntimeError = false
	return l
}

func main() {
	if len(os.Args) > 2 {
		fmt.Println("Usage: golox [script]")
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
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		l.run(text)
		l.HadError = false
	}
}

func (l Lox) run(source string) {
	s := scanner.NewScanner(source)
	tokens := s.ScanTokens()
	p := parser.NewParser(tokens)
	stmts := p.Parse()
	if l.HadError {
		return
	}
	i := &interpreter.Interpreter{}
	i.Interpret(stmts)
}
