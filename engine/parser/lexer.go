package parser

import (
	"strings"
	"unicode"
)

// Lexer implementation (basic tokenizer)
// Supports keywords, identifiers, numbers, strings, symbols

type TokenType string

const (
	TokEOF     TokenType = "EOF"
	TokIdent   TokenType = "IDENT"
	TokNumber  TokenType = "NUMBER"
	TokString  TokenType = "STRING"
	TokComma   TokenType = ","
	TokLParen  TokenType = "("
	TokRParen  TokenType = ")"
	TokStar    TokenType = "*"
	TokEqual   TokenType = "="
	TokKeyword TokenType = "KEYWORD"
)

type Token struct {
	Type  TokenType
	Value string
}

type Lexer struct {
	input []rune
	pos   int
}

func NewLexer(s string) *Lexer {
	return &Lexer{input: []rune(s)}
}

func (l *Lexer) next() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	r := l.input[l.pos]
	l.pos++
	return r
}

func (l *Lexer) peek() rune {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) skipSpaces() {
	for {
		if ch := l.peek(); ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			l.next()
			continue
		}
		break
	}
}

func (l *Lexer) readIdent() string {
	var out []rune
	for {
		ch := l.peek()
		if ch == 0 || !(unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_') {
			break
		}
		out = append(out, l.next())
	}
	return string(out)
}

func (l *Lexer) readNumber() string {
	var out []rune
	for unicode.IsDigit(l.peek()) {
		out = append(out, l.next())
	}
	return string(out)
}

func (l *Lexer) readString() string {
	l.next() // skip initial quote
	var out []rune
	for {
		ch := l.next()
		if ch == '\'' || ch == 0 {
			break
		}
		out = append(out, ch)
	}
	return string(out)
}

func (l *Lexer) NextToken() Token {
	l.skipSpaces()
	ch := l.peek()

	switch {
	case ch == 0:
		return Token{Type: TokEOF}
	case ch == ',':
		l.next()
		return Token{Type: TokComma, Value: ","}
	case ch == '(':
		l.next()
		return Token{Type: TokLParen, Value: "("}
	case ch == ')':
		l.next()
		return Token{Type: TokRParen, Value: ")"}
	case ch == '*':
		l.next()
		return Token{Type: TokStar, Value: "*"}
	case ch == '=':
		l.next()
		return Token{Type: TokEqual, Value: "="}
	case unicode.IsLetter(ch):
		ident := l.readIdent()
		upper := strings.ToUpper(ident)
		switch upper {
		case "SELECT", "INSERT", "INTO", "VALUES", "CREATE", "TABLE", "WHERE", "SET", "FROM", "UPDATE", "DELETE":
			return Token{Type: TokKeyword, Value: upper}
		default:
			return Token{Type: TokIdent, Value: ident}
		}
	case unicode.IsDigit(ch):
		num := l.readNumber()
		return Token{Type: TokNumber, Value: num}
	case ch == '\'':
		str := l.readString()
		return Token{Type: TokString, Value: str}
	}

	l.next()
	return Token{Type: TokIdent, Value: string(ch)}
}

// Parse is a convenience wrapper that uses the lexer and parser to parse SQL
func Parse(sql string) (Statement, error) {
	l := NewLexer(sql)
	p := NewParser(l)
	return p.ParseStatement()
}
