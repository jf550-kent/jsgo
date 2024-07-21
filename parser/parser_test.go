package parser

import (
	"testing"
)

func TestParserError(t *testing.T) {
	tests := []struct{
		filename string
		src []byte
	}{
		{"@", []byte("jdks@")},
		{"var statement", []byte("var 8988")},
		{"hel", []byte("hel")},
	}

	for _, tt := range tests {
		shouldPanic(t, tt.filename, tt.src)
	}
}

func shouldPanic(t *testing.T, filename string, src []byte) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Parse: should panic with filename: %s", filename)
		}
	}()
	Parse(filename, src)
}

// func TestVar(t *testing.T) {
// 	input := []byte("var apple = 89;")
// 	tests := []struct{
// 		input string
// 		expectedVariable string
// 		expectedExpression any
// 	}{
// 		{"var apple = 10;", "apple", }
// 	}

// 	main, err := Parse("", input)
// 	if err != nil {
// 		t.Fatal("Parse error:", err)
// 	}
// 	if len(main.Statements) != 1 {
// 		t.Errorf("main should have 1 statement. got=%d", len(main.Statements))
// 	}
// 	varStmt, ok := main.Statements[0].(*ast.VarStatement)
// 	if !ok {
// 		t.Errorf("wrong statement. expected=%T got=%T", &ast.VarStatement{}, main.Statements[0])
// 	}
// }
