package parser

import (
	"testing"
	"reflect"

	"github.com/jf550-kent/jsgo/ast"
)

func TestParserError(t *testing.T) {
	tests := []struct {
		filename string
		src      []byte
	}{
		{"@", []byte("jdks@")},
		{"var statement", []byte("var 8988")},
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

func TestVar(t *testing.T) {
	tests := []struct {
		input              string
		expectedVariable   string
		expectedExpression any
	}{
		{"var apple = 10;", "apple", 10},
		{"var yellow = true;", "yellow", true},
		{"var numApp = 10.1", "numApp", 10.1},
	}

	for i, tt := range tests {
		main := Parse("", []byte(tt.input))
		if len(main.Statements) != 1 {
			t.Errorf("main should have 1 statement. got=%d", len(main.Statements))
		}
		varStmt, ok := main.Statements[0].(*ast.VarStatement)
		if !ok {
			t.Errorf("wrong statement. expected=%T got=%T", &ast.VarStatement{}, main.Statements[0])
		}
	
		if varStmt.Variable.Literal != tt.expectedVariable {
			t.Errorf("wrong variable. expected=%s, got=%s", tests[0].expectedVariable, varStmt.Variable.Literal)
		}

		
	}
}
