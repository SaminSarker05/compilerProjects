package lexer

import (
	"testing"

	"gorilla_interpreter/token"
)

func TestNextTokenAll(t *testing.T) { // t holds test state
	input := `let five = 5;
		let ten = 10;
		let add = fn(x, y) {
		x + y;
		};
		let result = add(five, ten);
		!-/*5;
		5 < 10 > 5;

		if (5 < 10) {
			return true;
		} else {
			return false; 
		}

		10 == 10;
		10 != 9;
	`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.INT, "5"},
		{token.LESSTHAN, "<"},
		{token.INT, "10"},
		{token.GREATERTHAN, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LESSTHAN, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.INT, "10"},
		{token.EQUAL, "=="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},
		{token.INT, "10"},
		{token.NOT_EQUAL, "!="},
		{token.INT, "9"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
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

func TestNextTokenSingleCharLiterals(t *testing.T) {
	input := `=+(){},;`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
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

func TestNextTokenDoubleCharLiterals(t *testing.T) {
	input := `!=<>==`
	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.NOT_EQUAL, "!="},
		{token.LESSTHAN, "<"},
		{token.GREATERTHAN, ">"},
		{token.EQUAL, "=="},
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
