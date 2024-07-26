package ast

import (
	"testing"

	"github.com/jf550-kent/jsgo/token"
)

func TestIFExpression(t *testing.T) {
	ifToken := token.Token{TokenType: token.IF, Literal: "if", Start: token.Pos{Line: 1, Col: 1}, End: token.Pos{Line: 1, Col: 2}}
	number := token.Token{TokenType: token.NUMBER, Literal: "4", Start: token.Pos{Line: 1, Col: 4}, End: token.Pos{Line: 1, Col: 5}}
	n := &Number{Token: number, Value: 4}
	varStatements := []Statement{
		&VarStatement{Token: token.Token{Literal: "var", Start: token.Pos{Line: 2, Col: 1}, End: token.Pos{Line: 2, Col: 4}},
			Variable:   &Identifier{Token: token.Token{Literal: "x", Start: token.Pos{Line: 2, Col: 5}, End: token.Pos{Line: 2, Col: 6}}, Literal: "x"},
			Expression: &Number{Token: token.Token{Literal: "123", Start: token.Pos{Line: 2, Col: 9}, End: token.Pos{Line: 2, Col: 11}}, Value: 123},
		},
	}
	body := &BlockStatement{Token: token.Token{Literal: "{", Start: token.Pos{Line: 2, Col: 7}, End: token.Pos{Line: 2, Col: 8}}, Statements: varStatements}
	elseBlock := &BlockStatement{Token: token.Token{Literal: "{", Start: token.Pos{Line: 3, Col: 7}, End: token.Pos{Line: 3, Col: 8}}, Statements: varStatements}
	ifExpr := IFExpression{Token: ifToken, Condition: n, Body: body, Else: elseBlock}

	if ifExpr.String() != "if (4) {var x = 123; } else {var x = 123; };" {
		t.Errorf("error IFExpression.String got=%s, want=%s", ifExpr.String(), "if (4) {var x = 123; } else {var x = 123; };")
	}
}
