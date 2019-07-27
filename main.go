package main

import (
	"bufio"
	"fmt"
	"github.com/mohitranka/golox/scanner"
	"io/ioutil"
	"os"
	"strings"
)

type Lox struct {
	HadError bool
}

func NewLox() *Lox {
	l := new(Lox)
	l.HadError = false
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
	scanner := scanner.NewScanner(source)
	tokens := scanner.ScanTokens()
	for _, token := range tokens {
		fmt.Println(token)
	}
}
