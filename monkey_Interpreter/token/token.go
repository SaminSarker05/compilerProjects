package token

type TokenType string

// define possible token types
const (
	ILLEGAL = "ILLEGAL" // unrecognized token
	EOF     = "EOF"

	// literals
	IDENTIFIER = "IDENTIFIER"
	INT        = "INT"

	ASSIGN = "="
	PLUS   = "+"

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
)

var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// use keywords table to check if identifier is a reserved keyword
func LookupIdentifier(identifier string) TokenType {
	if token, ok := keywords[identifier]; ok {
		return token
	}
	return IDENTIFIER
}

type Token struct {
	Type    TokenType
	Literal string
}
