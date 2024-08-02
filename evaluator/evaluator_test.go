package evaluator

import (
	"os"
	"testing"

	"github.com/jf550-kent/jsgo/benchmark"
	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
)

func BenchmarkExample(b *testing.B) {
	byt, err := os.ReadFile(benchmark.EXAMPLE_FILE_PATH)
	if err != nil {
		b.Error(err)
	}
	main := parser.Parse("", byt)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		result := Eval(main)
		if result.Type() == object.ERROR_OBJECT {
			b.Error("failed to evalulate main")
		}
	}
}

func TestEval(t *testing.T) {
	b, _ := os.ReadFile("./example.js")
	main := parser.Parse("", b)
	a := Eval(main)
	print(a)
}

func TestVarStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 7;", 35},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		result := evalSetup(tt.input)
		testValue(t, result, tt.expected)
	}
}

func TestUnaryOperation(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"-5;", -5},
		{"--5;", 5},
		{"--5.1;", 5.1},
		{"!5;", false},
		{"!true;", false},
		{"!false;", true},
	}

	for _, tt := range tests {
		main := parser.Parse("", []byte(tt.input))
		result := Eval(main)
		testValue(t, result, tt.expected)
	}
}

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"5;", 5},
		{"10;", 10},
		{"-5;", -5},
		{"-10;", -10},
		{"5 + 5 + 5 + 5 - 10;", 10},
		{"2 * 2 * 2 * 2 * 2;", 32},
		{"-50 + 100 + -50;", 0},
		{"5 * 2 + 10;", 20},
		{"5 + 2 * 10;", 25},
		{"20 + 2 * -10;", 0},
		{"50 / 2 * 2 + 10;", 60},
		{"2 * (5 + 10);", 30},
		{"3 * 3 * 3 + 10;", 37},
		{"3 * (3 * 3) + 10;", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10;", 50},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		testValue(t, evaluated, tt.expected)
	}
}

func evalSetup(src string) object.Object {
	main := parser.Parse("", []byte(src))
	return eval(main, object.NewEnvironment())
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
		{"1 < 2;", true},
		{"1 > 2;", false},
		{"1 < 1;", false},
		{"1 > 1;", false},
		{"1 == 1;", true},
		{"1 != 1;", false},
		{"1 == 2;", false},
		{"1 != 2;", true},
		{"true == true;", true},
		{"false == false;", true},
		{"true == false;", false},
		{"true != false;", true},
		{"false != true;", true},
		{"(1 < 2) == true;", true},
		{"(1 < 2) == false;", false},
		{"(1 > 2) == true;", false},
		{"(1 > 2) == false;", true},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		testValue(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true;", false},
		{"true;", true},
		{"!5;", false},
		{"!!true;", true},
		{"!!false;", false},
		{"!!5;", true},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		testBool(t, evaluated, tt.expected)
	}
}

func TestIfElseCondition(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"if (true) { 40; };", 40},
		{"if (false) { 10; };", nil},
		{"if (1) { 10; };", 10},
		{"if (1 < 2) { 20; };", 20},
		{"if (1 > 2) { 10; };", nil},
		{"if (1 > 2) { 10; } else { 20; };", 20},
		{"if (1 < 2) { 10; } else { 20; };", 10},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		result, ok := tt.expected.(int)
		if ok {
			testNumber(t, evaluated, int64(result))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return 10;
				};
				return 1;
			};`, 10,
		},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		testValue(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: NUMBER + BOOLEAN"},
		{"5 + true; 5;", "type mismatch: NUMBER + BOOLEAN"},
		{"-true;", "unable to minus value: -true for type: BOOLEAN"},
		{"true + false;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5;", "unknown operator: BOOLEAN + BOOLEAN"},
		{"if (10 > 1) { true + false; };", "unknown operator: BOOLEAN + BOOLEAN"},
		{`if (10 > 1) {
			if (10 > 1) {
				return true - false;
			};
			return 1;
		};`, "unknown operator: BOOLEAN - BOOLEAN"},
		{"foobar;", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message: expected=%q, got=%q", tt.expectedMessage, errObj.Message)
			continue
		}
	}
}

func TestFunctionObject(t *testing.T) {
	input := "function(x) { x + 3; };"
	evaluated := evalSetup(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
var first = 10;
var second = 10;
var third = 10;

var ourFunction = function(first) {
  var second = 20;

  return first + second + third;
};

ourFunction(20) + first + second;`

	main := evalSetup(input)
	testNumber(t, main, int64(70))
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"var a = function(x) { x; }; a(5);", 5},
		{"var identity = function(x) { return x; }; identity(5);", 5},
		{"var double = function(x) { x * 2; }; double(5);", 10},
		{"var add = function(x, y) { x + y; }; add(5, 5);", 10},
		{"var add = function(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
	}

	for _, tt := range tests {
		testValue(t, evalSetup(tt.input), tt.expected)
	}
}

func TestAssignment()

func checkObject[expected any](t *testing.T, obj object.Object) expected {
	if obj == nil {
		t.Fatal("object is nil")
	}

	v, ok := obj.(expected)
	if !ok {
		t.Fatalf("object wrong type: got=%T expected=%T", obj, v)
	}
	return v
}

func testValue(t *testing.T, obj object.Object, expectedValue any) {
	switch v := expectedValue.(type) {
	case int:
		testNumber(t, obj, int64(v))
		return
	case float64:
		testFloat(t, obj, v)
		return
	case bool:
		testBool(t, obj, v)
		return
	}
	t.Errorf("invalid type for expected value")
}

func testBool(t *testing.T, obj object.Object, v bool) {
	b := checkObject[*object.Boolean](t, obj)

	if b.Value != v {
		t.Errorf("wrong boolean value: got=%v expected=%v", b.Value, v)
	}
}

func testNumber(t *testing.T, obj object.Object, v int64) {
	num := checkObject[*object.Number](t, obj)

	if num.Value != v {
		t.Errorf("wrong number value: got=%v expected=%v", num.Value, v)
	}
}

func testFloat(t *testing.T, obj object.Object, v float64) {
	fl := checkObject[*object.Float](t, obj)

	if fl.Value != v {
		t.Errorf("wrong number value: got=%v expected=%v", fl.Value, v)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
