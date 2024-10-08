package vm

import (
	"testing"

	"github.com/jf550-kent/jsgo/compiler"
	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
)

type vmTestCase struct {
	input    string
	expected any
}

func TestNumberOperation(t *testing.T) {
	tests := []vmTestCase{
		{"3", 3},
		{"9", 9},
		{"8 + 1", 9},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"20 << 10", 20480},
		{"99 ^ 8", 107},
		{"-50", -50},
		{"-90", -90},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	testVmTests(t, tests)
}

func TestFloatOperation(t *testing.T) {
	tests := []vmTestCase{
		{"1.0", 1.0},
		{"3.14", 3.14},
		{"1.5 + 2.5", 4.0},
		{"5.0 - 1.2", 3.8},
		{"2.5 * 4.0", 10.0},
		{"10.0 / 3.0", 3.3333333333333335},
		{"(2.0 + 3.0) * 4.0", 20.0},
		{"5.0 + 2.0 * 3.0", 11.0},
		{"(5.0 + 2.0) / 3.0", 2.3333333333333335},
		{"7.0 / (2.0 + 3.0)", 1.4},
		{"9.0 - (2.0 + 1.0) * 2.0", 3.0},
	}

	testVmTests(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []vmTestCase{
		{"true;", true},
		{"false;", false},
		{"3 < 4", true},
		{"3 > 4", false},
		{"3 < 3", false},
		{"3 > 3", false},
		{"3 == 3", true},
		{"3 != 3", false},
		{"3 == 4", false},
		{"3 != 4", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(3 < 4) == true", true},
		{"(3 < 4) == false", false},
		{"(3 > 4) == true", false},
		{"(3 > 4) == false", true},
		{"(3.10 > 4) == false", true},
		{"(3.10 < 4) == true", true},
		{"(3.10 == 4) == true", false},
		{"(3.10 != 4) == false", false},
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	testVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 100 }", 100},
		{"if (true) { 190 } else { 20 }", 190},
		{"if (false) { 10 } else { 33 } ", 33},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", NULL},
		{"if (false) { 10 }", NULL},
		{"!(if (false) { 5; })", true},
	}

	testVmTests(t, tests)
}

func TestGlobalStatement(t *testing.T) {
	tests := []vmTestCase{
		{"var apple = 99; apple", 99},
		{"var one = 80.9; var two = 2; one + two", 82.9},
		{"var one = 1; var two = one + one; one + two", 3},
		{"var one = 1; one = 29 one;", 29},
	}

	testVmTests(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []vmTestCase{
		{`"Hello world"`, "Hello world"},
	}
	testVmTests(t, tests)
}

func TestArray(t *testing.T) {
	tests := []vmTestCase{
		{"[]", []int{}},
		{"[8, 9, 10]", []int{8, 9, 10}},
		{"[1 + 3, 5 * 6, 9 + 1]", []int{4, 30, 10}},
	}

	testVmTests(t, tests)
}

func TestDictionary(t *testing.T) {
	tests := []vmTestCase{
		{
			"{}", map[object.Hash]int64{},
		},
		{
			"{1: 2, 2: 3}",
			map[object.Hash]int64{
				(&object.Number{Value: 1}).Hash(): 2,
				(&object.Number{Value: 2}).Hash(): 3,
			},
		},
		{
			"{1 + 1: 2 * 2, 3 + 3: 4 * 4}",
			map[object.Hash]int64{
				(&object.Number{Value: 2}).Hash(): 4,
				(&object.Number{Value: 6}).Hash(): 16,
			},
		},
	}

	testVmTests(t, tests)
}

func TestIndexing(t *testing.T) {
	tests := []vmTestCase{
		{"[1, 2, 3][1]", 2},
		{"[1, 2, 3][0 + 2]", 3},
		{"[[1, 1, 1]][0][0]", 1},
		{"[][0]", NULL},
		{"[1, 2, 3][99]", NULL},
		{"[1][-1]", NULL},
		{"{1: 1, 2: 2}[1]", 1},
		{"{1: 1, 2: 2}[2]", 2},
		{"{1: 1}[0]", NULL},
		{"{}[0]", NULL},
	}

	testVmTests(t, tests)
}

func TestCallingFunction(t *testing.T) {
	tests := []vmTestCase{
		{input: "var onePlusOne = function() { 1 + 1}; onePlusOne();", expected: 2},
		{input: `var one = function() { 1; }; var two = function() { 2; }; one() + two();`, expected: 3},
		{input: "var a = function() { 1 }; var b = function() { a() + 1 }; var c = function() { b() + 1 }; c();", expected: 3},
		{input: "var a = function() { return 1; 3; 2; }; a();", expected: 1},
		{input: "var a = function() {  1; return 3; return 2; }; a();", expected: 3},
		{input: "var a = function() { }; a();", expected: NULL},
	}
	testVmTests(t, tests)
}

func TestCallingFunctionsWithLocalScope(t *testing.T) {
	tests := []vmTestCase{
		{input: "var foo = function() { var apple = 98; apple; }; foo();", expected: 98},
		{input: "var sum = function() { var first = 9; var second = 9; return first + second; }; sum()", expected: 18},
		{input: "var sum = function() { var first = 9; var second = 9; return first + second; }; var ten = function() { var five = 5; return five + five;}; ten() + sum()", expected: 28},
		{input: "var globalNumber = 90; var minus = function() { var num = 1 return globalNumber - num; } var add = function() { var num = 8 return globalNumber + num; } minus() + add()", expected: 187},
		{input: "var sum = function(a, b) { return a + b; }; sum(2,2)", expected: 4},
		{input: "var sum = function(a, b) { var c = a + b; return c; }; sum(1, 2);", expected: 3},
		{input: "var sum = function(a, b) { var c = a + b; return c; }; sum(1, 2) + sum(3, 4);", expected: 10},
		{input: "var sum = function(a, b) { var c = a + b; return c; }; var outer = function() { return sum(1, 2) + sum(3, 4); }; outer();", expected: 10},
	}

	testVmTests(t, tests)
}

func TestClosure(t *testing.T) {
	tests := []vmTestCase{
		{input: "var close = function(a) { function() {a;}}; var closure = close(100); closure()", expected: 100},
	}
	testVmTests(t, tests)
}

func TestRecursiveFibonacci(t *testing.T) {
	tests := []vmTestCase{
		{
			input: `
			var fibonacci = function(x) {
					if (x == 0) {
							return 0;
					} else {
							if (x == 1) {
									return 1;
							} else {
									fibonacci(x - 1) + fibonacci(x - 2);
							}
					}
			};
			fibonacci(15);
			`,
			expected: 610,
		},
	}

	testVmTests(t, tests)
}

func TestBracket(t *testing.T) {
	tests := []vmTestCase{
		{input: "var arr = [10]; arr[1] = 90; arr;", expected: []int{10, 90}},
		{
			input: `var dic = { "next": 10}; dic["current"] = 20 dic;`,
			expected: map[object.Hash]int64{
				(&object.String{Value: "next"}).Hash():    10,
				(&object.String{Value: "current"}).Hash(): 20,
			},
		},
	}

	testVmTests(t, tests)
}

func TestForLoop(t *testing.T) {
	tests := []vmTestCase{
		// {input: "for (var i = 0; i < 10; i = i + 1) { 29; }; i;", expected: 10},
		// {input: "var i = 0 for (; i < 0;) { 29; }; i;", expected: 0},
		{input: "var i = 0 for (; i < 0;) { 29; }; i;", expected: 0},
		{input: `
		
		var isShorterThan = function(x, y) {

  for (;y != null;) {
    if (x == null) { return 8888; }
	return 90;
  }
  return 323;
}

isShorterThan(null, null)
		`, expected: 323},
		{`var y = null;

		var x = null
		for (;y != null;) {
    		if (x == null) { return 8888; }
			return 90;
  		}`, nil},
	}

	testVmTests(t, tests)
}

func TestDebug(t *testing.T) {
	input := `var y = null;

		var x = null
		for (;y != null;) {
    		if (x == null) { return 8888; }
			return 90;
  		}`
	main := parser.Parse("", []byte(input))

	com := compiler.New()
	if err := com.Compile(main); err != nil {
		t.Fatalf("compiler error: %s", err)
	}

	vm := New(com.ByteCode())
	if err := vm.Run(); err != nil {
		t.Fatalf("vm error: %s", err)
	}
}

func TestBuiltIn(t *testing.T) {
	tests := []vmTestCase{
		{`var arr = []; arr["push"](10); arr["push"](19); arr;`, []int{10, 19}},
		{`var arr = [1, 2, 3]; arr["size"];`, 3},
	}
	testVmTests(t, tests)
}

func testVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		main := parser.Parse("", []byte(tt.input))

		com := compiler.New()
		if err := com.Compile(main); err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(com.ByteCode())
		if err := vm.Run(); err != nil {
			t.Fatalf("vm error: %s", err)
		}

		result := vm.LastPopStack()

		testObject(t, tt.expected, result)
	}
}

