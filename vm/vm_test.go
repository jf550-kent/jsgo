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

func TestNumberAddition(t *testing.T) {
	tests := []vmTestCase{
		{"3", 3},
		{"9", 9},
		{"8 + 1", 9},
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

		result := vm.StackTop()

		testObject(t, tt.expected, result)
	}
}

func testObject(t *testing.T, expected any, actual object.Object) {
	switch expected := expected.(type) {
	case int:
		testNumberObject(t, int64(expected), actual)
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
