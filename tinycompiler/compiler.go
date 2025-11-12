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
			fmt.Printf("Unexpected character at pos %d: %s\n", curr_pos, string(char))
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
	Parameters []AstNode  // function definition parameters
	Body       []AstNode  // function definition body
	Arguments  *[]AstNode // function call arguments, CHANGED
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
	if curr_token.Class == "IDENTIFIER" {
		pc++
		return AstNode{
			NodeType: "Identifier",
			Name:     curr_token.Lexeme,
		}
	}
	if curr_token.Class == "LPAREN" {
		pc++
		curr_token = pt[pc]
		node := AstNode{ // start building CallExpression node
			NodeType:   "CallExpression",
			Name:       curr_token.Lexeme,
			Parameters: []AstNode{},
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

// visitor map type for AST traversal, holds functions for each node type
type visitor map[string]func(n *AstNode, parent *AstNode)

// used by Transformer to convert LISP-like AST to C-like AST
func traverser(a ast, v visitor) {
	// entry point for traversal
	traverseNode(AstNode(a), AstNode{}, v)
}

func traverseArray(node []AstNode, parent AstNode, v visitor) {
	for _, child := range node {
		traverseNode(child, parent, v)
	}
}

func traverseNode(node AstNode, parent AstNode, v visitor) {
	for k, va := range v { // call visitor function if node type matches
		if node.NodeType == k {
			va(&node, &parent)
		}
	}
	switch node.NodeType {
	case "Program":
		traverseArray(node.Body, node, v)
	case "CallExpression":
		traverseArray((node.Parameters), node, v)
	case "NumberLiteral":
		break
	case "Identifier":
		break
	default:
		log.Fatalf("Unknown node type during traversal: %s", node.NodeType)
	}
}

// convert LISP-like AST to C-like AST
func transformer(a ast) ast {
	newAst := ast{
		NodeType: "Program",
		Body:     []AstNode{},
	}

	// reference from old AST to new AST
	// specifies where new nodes should be added
	a.Context = &newAst.Body

	traverser(a, map[string]func(n *AstNode, parent *AstNode){
		"NumberLiteral": func(n *AstNode, parent *AstNode) {
			*parent.Context = append(*parent.Context, AstNode{
				NodeType: "NumberLiteral",
				Value:    n.Value,
			})
		},
		"CallExpression": func(n *AstNode, parent *AstNode) {
			expression := AstNode{
				NodeType: "CallExpression",
				Callee: &AstNode{
					NodeType: "Identifier",
					Name:     n.Name,
				},
				Arguments: new([]AstNode),
			}
			n.Context = expression.Arguments
			if parent.NodeType != "CallExpression" {
				es := AstNode{
					NodeType:   "ExpressionStatement",
					Expression: &expression,
				}
				*parent.Context = append(*parent.Context, es)
			} else {
				*parent.Context = append(*parent.Context, expression)
			}
		},
		"Identifier": func(n *AstNode, parent *AstNode) {
			*parent.Context = append(*parent.Context, AstNode{
				NodeType: "Identifier",
				Name:     n.Name,
			})
		},
	})
	return newAst
}

// code generator recursively converts AST to target code
func generator(n AstNode) string {
	switch n.NodeType {
	case "Program":
		var code []string
		for _, child := range n.Body {
			code = append(code, generator(child))
		}
		return strings.Join(code, "\n")
	case "ExpressionStatement":
		return generator(*n.Expression) + ";"
	case "CallExpression":
		// print callee and args in parentheses
		var code []string
		c := generator(*n.Callee)
		for _, arg := range *n.Arguments {
			code = append(code, generator(arg))
		}
		return c + "(" + strings.Join(code, ", ") + ")"
	case "NumberLiteral":
		return n.Value
	case "Identifier":
		return n.Name
	default:
		log.Fatalf("Unknown node type during code generation: %s", n.NodeType)
		return ""
	}
}

// link all compiler steps together
func compiler(input string) string {
	tokens := lexer(input)
	ast := parser(tokens)
	newAst := transformer(ast)
	output := generator(AstNode(newAst))
	return output
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
		if node.Name != "" && len(node.Parameters) > 0 {
			fmt.Printf("%sCallExpression [%s]\n", prefix, node.Name)
			for _, arg := range node.Parameters {
				printAST(arg, indent+1)
			}
		} else if node.Callee != nil { // for transformed AST
			fmt.Printf("%sCallExpression\n", prefix)
			fmt.Printf("%s. Callee:\n", prefix)
			printAST(*node.Callee, indent+2)
			if len(*node.Arguments) > 0 {
				fmt.Printf("%s. Arguments:\n", prefix)
				for _, arg := range *node.Arguments {
					printAST(arg, indent+2)
				}
			}
		}
	case "ExpressionStatement":
		fmt.Printf("%sExpressionStatement\n", prefix)
		printAST(*node.Expression, indent+1)
	case "NumberLiteral":
		fmt.Printf("%sNumberLiteral [%s]\n", prefix, node.Value)
	case "Identifier":
		fmt.Printf("%sIdentifier [%s]\n", prefix, node.Name)
	default:
		fmt.Printf("%sUnknown NodeType: %s\n", prefix, node.NodeType)
	}
}

func main() {
	// test lexer
	fmt.Println("=== Lexer Tests ===")
	testCases := []string{
		"(add 2 (subtract 4 2))",
		"(multiply (add 1 2) (divide 6 3))",
		"(define x 10)",
	}
	for _, input := range testCases {
		fmt.Printf("Input: %s\n", input)
		tokens := lexer(input)
		for _, token := range tokens {
			fmt.Printf("Class: %-10s Lexeme: %s\n", token.Class, token.Lexeme)
		}
		fmt.Println()
	}

	// test parser
	fmt.Println("=== Parser Tests ===")
	for _, input := range testCases {
		fmt.Printf("Input: %s\n", input)
		tokens := lexer(input)
		ast := parser(tokens)
		printAST(AstNode(ast), 0)
		fmt.Println()
	}

	// test transformer and generator
	fmt.Println("=== Transformer and Generator Tests ===")
	for _, input := range testCases {
		fmt.Printf("Input: %s\n", input)
		tokens := lexer(input)
		ast := parser(tokens)
		newAst := transformer(ast)
		fmt.Println("Transformed AST:")
		printAST(AstNode(newAst), 0)
		fmt.Println()
		output := generator(AstNode(newAst))
		fmt.Printf("Generated Code:\n%s\n\n", output)
	}
}
