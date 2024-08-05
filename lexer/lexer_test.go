package lexer

import (
	"os"
	"testing"

	"github.com/jf550-kent/jsgo/benchmark"
	"github.com/jf550-kent/jsgo/token"
)

func BenchmarkLex(b *testing.B) {
	byt, err := os.ReadFile(benchmark.EXAMPLE_FILE_PATH)
	if err != nil {
		b.Fatal("failed to read file", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l := New(byt)
		tok, err := l.Lex()
		if err != nil {
			b.Fatal("Lexer.Lex: error in Lex", err)
		}
		if tok.TokenType == token.EOF {
			break
		}
	}
}

func TestLexerLineAndCol(t *testing.T) {
	input := `How are you;
	
	hh`
	expected := []struct {
		char byte
		line int
		col  int
	}{
		{'H', 1, 1},
		{'o', 1, 2},
		{'w', 1, 3},
		{' ', 1, 4},
		{'a', 1, 5},
		{'r', 1, 6},
		{'e', 1, 7},
		{' ', 1, 8},
		{'y', 1, 9},
		{'o', 1, 10},
		{'u', 1, 11},
		{';', 1, 12},
		{'\n', 2, 0},
		{'\t', 2, 1},
		{'\n', 3, 0},
		{'\t', 3, 1},
		{'h', 3, 2},
		{'h', 3, 3},
	}

	l := New([]byte(input))

	for _, exp := range expected {

		if l.ch != exp.char {
			t.Errorf("Expected character '%c', got '%c'", exp.char, l.ch)
		}
		if l.line != exp.line {
			t.Errorf("Expected line %d, got %d", exp.line, l.line)
		}
		if l.col != exp.col {
			t.Errorf("Expected column %d, got %d", exp.col, l.col)
		}

		l.next() // Advance lexer to the next character
	}
}

func TestLexSingleToken(t *testing.T) {
	tests := []struct {
		input         string
		expectedToken token.Token
	}{
		{"+", token.Token{TokenType: token.ADD, Literal: "+", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"-", token.Token{TokenType: token.MINUS, Literal: "-", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"*", token.Token{TokenType: token.MUL, Literal: "*", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"/", token.Token{TokenType: token.DIVIDE, Literal: "/", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{",", token.Token{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{".", token.Token{TokenType: token.DOT, Literal: ".", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{":", token.Token{TokenType: token.COLON, Literal: ":", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{";", token.Token{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"(", token.Token{TokenType: token.LPAREN, Literal: "(", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{")", token.Token{TokenType: token.RPAREN, Literal: ")", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"{", token.Token{TokenType: token.LBRACE, Literal: "{", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"}", token.Token{TokenType: token.RBRACE, Literal: "}", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"=", token.Token{TokenType: token.ASSIGN, Literal: "=", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"!", token.Token{TokenType: token.BANG, Literal: "!", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 1}}},
		{"!=", token.Token{TokenType: token.NOT_EQUAL, Literal: "!=", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 2}}},
		{"==", token.Token{TokenType: token.EQUAL, Literal: "==", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 2}}},
		{"", token.Token{TokenType: token.EOF, Literal: "EOF", Start: token.Pos{Line: 1, Col: 0}, End: token.Pos{Line: 1, Col: 0}}},
		{"89", token.Token{TokenType: token.NUMBER, Literal: "89", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 2}}},
		{"hello", token.Token{TokenType: token.IDENT, Literal: "hello", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 5}}},
		{"89.2", token.Token{TokenType: token.FLOAT, Literal: "89.2", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 4}}},
		{"var", token.Token{TokenType: token.VAR, Literal: "var", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 3}}},
		{"function", token.Token{TokenType: token.FUNCTION, Literal: "function", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 8}}},
		{"if", token.Token{TokenType: token.IF, Literal: "if", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 2}}},
		{"else", token.Token{TokenType: token.ELSE, Literal: "else", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 4}}},
		{"elseif", token.Token{TokenType: token.ELSEIF, Literal: "elseif", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 6}}},
		{"return", token.Token{TokenType: token.RETURN, Literal: "return", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 6}}},
		{"false", token.Token{TokenType: token.FALSE, Literal: "false", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 5}}},
		{"true", token.Token{TokenType: token.TRUE, Literal: "true", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 4}}},
		{`"hello"`, token.Token{TokenType: token.STRING, Literal: "hello", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 7}}},
	}

	for _, test := range tests {
		l := New([]byte(test.input))
		tok, err := l.Lex()
		if err != nil {
			t.Fatalf("Lexer.Lex Lex error. %v", err)
		}
		if tok != test.expectedToken {
			t.Errorf("Lexer.Lex wrong token. got=%q, expected=%+v", tok, test.expectedToken)
		}
	}
}

func TestStringUnicode(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			input:    `"Hello, World!";`,
			expected: "Hello, World!",
		},
		{
			input:    `"Line1\nLine2\tTabbed";`,
			expected: "Line1\nLine2\tTabbed",
		},
		{
			input:    `"She said, \"Hello, World!\"";`,
			expected: `She said, "Hello, World!"`,
		},
		{
			input:    `"Path\\to\\file";`,
			expected: "Path\\to\\file",
		},
		{
			input:    `"This is a string\nthat spans multiple lines."`,
			expected: "This is a string\nthat spans multiple lines.",
		},
		{
			input:    `"Here is a Unicode character: \u2603 (Snowman) \U0001F600";`,
			expected: "Here is a Unicode character: ☃ (Snowman) 😀",
		},
		{
			input:    `"Mix of text, numbers (12345), and special chars: !@#$%^&*()";`,
			expected: "Mix of text, numbers (12345), and special chars: !@#$%^&*()",
		},
		{
			input:    `"";`,
			expected: "",
		},
		{
			input:    `"   This string has spaces at both ends.   ";`,
			expected: "   This string has spaces at both ends.   ",
		},
	}

	for _, test := range tests {
		l := New([]byte(test.input))
		tok, err := l.Lex()
		if err != nil {
			t.Fatalf("Lexer.Lex Lex error. %v", err)
		}
		if tok.Literal != test.expected {
			t.Errorf("Lexer.Lex wrong token literal. got=%s, expected=%+v", tok.Literal, test.expected)
		}
	}
}

func TestLexSourceFile(t *testing.T) {
	byt, err := os.ReadFile("./lexer_test_file.js")
	if err != nil {
		t.Fatal("failed to read file", err)
	}

	tests := []token.Token{
		// Line 1
		{TokenType: token.VAR, Literal: "var", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 3}},
		{TokenType: token.IDENT, Literal: "num", Start: token.Pos{Line: 1, Col: 5}, End: token.Pos{Line: 1, Col: 7}},
		{TokenType: token.ASSIGN, Literal: "=", Start: token.Pos{Line: 1, Col: 9}, End: token.Pos{Line: 1, Col: 9}},
		{TokenType: token.NUMBER, Literal: "89", Start: token.Pos{Line: 1, Col: 11}, End: token.Pos{Line: 1, Col: 12}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 1, Col: 13}, End: token.Pos{Line: 1, Col: 13}},

		// Line 2
		{TokenType: token.VAR, Literal: "var", Start: token.Pos{Line: 2, Col: 1}, End: token.Pos{Line: 2, Col: 3}},
		{TokenType: token.IDENT, Literal: "add", Start: token.Pos{Line: 2, Col: 5}, End: token.Pos{Line: 2, Col: 7}},
		{TokenType: token.ASSIGN, Literal: "=", Start: token.Pos{Line: 2, Col: 9}, End: token.Pos{Line: 2, Col: 9}},
		{TokenType: token.FUNCTION, Literal: "function", Start: token.Pos{Line: 2, Col: 11}, End: token.Pos{Line: 2, Col: 18}},
		{TokenType: token.LPAREN, Literal: "(", Start: token.Pos{Line: 2, Col: 19}, End: token.Pos{Line: 2, Col: 19}},
		{TokenType: token.IDENT, Literal: "a", Start: token.Pos{Line: 2, Col: 20}, End: token.Pos{Line: 2, Col: 20}},
		{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 2, Col: 21}, End: token.Pos{Line: 2, Col: 21}},
		{TokenType: token.IDENT, Literal: "b", Start: token.Pos{Line: 2, Col: 23}, End: token.Pos{Line: 2, Col: 23}},
		{TokenType: token.RPAREN, Literal: ")", Start: token.Pos{Line: 2, Col: 24}, End: token.Pos{Line: 2, Col: 24}},
		{TokenType: token.LBRACE, Literal: "{", Start: token.Pos{Line: 2, Col: 26}, End: token.Pos{Line: 2, Col: 26}},

		// Line 3
		{TokenType: token.RETURN, Literal: "return", Start: token.Pos{Line: 3, Col: 3}, End: token.Pos{Line: 3, Col: 8}},
		{TokenType: token.IDENT, Literal: "a", Start: token.Pos{Line: 3, Col: 10}, End: token.Pos{Line: 3, Col: 10}},
		{TokenType: token.ADD, Literal: "+", Start: token.Pos{Line: 3, Col: 12}, End: token.Pos{Line: 3, Col: 12}},
		{TokenType: token.IDENT, Literal: "b", Start: token.Pos{Line: 3, Col: 14}, End: token.Pos{Line: 3, Col: 14}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 3, Col: 15}, End: token.Pos{Line: 3, Col: 15}},

		// Line 4
		{TokenType: token.RBRACE, Literal: "}", Start: token.Pos{Line: 4, Col: 1}, End: token.Pos{Line: 4, Col: 1}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 4, Col: 2}, End: token.Pos{Line: 4, Col: 2}},

		// Line 6
		{TokenType: token.VAR, Literal: "var", Start: token.Pos{Line: 6, Col: 1}, End: token.Pos{Line: 6, Col: 3}},
		{TokenType: token.IDENT, Literal: "foo", Start: token.Pos{Line: 6, Col: 5}, End: token.Pos{Line: 6, Col: 7}},
		{TokenType: token.ASSIGN, Literal: "=", Start: token.Pos{Line: 6, Col: 9}, End: token.Pos{Line: 6, Col: 9}},
		{TokenType: token.FUNCTION, Literal: "function", Start: token.Pos{Line: 6, Col: 11}, End: token.Pos{Line: 6, Col: 18}},
		{TokenType: token.LPAREN, Literal: "(", Start: token.Pos{Line: 6, Col: 19}, End: token.Pos{Line: 6, Col: 19}},
		{TokenType: token.IDENT, Literal: "a", Start: token.Pos{Line: 6, Col: 20}, End: token.Pos{Line: 6, Col: 20}},
		{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 6, Col: 21}, End: token.Pos{Line: 6, Col: 21}},
		{TokenType: token.IDENT, Literal: "func", Start: token.Pos{Line: 6, Col: 23}, End: token.Pos{Line: 6, Col: 26}},
		{TokenType: token.RPAREN, Literal: ")", Start: token.Pos{Line: 6, Col: 27}, End: token.Pos{Line: 6, Col: 27}},
		{TokenType: token.LBRACE, Literal: "{", Start: token.Pos{Line: 6, Col: 29}, End: token.Pos{Line: 6, Col: 29}},

		// Line 7
		{TokenType: token.RETURN, Literal: "return", Start: token.Pos{Line: 7, Col: 3}, End: token.Pos{Line: 7, Col: 8}},
		{TokenType: token.IDENT, Literal: "func", Start: token.Pos{Line: 7, Col: 10}, End: token.Pos{Line: 7, Col: 13}},
		{TokenType: token.LPAREN, Literal: "(", Start: token.Pos{Line: 7, Col: 14}, End: token.Pos{Line: 7, Col: 14}},
		{TokenType: token.IDENT, Literal: "a", Start: token.Pos{Line: 7, Col: 15}, End: token.Pos{Line: 7, Col: 15}},
		{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 7, Col: 16}, End: token.Pos{Line: 7, Col: 16}},
		{TokenType: token.IDENT, Literal: "a", Start: token.Pos{Line: 7, Col: 18}, End: token.Pos{Line: 7, Col: 18}},
		{TokenType: token.RPAREN, Literal: ")", Start: token.Pos{Line: 7, Col: 19}, End: token.Pos{Line: 7, Col: 19}},
		{TokenType: token.MINUS, Literal: "-", Start: token.Pos{Line: 7, Col: 21}, End: token.Pos{Line: 7, Col: 21}},
		{TokenType: token.IDENT, Literal: "a", Start: token.Pos{Line: 7, Col: 23}, End: token.Pos{Line: 7, Col: 23}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 7, Col: 24}, End: token.Pos{Line: 7, Col: 24}},

		// Line 8
		{TokenType: token.RBRACE, Literal: "}", Start: token.Pos{Line: 8, Col: 1}, End: token.Pos{Line: 8, Col: 1}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 8, Col: 2}, End: token.Pos{Line: 8, Col: 2}},

		// Line 10
		{TokenType: token.IDENT, Literal: "foo", Start: token.Pos{Line: 10, Col: 1}, End: token.Pos{Line: 10, Col: 3}},
		{TokenType: token.LPAREN, Literal: "(", Start: token.Pos{Line: 10, Col: 4}, End: token.Pos{Line: 10, Col: 4}},
		{TokenType: token.NUMBER, Literal: "4", Start: token.Pos{Line: 10, Col: 5}, End: token.Pos{Line: 10, Col: 5}},
		{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 10, Col: 6}, End: token.Pos{Line: 10, Col: 6}},
		{TokenType: token.IDENT, Literal: "add", Start: token.Pos{Line: 10, Col: 8}, End: token.Pos{Line: 10, Col: 10}},
		{TokenType: token.RPAREN, Literal: ")", Start: token.Pos{Line: 10, Col: 11}, End: token.Pos{Line: 10, Col: 11}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 10, Col: 12}, End: token.Pos{Line: 10, Col: 12}},

		// Line 12
		{TokenType: token.VAR, Literal: "var", Start: token.Pos{Line: 12, Col: 1}, End: token.Pos{Line: 12, Col: 3}},
		{TokenType: token.IDENT, Literal: "total", Start: token.Pos{Line: 12, Col: 5}, End: token.Pos{Line: 12, Col: 9}},
		{TokenType: token.ASSIGN, Literal: "=", Start: token.Pos{Line: 12, Col: 11}, End: token.Pos{Line: 12, Col: 11}},
		{TokenType: token.IDENT, Literal: "num", Start: token.Pos{Line: 12, Col: 13}, End: token.Pos{Line: 12, Col: 15}},
		{TokenType: token.MUL, Literal: "*", Start: token.Pos{Line: 12, Col: 17}, End: token.Pos{Line: 12, Col: 17}},
		{TokenType: token.NUMBER, Literal: "90", Start: token.Pos{Line: 12, Col: 19}, End: token.Pos{Line: 12, Col: 20}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 12, Col: 21}, End: token.Pos{Line: 12, Col: 21}},

		// Line 14
		{TokenType: token.STRING, Literal: "hello", Start: token.Pos{Line: 14, Col: 1}, End: token.Pos{Line: 14, Col: 7}},

		// Line 16
		{TokenType: token.LBRACKET, Literal: "[", Start: token.Pos{Line: 16, Col: 1}, End: token.Pos{Line: 16, Col: 1}},
		{TokenType: token.NUMBER, Literal: "1", Start: token.Pos{Line: 16, Col: 2}, End: token.Pos{Line: 16, Col: 2}},
		{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 16, Col: 3}, End: token.Pos{Line: 16, Col: 3}},
		{TokenType: token.NUMBER, Literal: "3", Start: token.Pos{Line: 16, Col: 5}, End: token.Pos{Line: 16, Col: 5}},
		{TokenType: token.COMMA, Literal: ",", Start: token.Pos{Line: 16, Col: 6}, End: token.Pos{Line: 16, Col: 6}},
		{TokenType: token.NUMBER, Literal: "5", Start: token.Pos{Line: 16, Col: 8}, End: token.Pos{Line: 16, Col: 8}},
		{TokenType: token.RBRACKET, Literal: "]", Start: token.Pos{Line: 16, Col: 9}, End: token.Pos{Line: 16, Col: 9}},
		{TokenType: token.SEMICOLON, Literal: ";", Start: token.Pos{Line: 16, Col: 10}, End: token.Pos{Line: 16, Col: 10}},
	}

	l := New(byt)

	for _, test := range tests {
		tok, err := l.Lex()
		if err != nil {
			t.Fatal("Lexer.Lex: error in Lex", err)
		}
		if tok != test {
			t.Errorf("Lexer.Lex wrong token, got=%+v, expected=%+v", tok, test)
		}

		if tok.Start != test.Start {
			t.Errorf("wrong start, got=%+v, expected=%+v", tok.Start, test.Start)
		}
		if tok.End != test.End {
			t.Errorf("wrong end, got=%+v, expected=%+v", tok.End, test.End)
		}
	}
}
