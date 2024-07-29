package evaluator

import (
	"os"
	"testing"

	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
)

func TestEval(t *testing.T) {
	b, _ := os.ReadFile("./example.js")
	main := parser.Parse("", b)
	a := Eval(main)
	print(a)
}
func TestVarStatement(t *testing.T) {
	input := "var a = 89;"
	main := parser.Parse("", []byte(input))
	result := eval(main, object.NewEnvironment())
	num := checkObject[*object.Number](t, result)
	if num.Value != 89 {
		t.Fail()
	}
}

func TestUnaryOperation(t *testing.T) {
	tests := []struct{
		input string
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