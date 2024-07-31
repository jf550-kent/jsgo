package parser

import (
	"fmt"
	"os"
	"testing"

	"github.com/jf550-kent/jsgo/ast"
)

func BenchmarkExample(b *testing.B) {
	byt, err := os.ReadFile("./../benchmark/example.js")
	if err != nil {
		b.Fatal("failed to read file", err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		main := Parse("", byt)
		if len(main.Statements) != 5 {
			b.Fatal("parser failed")
		}
	}
}
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
		{"var numApp = apple;", "numApp", "apple"},
	}

	for _, tt := range tests {
		main := Parse("", []byte(tt.input))
		if len(main.Statements) != 1 {
			t.Errorf("main should have 1 statement. got=%d", len(main.Statements))
		}

		varStmt := checkStatement[*ast.VarStatement](t, main.Statements[0])
		if varStmt.Variable.Literal != tt.expectedVariable {
			t.Errorf("wrong variable. expected=%s, got=%s", tt.expectedVariable, varStmt.Variable.Literal)
		}

		testValueExpression(t, varStmt.Expression, tt.expectedExpression)
	}
}

func TestReturn(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 10;", 10},
		{"return true;", true},
		{"return apple;", "apple"},
	}

	for _, tt := range tests {
		main := Parse("", []byte(tt.input))

		if len(main.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(main.Statements))
		}

		returnStmt := checkStatement[*ast.ReturnStatement](t, main.Statements[0])
		if returnStmt.Token.String() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q", returnStmt.Token.String())
		}

		testValueExpression(t, returnStmt.ReturnExpression, tt.expectedValue)
	}
}

func TestBinaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		left     any
		operator string
		right    any
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"foobar + barfoo;", "foobar", "+", "barfoo"},
		{"foobar - barfoo;", "foobar", "-", "barfoo"},
		{"foobar * barfoo;", "foobar", "*", "barfoo"},
		{"foobar / barfoo;", "foobar", "/", "barfoo"},
		{"foobar > barfoo;", "foobar", ">", "barfoo"},
		{"foobar < barfoo;", "foobar", "<", "barfoo"},
		{"foobar == barfoo;", "foobar", "==", "barfoo"},
		{"foobar != barfoo;", "foobar", "!=", "barfoo"},
		{"true == true;", true, "==", true},
		{"true != false;", true, "!=", false},
		{"false == false;", false, "==", false},
	}

	for _, tt := range tests {
		main := Parse(tt.input, []byte(tt.input))
		if len(main.Statements) != 1 {
			t.Fatal("number of main Statements is not 1")
		}
		exprStmt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
		testBinaryExpression(t, exprStmt.Expression, tt.left, tt.operator, tt.right)
	}
}

func TestUnaryExpression(t *testing.T) {
	tests := []struct {
		input    string
		operator string
		value    any
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!foobar;", "!", "foobar"},
		{"-foobar;", "-", "foobar"},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range tests {
		main := Parse(tt.input, []byte(tt.input))

		if len(main.Statements) != 1 {
			t.Fatalf("main.Statements does not contain %d statements. got=%d\n", 1, len(main.Statements))
		}
		exprStmt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
		expr := checkExpression[*ast.UnaryExpression](t, exprStmt.Expression)

		if expr.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q, got=%q", tt.operator, expr.Operator)
		}
		testValueExpression(t, expr.Expression, tt.value)
	}
}

func TestFunctionDeclaration(t *testing.T) {
	input := "function (a, b) { x; };"
	expectedParameter := []string{"a", "b"}

	main := Parse("func", []byte(input))
	if len(main.Statements) != 1 {
		t.Fatal("statement is not one")
	}

	exprStmt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
	funcExprs := checkExpression[*ast.FunctionDeclaration](t, exprStmt.Expression)

	if funcExprs.Token.Literal != "function" {
		t.Errorf("wrong function literal, got=%s, expected=function", funcExprs.Token.Literal)
	}

	if len(funcExprs.Parameters) != 2 {
		t.Error("functions paramters is not 2")
	}

	for i, p := range funcExprs.Parameters {
		testValueExpression(t, p, expectedParameter[i])
	}

	block := checkStatement[*ast.BlockStatement](t, funcExprs.Body)
	if len(block.Statements) != 1 {
		t.Error("block statement is not 1")
	}
	blockExprStmt := checkStatement[*ast.ExpressionStatement](t, block.Statements[0])
	blockExpr := checkExpression[*ast.Identifier](t, blockExprStmt.Expression)

	testValueExpression(t, blockExpr, "x")
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b;",
			"((-a) * (b))",
		},
		{
			"!-a;",
			"(!(-(a)))",
		},
		{
			"a + b + c;",
			"((a + b) + (c))",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4 - 5 * 5;",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true;",
			"true",
		},
		{
			"false;",
			"false",
		},
		{
			"3 > 5 == false;",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true;",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4;",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2;",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5);",
			"(2 / (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5);",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5);",
			"(-(5 + 5))",
		},
		{
			"!(true == true);",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		main := Parse("", []byte(tt.input))
		if len(main.Statements) != 1 {
			t.Error("statement should be one")
		}

		checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	main := Parse("", []byte(input))

	if len(main.Statements) != 1 {
		t.Fatalf("main.Statements does not contain %d statements. got=%d\n",
			1, len(main.Statements))
	}

	stmt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
	expr := checkExpression[*ast.CallExpression](t, stmt.Expression)

	testValueExpression(t, expr.Function, "add")
	if len(expr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(expr.Arguments))
	}

	testValueExpression(t, expr.Arguments[0], 1)
	testBinaryExpression(t, expr.Arguments[1], 2, "*", 3)
	testBinaryExpression(t, expr.Arguments[2], 4, "+", 5)
}

func testBinaryExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	binExpr := checkExpression[*ast.BinaryExpression](t, exp)
	testValueExpression(t, binExpr.Left, left)
	if binExpr.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q", operator, binExpr.Operator)
		return false
	}
	testValueExpression(t, binExpr.Right, right)
	return true
}

func TestIfExpression(t *testing.T) {
	input := `if (x) { x; };`
	main := Parse("", []byte(input))

	if len(main.Statements) != 1 {
		t.Fatalf("main.Body does not contain 1 statement. got=%d\n", len(main.Statements))
	}
	exprSt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
	ifExpr := checkExpression[*ast.IFExpression](t, exprSt.Expression)

	if ifExpr.Condition == nil {
		t.Fatalf("*ast.IfExpression condition is nil")
	}
	testValueExpression(t, ifExpr.Condition, "x")

	body := checkStatement[*ast.BlockStatement](t, ifExpr.Body)
	if len(body.Statements) != 1 {
		t.Errorf("body has more than one statement")
	}

	bodyExpr := checkStatement[*ast.ExpressionStatement](t, body.Statements[0])
	testValueExpression(t, bodyExpr.Expression, "x")

	if ifExpr.Else != nil {
		t.Errorf("*ast.IfExpression has an unexpected else clause")
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x) { x; } else { 10; };`
	main := Parse("", []byte(input))

	if len(main.Statements) != 1 {
		t.Fatalf("main.Body does not contain 1 statement. got=%d\n", len(main.Statements))
	}
	exprSt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])
	ifExpr := checkExpression[*ast.IFExpression](t, exprSt.Expression)

	if ifExpr.Condition == nil {
		t.Fatalf("*ast.IfExpression condition is nil")
	}
	testValueExpression(t, ifExpr.Condition, "x")

	body := checkStatement[*ast.BlockStatement](t, ifExpr.Body)
	if len(body.Statements) != 1 {
		t.Errorf("body has more than one statement")
	}

	bodyExpr := checkStatement[*ast.ExpressionStatement](t, body.Statements[0])
	testValueExpression(t, bodyExpr.Expression, "x")

	elseSt := checkStatement[*ast.BlockStatement](t, ifExpr.Else)
	if len(elseSt.Statements) != 1 {
		t.Error("else block has more than one statement")
	}

	elseExpr := checkStatement[*ast.ExpressionStatement](t, elseSt.Statements[0])
	testValueExpression(t, elseExpr.Expression, 10)
}

// add negative test case
func TestExpressionStatement(t *testing.T) {
	tests := []struct {
		input string
		value any
	}{
		{"5;", 5},
		{"10;", 10},
		{"900;", 900},
		{"9223372036854775807;", 9223372036854775807},
		{"5.98;", 5.98},
		{"7.89;", 7.89},
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		main := Parse(tt.input, []byte(tt.input))

		if len(main.Statements) != 1 {
			t.Fatal("number of main Statements is not 1")
		}
		exprStmt := checkStatement[*ast.ExpressionStatement](t, main.Statements[0])

		if exprStmt.Expression == nil {
			t.Fatal("expression is nil")
		}
		testValueExpression(t, exprStmt.Expression, tt.value)
	}
}

func checkStatement[expected any](t *testing.T, stmt ast.Statement) expected {
	if stmt == nil {
		t.Fatal("statement is nil")
	}
	v, ok := stmt.(expected)
	if !ok {
		t.Fatalf("statement wrong type: got=%T expected=%T", stmt, v)
	}
	return v
}

func checkExpression[expected any](t *testing.T, expr ast.Expression) expected {
	if expr == nil {
		t.Fatal("expression is nil")
	}

	v, ok := expr.(expected)
	if !ok {
		t.Fatalf("exppresion wrong type: got=%T expected=%T", expr, v)
	}
	return v
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
