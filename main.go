package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/huboleo/monkey-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Welcome to monkey REPL %s!", user.Username)
	repl.Start(os.Stdin, os.Stdout)
}
