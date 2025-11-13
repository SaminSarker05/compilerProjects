package lexer

import (
	"testing"

	"monkey_interpreter/token"
)

func TestNextToken(t *testing.T) { // t holds test state
	input := `let five = 5;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
	}

	lexer := New(input)

	for pos, test := range tests {
		token := lexer.NextToken()

		if token.Type != test.expectedType {
			t.Fatalf("index[%d] - wrong token type. expected=%q, got=%q",
				pos, test.expectedType, token.Type)
		}
		if token.Literal != test.expectedLiteral {
			t.Fatalf("index[%d] - wrong token literal, expected=%q, got=%q",
				pos, test.expectedLiteral, token.Literal)
		}
	}
}