func testObject(t *testing.T, expected any, actual object.Object) {
	switch expected := expected.(type) {
	case int:
		testNumberObject(t, int64(expected), actual)
	case float64:
		testFloat(t, expected, actual)
	case bool:
		testBoolean(t, expected, actual)
	case string:
		testString(t, expected, actual)
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object not Array: %T (%+v)", actual, actual)
		}
		if len(array.Body) != len(expected) {
			t.Errorf("wrong num of elements. want=%d, got=%d", len(expected), len(array.Body))
		}
		for i, expectedElem := range expected {
			testNumberObject(t, int64(expectedElem), array.Body[i])
		}
	case map[object.Hash]int64:
		dic, ok := actual.(*object.Dictionary)
		if !ok {
			t.Errorf("object is not Dictinary. got=%T (%+v)", actual, actual)
		}

		if len(dic.Value) != len(expected) {
			t.Errorf("dictionary has wrong number of key-value pair. want=%d, got=%d", len(expected), len(dic.Value))
		}

		for expectedKey, expectedValue := range expected {
			pair, ok := dic.Value[expectedKey]
			if !ok {
				t.Errorf("no pair for given key in Pairs")
			}
			testNumberObject(t, expectedValue, pair.Value)
		}
	case *object.Null:
		if expected != NULL {
			t.Errorf("expected null got=%v", expected)
		}
	}
}

func testNumberObject(t *testing.T, constant int64, act object.Object) {
	actual, ok := act.(*object.Number)
	if !ok {
		t.Errorf("expecting number got=%T", act)
	}

	if constant != actual.Value {
		t.Errorf("wrong number value: got=%d expected=%d", actual.Value, constant)
	}
}

func testFloat(t *testing.T, constant float64, obj object.Object) {
	fl := checkObject[*object.Float](t, obj)

	if fl.Value != constant {
		t.Errorf("wrong float value: got=%v expected=%v", fl.Value, constant)
	}
}

func testBoolean(t *testing.T, constant bool, obj object.Object) {
	fl := checkObject[*object.Boolean](t, obj)

	if fl.Value != constant {
		t.Errorf("wrong boolean value: got=%v expected=%v", fl.Value, constant)
	}
}

func testString(t *testing.T, v string, obj object.Object) {
	s := checkObject[*object.String](t, obj)

	if s.Value != v {
		t.Errorf("wrong string value: got=%v expected=%v", s.Value, v)
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
