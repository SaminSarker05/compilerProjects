package repl

import (
	"bufio"
	"fmt"
	"gorilla_interpreter/lexer"
	"gorilla_interpreter/token"
	"io"
)

func StartRepl(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Printf(">> ")
		userInput := scanner.Scan()

		if !userInput {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Printf("%+v\n", tok)
		}
	}
}
