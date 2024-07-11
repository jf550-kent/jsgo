package parser

import (
	"testing"

	"github.com/jun-hf/jsgo/ast"
)

func TestVar(t *testing.T) {
	input := []byte("var apple = 89;")
	tests := []struct{
		input string
		expectedVariable string
		expectedExpression any
	}{
		{"var apple = 10;", "apple", }
	}

	main, err := Parse("", input)
	if err != nil {
		t.Fatal("Parse error:", err)
	}
	if len(main.Statements) != 1 {
		t.Errorf("main should have 1 statement. got=%d", len(main.Statements))
	}
	varStmt, ok := main.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Errorf("wrong statement. expected=%T got=%T", &ast.VarStatement{}, main.Statements[0])
	}
}

func testVarStatement(t *testing.T, s ast.Statement,)