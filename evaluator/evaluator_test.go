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
		expected any
	}{
		{"var a = 5; a;", 5},
		{"var a = 5 * 7;", 35},
		{"var a = 5; var b = a; b;", 5},
		{"var a = 5; var b = a; var c = a + b + 5; c;", 15},
		{`var hello = "Hello world"; hello;`, "Hello world"},
		{`var hello = null; hello;`, nil},
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

// add for loop checking
func TestAssignment(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"var a = 90; if (a) { a = 199 }; a;", 199},
		{"var a = 90; if (a) { var a = 199; }; a; ", 199},
		{"var a = 90; var add = function() { var a = true; a = 100; }; add(); a;", 90},
		{"var a = 90; var add = function() { var a = true; return a; }; add();", true},
	}

	for _, tt := range tests {
		testValue(t, evalSetup(tt.input), tt.expected)
	}
}

func TestFor(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"for (var i = 0; i < 5; i = i + 1) {}; i;", 5},
		{"var i = 0; for (; i < 5; i = i + 1) {}; i;", 5},
		{"for (var i = 5; i < 0; i = i + 1) {}; i;", 5},
		{"var sum = 0; for (var i = 0; i < 3; i = i + 1) { for (var j = 0; j < 2; j = j + 1) { sum = sum + 1; } }; sum;", 6},
		{"var sum = 0; for (var i = 0; i < 3; i = i + 1) { for (var j = 0; j < 2; j = j + 1) { for (var k = 0; k < 2; k = k + 1) { sum = sum + 1; } } }; sum;", 12},
		{"var three = function() { for (var i = 0; i < 5; i = i + 1) { if (i == 3) { return i; } } }; three();", 3},
		{"for (;false;) {};", NULL},
		{"var flag = true; for (var i = 0; flag; i = i + 1) { if (i == 200) { flag = false; }}; i;", 201},
		{"for (var i = 0; i != 10; i = i + 2) {}; i;", 10},
		{"for (var i = 0; i < 0; i = i + 1) {}; i;", 0},
		{"var sum = 0; for (var i = 0; i < 10000; i = i + 1) { for (var j = 10; i < 10; j = 20) { i = 20 + j; } if (i > 200) { sum = i; }; }; sum;", 9999},
	}

	for _, tt := range tests {
		eval := evalSetup(tt.input)
		testValue(t, eval, tt.expected)
	}
}

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

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := evalSetup(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Body) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Body))
	}

	testValue(t, result.Body[0], 1)
	testValue(t, result.Body[1], 4)
	testValue(t, result.Body[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{"[1, 2, 3][0]", 1},
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][2]", 3},
		{"var i = 0; [1][i];", 1},
		{"[1, 2, 3][1 + 1];", 3},
		{"var myArray = [1, 2, 3]; myArray[2];", 3},
		{"var myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];", 6},
		{"var myArray = [1, 2, 3]; var i = myArray[0]; myArray[i]", 2},
		{"[1, 2, 3][3]", nil},
		{"[1, 2, 3][-1]", nil},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testValue(t, evaluated, integer)
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestArrayLength(t *testing.T) {
	input := `var arr = [1, 3, 4]; arr["length"];`
	evaluated := evalSetup(input)
	testValue(t, evaluated, 3)
}

func TestArrayPush(t *testing.T) {
	input := `var arr = [1, 3, 4]; arr["push"](9); arr["length"];`
	evaluated := evalSetup(input)
	testValue(t, evaluated, 4)
}

func TestArrayFunctionCall(t *testing.T) {
	input := `var add = function (a) { return a + a; }; var arr = [1, 3, 4, add]; arr[3](9);`
	evaluated := evalSetup(input)
	testValue(t, evaluated, 18)
}

func TestDictionaryExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`{"foo": 5}["foo"]`, 5},
		{`{"foo": 5}["bar"]`, nil},
		{`var key = "foo"; {"foo": 5}[key]`, 5},
		{`{}["foo"]`, nil},
		{`{5: 5}[5]`, 5},
		{`{true: 5}[true]`, 5},
		{`{false: 5}[false]`, 5},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		testValue(t, evaluated, tt.expected)
	}
}

func TestDictionaryDeclaration(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{`var apple = {"color": "red"}; apple["taste"] = "red"; apple["taste"];`, "red"},
	}

	for _, tt := range tests {
		evaluated := evalSetup(tt.input)
		testValue(t, evaluated, tt.expected)
	}
}

func testValue(t *testing.T, obj object.Object, expectedValue any) {
	switch v := expectedValue.(type) {
	case int:
		testNumber(t, obj, int64(v))
	case float64:
		testFloat(t, obj, v)
	case bool:
		testBool(t, obj, v)
	case *object.Null:
		if obj != v {
			t.Errorf("expecting null got=%v", v)
		}
	case string:
		testString(t, obj, v)
	case nil:
		testNullObject(t, obj)
	default:
		t.Errorf("invalid type for expected value")
	}
}

func testBool(t *testing.T, obj object.Object, v bool) {
	b := checkObject[*object.Boolean](t, obj)

	if b.Value != v {
		t.Errorf("wrong boolean value: got=%v expected=%v", b.Value, v)
	}
}

func testString(t *testing.T, obj object.Object, v string) {
	s := checkObject[*object.String](t, obj)

	if s.Value != v {
		t.Errorf("wrong string value: got=%v expected=%v", s.Value, v)
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
