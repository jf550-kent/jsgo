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

func TestCompileOperation(t *testing.T) {
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
		{input: "90 - 45;", expectedConstants: []any{90, 45},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "90 / 45;", expectedConstants: []any{90, 45},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDiv),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "90 * 45;", expectedConstants: []any{90, 45},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "90 << 45;", expectedConstants: []any{90, 45},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSHL),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "90 ^ 45;", expectedConstants: []any{90, 45},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpXOR),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "-1",
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpMinus),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestBooleanExpression(t *testing.T) {
	tests := []compilerTestCase{
		{input: "true;", expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{input: "false;", expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 > 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 < 2",
			expectedConstants: []any{2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpGreaterThan),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 == 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "1 != 2",
			expectedConstants: []any{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "true == false",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "true != false",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpFalse),
				bytecode.Make(bytecode.OpNotEqual),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "!true",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpBang),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	testCompilerTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
				if (true) { 10 } else { 20 }; 3333;
				`,
			expectedConstants: []any{10, 20, 3333},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),
				bytecode.Make(bytecode.OpJumpNotTrue, 10),
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpJump, 13),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},

		{
			input: `
			if (true) { 10 }; 3333;
			`,
			expectedConstants: []any{10, 3333},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpTrue),            // 0
				bytecode.Make(bytecode.OpJumpNotTrue, 10), // 1
				bytecode.Make(bytecode.OpConstant, 0),     // 4
				bytecode.Make(bytecode.OpJump, 11),        // 7
				bytecode.Make(bytecode.OpNull),            // 10
				bytecode.Make(bytecode.OpPop),             // 11
				bytecode.Make(bytecode.OpConstant, 1),     // 12
				bytecode.Make(bytecode.OpPop),             // 15
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestVarStatements(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: `
            var one = 1;
            var two = 2;
            `,
			expectedConstants: []interface{}{1, 2},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 1),
			},
		},
		{
			input: `
            var one = 1;
            one;
            `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
            var one = 1;
            var two = one;
            two;
            `,
			expectedConstants: []interface{}{1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpSetGlobal, 1),
				bytecode.Make(bytecode.OpGetGlobal, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
            var one = 1;
            one = 90;
            one;
            `,
			expectedConstants: []interface{}{1, 90},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			var num = 55;
			function() { num }
			`,
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetGlobal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			function() {
					var num = 55;
					num
			}
			`,
			expectedConstants: []interface{}{
				55,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: `
			function() {
					var a = 55;
					var b = 77;
					a + b
			}
			`,
			expectedConstants: []interface{}{
				55,
				77,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpSetLocal, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpSetLocal, 1),
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestStringExpression(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             `"Hello world!"`,
			expectedConstants: []any{"Hello world!"},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestArrayLiterals(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[]",
			expectedConstants: []interface{}{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpArray, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "[10, 8, 9]",
			expectedConstants: []interface{}{10, 8, 9},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "[1 + 2, 3 - 4, 5 * 6]",
			expectedConstants: []interface{}{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestDictionary(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "{}",
			expectedConstants: []any{},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpDic, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "{1: 2, 3: 4, 5: 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpDic, 6),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "{1: 2 + 3, 4: 5 * 6}",
			expectedConstants: []any{1, 2, 3, 4, 5, 6},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpConstant, 5),
				bytecode.Make(bytecode.OpMul),
				bytecode.Make(bytecode.OpDic, 4),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestIndexing(t *testing.T) {
	tests := []compilerTestCase{
		{
			input:             "[1, 2, 3][1 + 1]",
			expectedConstants: []any{1, 2, 3, 1, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpArray, 3),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpConstant, 4),
				bytecode.Make(bytecode.OpAdd),
				bytecode.Make(bytecode.OpIndex),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input:             "{1: 2}[2 - 1]",
			expectedConstants: []any{1, 2, 2, 1},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpDic, 2),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpSub),
				bytecode.Make(bytecode.OpIndex),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}
	testCompilerTests(t, tests)
}

func TestFunction(t *testing.T) {
	tests := []compilerTestCase{
		{
			input: "function() { return 19 + 20; }",
			expectedConstants: []any{
				19, 20,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpAdd),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "function() { 10 }",
			expectedConstants: []any{
				10,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "function() {89; 99;};",
			expectedConstants: []any{
				89,
				99,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpConstant, 1),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "function() {};",
			expectedConstants: []any{
				bytecode.Make(bytecode.OpReturn),
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "function() {90}();",
			expectedConstants: []any{
				90,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "var foo = function() {90}; foo();",
			expectedConstants: []any{
				90,
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpConstant, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpCall, 0),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "var foo = function(a) { a; }; foo(39)",
			expectedConstants: []any{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpReturnValue),
				},
				39,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpCall, 1),
				bytecode.Make(bytecode.OpPop),
			},
		},
		{
			input: "var foos = function(a, b, c) { a; b; c; }; foos(10, 20, 30);",
			expectedConstants: []any{
				[]bytecode.Instructions{
					bytecode.Make(bytecode.OpGetLocal, 0),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 1),
					bytecode.Make(bytecode.OpPop),
					bytecode.Make(bytecode.OpGetLocal, 2),
					bytecode.Make(bytecode.OpReturnValue),
				},
				10,
				20,
				30,
			},
			expectedInstructions: []bytecode.Instructions{
				bytecode.Make(bytecode.OpConstant, 0),
				bytecode.Make(bytecode.OpSetGlobal, 0),
				bytecode.Make(bytecode.OpGetGlobal, 0),
				bytecode.Make(bytecode.OpConstant, 1),
				bytecode.Make(bytecode.OpConstant, 2),
				bytecode.Make(bytecode.OpConstant, 3),
				bytecode.Make(bytecode.OpCall, 3),
				bytecode.Make(bytecode.OpPop),
			},
		},
	}

	testCompilerTests(t, tests)
}

func TestCompilerScopes(t *testing.T) {
	compiler := New()
	if compiler.scopeIndex != 0 {
		t.Errorf("scopeIndex should be 0. got=%d", compiler.scopeIndex)
	}

	globalSymbolTable := compiler.symbolTable
	compiler.emit(bytecode.OpMinus)
	compiler.enterScope()
	if compiler.scopeIndex != 1 {
		t.Errorf("scopeIndex should be1. got=%d", compiler.scopeIndex)
	}

	compiler.emit(bytecode.OpMul)
	if len(compiler.scopesStack[compiler.scopeIndex].instructions) != 1 {
		t.Errorf("instructions length wrong got=%d", len(compiler.scopesStack[compiler.scopeIndex].instructions))
	}

	last := compiler.scopesStack[compiler.scopeIndex].lastInstruction
	if last.Opcode != bytecode.OpMul {
		t.Errorf("last instructions.Opcode wrong: got=%d expected=%d", last.Opcode, bytecode.OpMul)
	}

	if compiler.symbolTable.Outer != globalSymbolTable {
		t.Errorf("outer scope did not set correctly")
	}

	compiler.leaveScope()
	if compiler.scopeIndex != 0 {
		t.Errorf("scope index should be 1 : got=%d", compiler.scopeIndex)
	}

	if compiler.symbolTable != globalSymbolTable {
		t.Errorf("compiler did not restore global symbol table")
	}
	if compiler.symbolTable.Outer != nil {
		t.Errorf("compiler modified global symbol table incorrectly")
	}

	compiler.emit(bytecode.OpAdd)
	if len(compiler.scopesStack[compiler.scopeIndex].instructions) != 2 {
		t.Errorf("should have 2 instructions: got =%d", len(compiler.scopesStack[compiler.scopeIndex].instructions))
	}

	last = compiler.scopesStack[compiler.scopeIndex].lastInstruction
	if last.Opcode != bytecode.OpAdd {
		t.Errorf("last instructions.Opcode wrong: got=%d expected=%d", last.Opcode, bytecode.OpAdd)
	}
	previous := compiler.scopesStack[compiler.scopeIndex].previousInstruction
	if previous.Opcode != bytecode.OpMinus {
		t.Errorf("last instructions.Opcode wrong: got=%d expected=%d", last.Opcode, bytecode.OpMinus)
	}
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
		case string:
			testStringObject(t, constant, actual[i])
		case []bytecode.Instructions:
			function, ok := actual[i].(*object.BytecodeFunction)
			if !ok {
				t.Errorf("constant not a function")
			}
			testInstructions(t, constant, function.Instructions)
		}
	}
}

func testStringObject(t *testing.T, constant string, act object.Object) {
	actual, ok := act.(*object.String)
	if !ok {
		t.Errorf("expecting number got=%T", act)
	}
	if constant != actual.Value {
		t.Errorf("wrong string literal got=%s expected=%s", actual.Value, constant)
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
