package main

import (
	"fmt"
	"gorilla_interpreter/repl"
	"os"
)

func main() {
	fmt.Printf("Type in commands\n")
	repl.StartRepl(os.Stdin, os.Stdout)
}
