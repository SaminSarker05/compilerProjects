package lexer

import (
	"gorilla_interpreter/token"
)

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

// consumes next char and update position fields
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.char = 0 // EOF
	} else {
		l.char = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

// returns next token from lexer and advance position
func (l *Lexer) NextToken() token.Token {
	var tok token.Token
	l.skipWhiteSpaces() // eat whitespace characters

	switch l.char {
	case '=':
		// peak ahead to check if == token or = token
		if l.peek() == '=' {
			tok = token.Token{
				Type:    token.EQUAL,
				Literal: string(l.char) + string(l.peek()),
			}
			l.readChar()
		} else {
			tok = newToken(token.ASSIGN, l.char)
		}
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
	case ',':
		tok = newToken(token.COMMA, l.char)
	case '-':
		tok = newToken(token.MINUS, l.char)
	case '/':
		tok = newToken(token.SLASH, l.char)
	case '<':
		tok = newToken(token.LESSTHAN, l.char)
	case '>':
		tok = newToken(token.GREATERTHAN, l.char)
	case '*':
		tok = newToken(token.ASTERISK, l.char)
	case '!':
		// peak ahead to check if != token or ! token
		if l.peek() == '=' {
			tok = token.Token{
				Type:    token.NOT_EQUAL,
				Literal: string(l.char) + string(l.peek()),
			}
			l.readChar()
		} else {
			tok = newToken(token.BANG, l.char)
		}
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.char) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.char) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.char)
		}
	}

	l.readChar() // read next character
	return tok
}

// peeks ahead at next char without eating
func (l *Lexer) peek() byte {
	if l.readPosition >= len(l.input) {
		return 0 // EOF
	} else {
		return l.input[l.readPosition]
	}
}

// creates a new token of the given token class and lexeme
func newToken(tokenType token.TokenType, char byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(char)}
}

// reads chars to build identifier until non-letter encountered
func (l *Lexer) readIdentifier() string {
	start_pos := l.position
	for isLetter(l.char) {
		l.readChar()
	}
	return l.input[start_pos:l.position]
}

func (l *Lexer) skipWhiteSpaces() {
	for l.char == ' ' || l.char == '\t' || l.char == '\n' || l.char == '\r' {
		l.readChar()
	}
}

// reads chars to build number
func (l *Lexer) readNumber() string {
	start_pos := l.position
	for isDigit(l.char) {
		l.readChar()
	}
	return l.input[start_pos:l.position]
}
