package lexer

import "monkey_interpreter/token"

type Lexer struct {
	input        string
	position     int  // current pos in input where last read char
	readPosition int  // current read pos for lookahead, always ahead of position
	char         byte // char under examination, bytes instead of rune because unicode more complex
}

// constructor for lexer struct
func New(input string) *Lexer {
	l := &Lexer{input: input} // create new lexer on heap
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// return next token from lexer and advance position
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	switch l.char {
	case '=':
		tok = newToken(token.ASSIGN, l.char)
	case ';':
		tok = newToken(token.SEMICOLON, l.char)
	case '(':
		tok = newToken(token.LPAREN, l.char)
	case ')':
		tok = newToken(token.RPAREN, l.char)
	case '{':
		tok = newToken(token.LBRACE, l.char)
	case '}':
		tok = newToken(token.RBRACE, l.char)
	case '+':
		tok = newToken(token.PLUS, l.char)
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar() // read next character
	return tok
}

// creates a new token of the given token class and lexeme
func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

// check if byte char is a letter
func isLetter(char byte) bool {
	// allow _, ?, and ! in identifiers
	return ('a' <= char && char <= 'z') || ('A' <= char && char <= 'Z') || char == '_' || char == '!' || char == '?'
}

// read chars to build identifier until non-letter char encountered
func (l *Lexer) readIdentifier() string {
	start_pos := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[start_pos:l.position]
}
