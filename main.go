package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/huboleo/monkey-interpreter/evaluator"
	"github.com/huboleo/monkey-interpreter/lexer"
	"github.com/huboleo/monkey-interpreter/object"
	"github.com/huboleo/monkey-interpreter/parser"
	"github.com/huboleo/monkey-interpreter/repl"
)

func main() {
	switch len(os.Args) {
	case 1:
		startRepl()
	case 2:
		os.Exit(runFile(os.Args[1]))
	default:
		fmt.Fprintln(os.Stderr, "usage: monkey [file.monkey]")
		os.Exit(2)
	}
}

func startRepl() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to monkey REPL %s!", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}

// runFile evaluates a whole Monkey source file and returns the process exit code.
func runFile(path string) int {
	src, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	p := parser.New(lexer.New(string(src)))
	program := p.ParseProgram()
	if errs := p.Errors(); len(errs) != 0 {
		for _, msg := range errs {
			fmt.Fprintf(os.Stderr, "%s: parser error: %s\n", path, msg)
		}
		return 1
	}

	result := evaluator.Eval(program, object.NewEnvironment())
	if errObj, ok := result.(*object.Error); ok {
		fmt.Fprintf(os.Stderr, "%s: %s\n", path, errObj.Inspect())
		return 1
	}

	return 0
}
