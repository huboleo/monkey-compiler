package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/huboleo/monkey-compiler/compiler"
	"github.com/huboleo/monkey-compiler/lexer"
	"github.com/huboleo/monkey-compiler/parser"
	"github.com/huboleo/monkey-compiler/repl"
	"github.com/huboleo/monkey-compiler/vm"
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

// runFile compiles a whole Monkey source file to bytecode, runs it on the VM
// and returns the process exit code.
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

	comp := compiler.New()
	if err := comp.Compile(program); err != nil {
		fmt.Fprintf(os.Stderr, "%s: compiler error: %s\n", path, err)
		return 1
	}

	machine := vm.New(comp.Bytecode())
	if err := machine.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s: vm error: %s\n", path, err)
		return 1
	}

	return 0
}
