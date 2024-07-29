package evaluator

import (
	"testing"

	"github.com/jf550-kent/jsgo/parser"
	"github.com/jf550-kent/jsgo/object"
)

func TestVarStatement(t *testing.T) {
	input := "var a = 89;"
	main := parser.Parse("", []byte(input))
	result := Eval(main, object.NewEnvironment())
	num := checkObject[*object.Number](t, result)
	if num.Value != 89 {
		t.Fail()
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