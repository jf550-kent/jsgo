package token

import (
	"strconv"
)

type Pos struct {
	Line int
	Col  int
}

// Token is samllest valid element of the language
type Token struct {
	TokenType TokenType
	Literal   string
	Start     Pos
	End       Pos
}

// TokenType represent the set of valid token type in the language
type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	literalBegin

	IDENT  // apple
	NUMBER // 89
	STRING // "number"
	FLOAT  // 81.0

	literalEnd

	operatorBegin // operator and delimiter

	ADD    // +
	MINUS  // -
	MUL    // *
	DIVIDE // /

	LSS    // <
	GTR    // >
	BANG   // !
	ASSIGN // =

	NOT_EQUAL // !=
	EQUAL     // ==

	COMMA     // ,
	SEMICOLON // ;
	DOT       // .
	COLON     // :

	LPAREN   // (
	RPAREN   // )
	LBRACE   // {
	RBRACE   // }
	LBRACKET // [
	RBRACKET // ]

	operatorEnd

	keywordBegin // keyword in the language of jsgo

	FUNCTION // function
	VAR      // var
	IF       // if
	ELSE     // else
	ELSEIF   // elseif
	RETURN   // return
	TRUE     // true
	FALSE    // false
	FOR      // for

	keywordEnd
)

// keywords maps the only valid keywords in the language
var keywords = map[string]TokenType{
	"var":      VAR,
	"function": FUNCTION,
	"if":       IF,
	"else":     ELSE,
	"elseif":   ELSEIF,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
	"for":      FOR,
}

// tokens store the repective string representation of the token
var tokens = [...]string{
	ILLEGAL:   "ILLEGAL",
	EOF:       "EOF",
	IDENT:     "IDENTIFIER",
	NUMBER:    "NUMBER",
	STRING:    "STRING",
	FLOAT:     "FLOAT",
	ADD:       "+",
	MINUS:     "-",
	MUL:       "*",
	DIVIDE:    "/",
	LSS:       "<",
	GTR:       ">",
	BANG:      "!",
	ASSIGN:    "=",
	NOT_EQUAL: "!=",
	EQUAL:     "==",
	COMMA:     ",",
	SEMICOLON: ";",
	DOT:       ".",
	COLON:     ":",
	LPAREN:    "(",
	RPAREN:    ")",
	LBRACE:    "{",
	RBRACE:    "}",
	LBRACKET:  "[",
	RBRACKET:  "]",
	FUNCTION:  "function",
	VAR:       "var",
	IF:        "if",
	ELSE:      "else",
	ELSEIF:    "elseif",
	RETURN:    "return",
	TRUE:      "true",
	FALSE:     "false",
	FOR:       "for",
}

func (t Token) Precedence() int {
	switch t.TokenType {
	case EQUAL, NOT_EQUAL:
		return 2
	case LSS, GTR:
		return 3
	case ADD, MINUS:
		return 4
	case DIVIDE, MUL:
		return 5
	case LPAREN:
		return 7
	}
	return 1
}

// String return the string representation of the token
func (t Token) String() string {
	s := ""

	if t.TokenType >= 0 && int(t.TokenType) < len(tokens) {
		s = tokens[t.TokenType]
	}

	if s == "" {
		s = "token(" + strconv.Itoa(int(t.TokenType)) + ") does not exist"
	}
	return s
}

func (t TokenType) String() string {
	s := ""

	if t >= 0 && int(t) < len(tokens) {
		s = tokens[t]
	}

	if s == "" {
		s = "token(" + strconv.Itoa(int(t)) + ") does not exist"
	}
	return s
}

// IsLiteral determine if the token is a literal in the language
func (t Token) IsLiteral() bool { return literalBegin < t.TokenType && t.TokenType < literalEnd }

// IsOperator determine if the token is an operator in the language
func (t Token) IsOperator() bool { return operatorBegin < t.TokenType && t.TokenType < operatorEnd }

// IsKeyword determine if the token is a keyword
func (t Token) IsKeyword() bool {
	s := tokens[t.TokenType]
	if s != t.Literal {
		panic("Token literal and tokenType mismatch")
	}

	return IsKeyword(s) && keywordBegin < t.TokenType && t.TokenType < keywordEnd
}

func Keyword(word string) (TokenType, bool) {
	ty, ok := keywords[word]
	return ty, ok
}

// IsKeyword determine if the string pass into it is part of a keyword in the language
func IsKeyword(word string) bool {
	_, ok := keywords[word]
	return ok
}

// IsIdentifier determine if the string pass into it is an identifier
func IsIdentifier(word string) bool {
	if len(word) == 0 || IsKeyword(word) {
		return false
	}
	return true
}
