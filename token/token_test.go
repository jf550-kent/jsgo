package token

import (
	"testing"
)

func TestTokenString(t *testing.T) {
	tests := []struct {
		token    TokenType
		expected string
	}{
		{ILLEGAL, "ILLEGAL"},
		{EOF, "EOF"},
		{IDENT, "IDENTIFIER"},
		{NUMBER, "NUMBER"},
		{STRING, "STRING"},
		{FLOAT, "FLOAT"},
		{ADD, "+"},
		{MINUS, "-"},
		{MUL, "*"},
		{DIVIDE, "/"},
		{LSS, "<"},
		{GTR, ">"},
		{BANG, "!"},
		{ASSIGN, "="},
		{NOT_EQUAL, "!="},
		{EQUAL, "=="},
		{COMMA, ","},
		{SEMICOLON, ";"},
		{DOT, "."},
		{COLON, ":"},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{FUNCTION, "function"},
		{VAR, "var"},
		{IF, "if"},
		{ELSE, "else"},
		{ELSEIF, "elseif"},
		{RETURN, "return"},
		{TRUE, "true"},
		{FALSE, "false"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{NULL, "null"},
	}

	for _, tt := range tests {
		token := Token{TokenType: tt.token}
		if token.String() != tt.expected {
			t.Errorf("expected %s, got %s", tt.expected, token.String())
		}
	}
}

func TestIsLiteral(t *testing.T) {
	literals := []TokenType{IDENT, NUMBER, STRING, FLOAT}
	for _, tt := range literals {
		token := Token{TokenType: tt}
		if !token.IsLiteral() {
			t.Errorf("expected true for literal %v", token)
		}
	}

	nonLiterals := []TokenType{ADD, FUNCTION}
	for _, tt := range nonLiterals {
		token := Token{TokenType: tt}
		if token.IsLiteral() {
			t.Errorf("expected false for non-literal %v", token)
		}
	}
}

func TestIsOperator(t *testing.T) {
	operators := []TokenType{ADD, MINUS, MUL, DIVIDE, LSS, GTR, BANG, ASSIGN, NOT_EQUAL, EQUAL, LBRACKET, RBRACKET}
	for _, tt := range operators {
		token := Token{TokenType: tt}
		if !token.IsOperator() {
			t.Errorf("expected true for operator %v", token)
		}
	}

	nonOperators := []TokenType{IDENT, FUNCTION}
	for _, tt := range nonOperators {
		token := Token{TokenType: tt}
		if token.IsOperator() {
			t.Errorf("expected false for non-operator %v", token)
		}
	}
}

func TestIsKeyword(t *testing.T) {
	keywords := []string{"function", "var", "if", "else", "elseif", "return", "true", "false", "for", "null"}

	for _, keyword := range keywords {
		if !IsKeyword(keyword) {
			t.Errorf("expected true for keyword %s", keyword)
		}
	}

	nonKeywords := []string{"apple", "banana", "IDENT", "NUMBER"}

	for _, nonKeyword := range nonKeywords {
		if IsKeyword(nonKeyword) {
			t.Errorf("expected false for non-keyword %s", nonKeyword)
		}
	}
}

func TestIsIdentifier(t *testing.T) {
	identifiers := []string{"apple", "banana", "myVar"}

	for _, identifier := range identifiers {
		if !IsIdentifier(identifier) {
			t.Errorf("expected true for identifier %s", identifier)
		}
	}

	nonIdentifiers := []string{"", "var", "if", "null"}

	for _, nonIdentifier := range nonIdentifiers {
		if IsIdentifier(nonIdentifier) {
			t.Errorf("expected false for non-identifier %s", nonIdentifier)
		}
	}
}
