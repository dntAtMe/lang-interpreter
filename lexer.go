package main

// Lexer represents the lexical analyzer
type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte
	line         int
	column       int
}

// New creates a new lexer instance
func New(input string) *Lexer {
	l := &Lexer{
		input:  input,
		line:   1,
		column: 0,
	}
	l.readChar()
	return l
}

// readChar reads the next character and advances the position
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++

	// Track line and column numbers
	if l.ch == '\n' {
		l.line++
		l.column = 0
	} else {
		l.column++
	}
}

// peekChar returns the next character without advancing the position
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken scans the input and returns the next token
func (l *Lexer) NextToken() Token {
	var tok Token

	// Skip whitespace (except newlines, which we track for line numbers)
	l.skipWhitespace()

	// Store current position for token
	line := l.line
	column := l.column

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(ASSIGN, l.ch, line, column)
		}
	case '+':
		tok = newToken(PLUS, l.ch, line, column)
	case '-':
		tok = newToken(MINUS, l.ch, line, column)
	case '*':
		tok = newToken(MULTIPLY, l.ch, line, column)
	case '/':
		// Check for single-line comments
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken()
		}
		tok = newToken(DIVIDE, l.ch, line, column)
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NOT_EQ, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(NOT, l.ch, line, column)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(LT, l.ch, line, column)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(GT, l.ch, line, column)
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: AND, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	case '|':
		if l.peekChar() == '|' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: OR, Literal: string(ch) + string(l.ch), Line: line, Column: column}
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	case ',':
		tok = newToken(COMMA, l.ch, line, column)
	case ';':
		tok = newToken(SEMICOLON, l.ch, line, column)
	case '(':
		tok = newToken(LPAREN, l.ch, line, column)
	case ')':
		tok = newToken(RPAREN, l.ch, line, column)
	case '{':
		tok = newToken(LBRACE, l.ch, line, column)
	case '}':
		tok = newToken(RBRACE, l.ch, line, column)
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
		tok.Line = line
		tok.Column = column
	case 0:
		tok.Literal = ""
		tok.Type = EOF
		tok.Line = line
		tok.Column = column
	default:
		if isLetter(l.ch) {
			tok.Line = line
			tok.Column = column
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			tok.Line = line
			tok.Column = column
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch, line, column)
		}
	}

	l.readChar()
	return tok
}

// newToken creates a new token with the given type and character
func newToken(tokenType TokenType, ch byte, line, column int) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: line, Column: column}
}

// readIdentifier reads an identifier or keyword
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber reads a number
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString reads a string literal
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// skipWhitespace skips whitespace characters except newlines
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// skipComment skips single-line comments (// comment)
func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// isLetter checks if a character is a letter or underscore
func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

// isDigit checks if a character is a digit
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// TokenizeAll returns all tokens from the input as a slice
func (l *Lexer) TokenizeAll() []Token {
	var tokens []Token

	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	return tokens
}
