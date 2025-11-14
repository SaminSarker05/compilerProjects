package token

type TokenType string

// define possible token types
const (
	ILLEGAL = "ILLEGAL" // unrecognized token
	EOF     = "EOF"

	// literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	// operators
	ASSIGN = "="
	PLUS   = "+"

	// double char operators
	EQUAL     = "=="
	NOT_EQUAL = "!="

	// operators
	MINUS       = "-"
	BANG        = "!"
	ASTERISK    = "*"
	SLASH       = "/"
	LESSTHAN    = "<"
	GREATERTHAN = ">"

	// delimiters
	COMMA     = ","
	SEMICOLON = ";"
	LPAREN    = "("
	RPAREN    = ")"
	LBRACE    = "{"
	RBRACE    = "}"

	// keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE     = "TRUE"
	FALSE    = "FALSE"
	RETURN   = "RETURN"
	IF       = "IF"
	ELSE     = "ELSE"
)

var keywords = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"false":  FALSE,
	"true":   TRUE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

type Token struct {
	Type    TokenType
	Literal string
}

// use keywords table to check if identifier is a reserved keyword
func LookupIdentifier(identifier string) TokenType {
	if tokentype, ok := keywords[identifier]; ok {
		return tokentype
	}
	return IDENTIFIER
}
