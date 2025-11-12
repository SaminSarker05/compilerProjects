package main

import (
	"fmt"
	"log"
	"strings"
	"unicode"
)

type Token struct {
	Class  string
	Lexeme string
}

func lexer(input string) []Token {
	tokens := []Token{}

	input += "\n"          // sentinel to mark end of input
	curr_pos := 0          // track current position in input string
	runes := []rune(input) // convert input string to rune slice for proper char indexing

	for curr_pos < len(runes) {
		char := runes[curr_pos]

		// skip any whitespace characters
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			curr_pos++
		} else if char == '(' {
			tokens = append(tokens, Token{
				Class:  "LPAREN",
				Lexeme: "(",
			})
			curr_pos++
		} else if char == ')' {
			tokens = append(tokens, Token{
				Class:  "RPAREN",
				Lexeme: ")",
			})
			curr_pos++
		} else if unicode.IsDigit(char) {
			start_pos := curr_pos
			for curr_pos < len(runes) && unicode.IsDigit(runes[curr_pos]) {
				curr_pos++
			}
			number := string(runes[start_pos:curr_pos])
			tokens = append(tokens, Token{
				Class:  "NUMBER",
				Lexeme: number,
			})
		} else if unicode.IsLetter(char) {
			start_pos := curr_pos
			for curr_pos < len(runes) && unicode.IsLetter(runes[curr_pos]) {
				curr_pos++
			}
			identifier := string(runes[start_pos:curr_pos])
			tokens = append(tokens, Token{
				Class:  "IDENTIFIER",
				Lexeme: identifier,
			})
		} else {
			fmt.Printf("Unexpected character: %s\n", char)
			break
		}
	}
	return tokens
}

type AstNode struct {
	NodeType   string     // what node represents
	Value      string     // literal values if any
	Name       string     // function/variable names
	Callee     *AstNode   // call target for higher order functions
	Expression *AstNode   // expressions
	Arguments  []AstNode  // function call arguments, CHANGED
	Parameters []AstNode  // function definition parameters
	Body       []AstNode  // function definition body
	Context    *[]AstNode // used in traversal for execution context

}

// type alias for root of AST
type ast AstNode

var pc int     // parser counter, track current position in tokens slice
var pt []Token // parser tokens, slice of tokens being parsed

func parser(tokens []Token) ast {
	pc = 0
	pt = tokens
	// create root of AST
	ast := ast{
		NodeType: "Program",
		Body:     []AstNode{},
	}
	// call recursive walk functio to build AST
	for pc < len(pt) {
		ast.Body = append(ast.Body, walk())
	}
	return ast
}

// recursive descent parsing function
func walk() AstNode {
	curr_token := pt[pc]
	if curr_token.Class == "NUMBER" {
		pc++
		return AstNode{
			NodeType: "NumberLiteral",
			Value:    curr_token.Lexeme,
		}
	}
	if curr_token.Class == "LPAREN" {
		pc++
		curr_token = pt[pc]
		node := AstNode{ // start building CallExpression node
			NodeType:  "CallExpression",
			Name:      curr_token.Lexeme,
			Arguments: []AstNode{},
		}
		pc++
		curr_token = pt[pc]

		for curr_token.Class != "RPAREN" {
			node.Parameters = append(node.Parameters, walk())
			curr_token = pt[pc]
		}
		pc++
		return node
	}
	log.Fatal("Unexpected token type: " + curr_token.Class)
	return AstNode{}
}

// helper function to print AST in readable format
func printAST(node AstNode, indent int) {
	prefix := strings.Repeat("  ", indent)
	switch node.NodeType {
	case "Program":
		fmt.Printf("%sProgram\n", prefix)
		for _, child := range node.Body {
			printAST(child, indent+1)
		}
	case "CallExpression":
		fmt.Printf("%sCallExpression [%s]\n", prefix, node.Name)
		for _, arg := range node.Arguments {
			printAST(arg, indent+1)
		}
	case "NumberLiteral":
		fmt.Printf("%sNumberLiteral [%s]\n", prefix, node.Value)
	default:
		fmt.Printf("%sUnknown NodeType: %s\n", prefix, node.NodeType)
	}
}

func main() {
	// test lexer
	testCasesLexer := []string{
		"(add 2 (subtract 4 2))",
		"(multiply (add 1 2) (divide 6 3))",
		"(define x 10)",
	}
	for _, input := range testCasesLexer {
		fmt.Printf("Input: %s\n", input)
		tokens := lexer(input)
		for _, token := range tokens {
			fmt.Printf("Class: %-10s Lexeme: %s\n", token.Class, token.Lexeme)
		}
		fmt.Println()
	}

	// test parser
	testCasesParser := []string{
		"(add 2 (subtract 4 2))",
		"(multiply (add 1 2) (divide 6 3))",
		"(define x 10)",
	}
	for _, input := range testCasesParser {
		fmt.Printf("Input: %s\n", input)
		tokens := lexer(input)
		ast := parser(tokens)
		printAST(AstNode(ast), 0)
	}
}
