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
		i := com.ByteCode().Instructions.String()
		print(i)

		result := vm.lastPopStack()

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
