// compiler/compiler_test.go

package compiler

import (
	"testing"

	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
)

type compilerTestCase struct {
	input                string
	expectedConstants    []any
	expectedInstructions []bytecode.Instructions
}

// [CONSTANT 0 CONSTANT 1 ADD POP]
func TestNumberAddition(t *testing.T) {
	tests := []compilerTestCase{
		{input: "1 + 2", expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "90; 45;", expectedConstants: []any{90, 45},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func testCompilerTests(t *testing.T, tests []compilerTestCase) {
	t.Helper()

	for _, tt := range tests {
		main := parser.Parse("", []byte(tt.input))

		compiler := New()
		err := compiler.Compile(main)
		if err != nil {
			t.Fatalf("compiler error: %v", err)
		}

		bytecode := compiler.ByteCode()

		testInstructions(t, tt.expectedInstructions, bytecode.Instructions)
		testConstants(t, tt.expectedConstants, bytecode.Constants)
	}
}

func testInstructions(t *testing.T, expected []bytecode.Instructions, actual bytecode.Instructions) {
	expectedInstructions := mergeInstructions(expected)

	if len(expectedInstructions) != len(actual) {
		t.Errorf("wrong instructions length: got=%d expected=%d", len(actual), len(expectedInstructions))
	}

	for i, ins := range expectedInstructions {
		if actual[i] != ins {
			t.Errorf("wrong instruction at %d got=%q expected=%q", i, actual[i], ins)
		}
	}
}

func testConstants(t *testing.T, expected []any, actual []object.Object) {
	if len(expected) != len(actual) {
		t.Errorf("wrong constant length: got=%d expected=%d", len(actual), len(expected))
	}

	for i, constant := range expected {
		switch constant := constant.(type) {
		case int:
			testNumberObject(t, int64(constant), actual[i])
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

func mergeInstructions(s []bytecode.Instructions) bytecode.Instructions {
	result := bytecode.Instructions{}

	for _, i := range s {
		result = append(result, i...)
	}

	return result
}
