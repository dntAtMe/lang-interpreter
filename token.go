package main

import "fmt"

// TokenType represents the type of token
type TokenType int

// Token types enumeration
const (
	// Special tokens
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT  // variable names, function names
	INT    // integers like 123
	STRING // strings like "hello"

	// Operators
	ASSIGN   // =
	PLUS     // +
	MINUS    // -
	MULTIPLY // *
	DIVIDE   // /

	// Comparison operators
	EQ     // ==
	NOT_EQ // !=
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=

	// Logical operators
	AND // &&
	OR  // ||
	NOT // !

	// Delimiters
	COMMA     // ,
	SEMICOLON // ;
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }

	// Keywords
	FUNCTION // func
	LET      // let
	TRUE     // true
	FALSE    // false
	IF       // if
	ELSE     // else
	RETURN   // return
)

// Token represents a single token
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// tokenTypeNames maps token types to their string representations
var tokenTypeNames = map[TokenType]string{
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",

	IDENT:  "IDENT",
	INT:    "INT",
	STRING: "STRING",

	ASSIGN:   "=",
	PLUS:     "+",
	MINUS:    "-",
	MULTIPLY: "*",
	DIVIDE:   "/",

	EQ:     "==",
	NOT_EQ: "!=",
	LT:     "<",
	GT:     ">",
	LTE:    "<=",
	GTE:    ">=",

	AND: "&&",
	OR:  "||",
	NOT: "!",

	COMMA:     ",",
	SEMICOLON: ";",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",

	FUNCTION: "FUNCTION",
	LET:      "LET",
	TRUE:     "TRUE",
	FALSE:    "FALSE",
	IF:       "IF",
	ELSE:     "ELSE",
	RETURN:   "RETURN",
}

// String returns the string representation of a token type
func (tt TokenType) String() string {
	if name, ok := tokenTypeNames[tt]; ok {
		return name
	}
	return fmt.Sprintf("TokenType(%d)", int(tt))
}

// String returns the string representation of a token
func (t Token) String() string {
	return fmt.Sprintf("{Type: %s, Literal: %q, Line: %d, Column: %d}",
		t.Type, t.Literal, t.Line, t.Column)
}

// keywords maps string literals to their corresponding TokenType
var keywords = map[string]TokenType{
	"func":   FUNCTION,
	"let":    LET,
	"true":   TRUE,
	"false":  FALSE,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}

// LookupIdent checks if an identifier is a keyword
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
