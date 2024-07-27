package parser

import (
	"fmt"
	"testing"

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
		{"var numApp = 10.1;", "numApp", 10.1},
		{"var numApp = apple;", "numApp", 10.1},
	}

	for _, tt := range tests {
		main := Parse("", []byte(tt.input))
		if len(main.Statements) != 1 {
			t.Errorf("main should have 1 statement. got=%d", len(main.Statements))
		}
		varStmt, ok := main.Statements[0].(*ast.VarStatement)
		if !ok {
			t.Errorf("wrong statement. expected=%T got=%T", &ast.VarStatement{}, main.Statements[0])
		}
	
		if varStmt.Variable.Literal != tt.expectedVariable {
			t.Errorf("wrong variable. expected=%s, got=%s", tt.expectedVariable, varStmt.Variable.Literal)
		}

		if !testValueExpression(t, varStmt.Expression, tt.expectedExpression) { return }
	}
}

func testValueExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testNumberValue(t, exp, int64(v))
	case float64:
		return testFloatValue(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}
	return false
}

func testNumberValue(t *testing.T, exp ast.Expression, num int64) bool {
	val, ok := exp.(*ast.Number)
	if !ok {
		t.Errorf("node not *ast.Number, got=%TT", exp)
		return false
	}

	if val.Value != num {
		t.Errorf("wrong number value got=%v expected=%v", val.Value, num)
		return false
	}

	if val.Token.Literal != fmt.Sprintf("%d", num) {
		t.Errorf("wrong number token literal got=%s, expected=%d", val.Token.Literal, num)
		return false
	}

	return true
}

func testFloatValue(t *testing.T, exp ast.Expression, f float64) bool {
	val, ok := exp.(*ast.Float)
	if !ok {
		t.Errorf("node not *ast.Float, got=%T", exp)
	}

	if val.Value != f {
		t.Errorf("wrong float value got=%T expected=%T", val.Value, f)
	}

	if val.Token.Literal != fmt.Sprintf("%v", f) {
		t.Errorf("wrong number token literal got=%s, expected=%v", val.Token.Literal, f)
		return false
	}
	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, i string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("node not *ast.Identifier: got=%T", exp)
		return false
	}

	if ident.String() != i {
		t.Errorf("*ast.Identifier wrong Literal: got=%s, expected=%s", ident.String(), i)
		return false
	}
	return true
}

func testBoolean(t *testing.T, exp ast.Expression, b bool) bool {
	boo, ok := exp.(*ast.Boolean) 
	if !ok {
		t.Errorf("node not *ast.Boolean: got=%T", exp)
		return false
	}

	if boo.Value != b {
		t.Errorf("*ast.Boolean wrong Value: got=%t, expected=%t", boo.Value, b)
		return false
	}
	return true
}