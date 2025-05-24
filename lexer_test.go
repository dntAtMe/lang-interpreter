package main

import (
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = func(x, y) {
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
"foobar"
"foo bar"
// This is a comment
let x = 42; // Another comment`

	tests := []struct {
		expectedType    TokenType
		expectedLiteral string
	}{
		{LET, "let"},
		{IDENT, "five"},
		{ASSIGN, "="},
		{INT, "5"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "ten"},
		{ASSIGN, "="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "add"},
		{ASSIGN, "="},
		{FUNCTION, "func"},
		{LPAREN, "("},
		{IDENT, "x"},
		{COMMA, ","},
		{IDENT, "y"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{IDENT, "x"},
		{PLUS, "+"},
		{IDENT, "y"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{SEMICOLON, ";"},
		{LET, "let"},
		{IDENT, "result"},
		{ASSIGN, "="},
		{IDENT, "add"},
		{LPAREN, "("},
		{IDENT, "five"},
		{COMMA, ","},
		{IDENT, "ten"},
		{RPAREN, ")"},
		{SEMICOLON, ";"},
		{NOT, "!"},
		{MINUS, "-"},
		{DIVIDE, "/"},
		{MULTIPLY, "*"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{GT, ">"},
		{INT, "5"},
		{SEMICOLON, ";"},
		{IF, "if"},
		{LPAREN, "("},
		{INT, "5"},
		{LT, "<"},
		{INT, "10"},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{TRUE, "true"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{ELSE, "else"},
		{LBRACE, "{"},
		{RETURN, "return"},
		{FALSE, "false"},
		{SEMICOLON, ";"},
		{RBRACE, "}"},
		{INT, "10"},
		{EQ, "=="},
		{INT, "10"},
		{SEMICOLON, ";"},
		{INT, "10"},
		{NOT_EQ, "!="},
		{INT, "9"},
		{SEMICOLON, ";"},
		{STRING, "foobar"},
		{STRING, "foo bar"},
		{LET, "let"},
		{IDENT, "x"},
		{ASSIGN, "="},
		{INT, "42"},
		{SEMICOLON, ";"},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifierRecognition(t *testing.T) {
	tests := []struct {
		input    string
		expected TokenType
	}{
		{"let", LET},
		{"func", FUNCTION},
		{"if", IF},
		{"else", ELSE},
		{"return", RETURN},
		{"true", TRUE},
		{"false", FALSE},
		{"myVariable", IDENT},
		{"test123", IDENT},
		{"_private", IDENT},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != tt.expected {
			t.Errorf("For input %q: expected %s, got %s",
				tt.input, tt.expected, tok.Type)
		}
	}
}

func TestLineAndColumnTracking(t *testing.T) {
	input := `let x = 5;
let y = 10;`

	l := New(input)

	// First line tokens
	tok := l.NextToken()
	if tok.Line != 1 || tok.Column != 1 {
		t.Errorf("Expected line 1, column 1, got line %d, column %d", tok.Line, tok.Column)
	}

	l.NextToken()
	l.NextToken()
	l.NextToken()
	l.NextToken()

	// Second line tokens
	tok = l.NextToken()
	if tok.Line != 2 || tok.Column != 1 {
		t.Errorf("Expected line 2, column 1, got line %d, column %d", tok.Line, tok.Column)
	}
}

func TestCommentSkipping(t *testing.T) {
	input := `let x = 5; // This is a comment
// Full line comment
let y = 10;`

	l := New(input)

	tokens := []TokenType{LET, IDENT, ASSIGN, INT, SEMICOLON, LET, IDENT, ASSIGN, INT, SEMICOLON, EOF}

	for i, expectedType := range tokens {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Errorf("Token %d: expected %s, got %s", i, expectedType, tok.Type)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello"`, "hello"},
		{`"hello world"`, "hello world"},
		{`""`, ""},
		{`"with spaces and 123 numbers"`, "with spaces and 123 numbers"},
	}

	for _, tt := range tests {
		l := New(tt.input)
		tok := l.NextToken()

		if tok.Type != STRING {
			t.Errorf("Expected STRING token, got %s", tok.Type)
		}

		if tok.Literal != tt.expected {
			t.Errorf("Expected literal %q, got %q", tt.expected, tok.Literal)
		}
	}
}

func TestOperators(t *testing.T) {
	input := `= + - * / == != < > && || !`

	expected := []TokenType{
		ASSIGN, PLUS, MINUS, MULTIPLY, DIVIDE,
		EQ, NOT_EQ, LT, GT, AND, OR, NOT, EOF,
	}

	l := New(input)

	for i, expectedType := range expected {
		tok := l.NextToken()
		if tok.Type != expectedType {
			t.Errorf("Token %d: expected %s, got %s", i, expectedType, tok.Type)
		}
	}
}
