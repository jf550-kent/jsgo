package benchmark

import (
	"os"
	"testing"

	"github.com/jf550-kent/jsgo/compiler"
	"github.com/jf550-kent/jsgo/evaluator"
	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
	"github.com/jf550-kent/jsgo/vm"
)

func BenchmarkListTree(b *testing.B) {
	byt := setUpFile(b, "./list.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "list", byt)
	}
}

func BenchmarkListTreeDebug(b *testing.B) {
	byt := setUpFile(b, "./list.js")

	for i := 0; i < b.N; i++ {
		testEvalDebug(b, "list", byt)
	}
}

func BenchmarkListBytecode(b *testing.B) {
	byt := setUpFile(b, "./list.js")
	testBytecode(b, "list", byt)
	for i := 0; i < b.N; i++ {
		testBytecode(b, "list", byt)
	}
}

func BenchmarkTowerTree(b *testing.B) {
	byt := setUpFile(b, "./tower.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "tower", byt)
	}
}

func BenchmarkTowerTreeDebug(b *testing.B) {
	byt := setUpFile(b, "./tower.js")

	for i := 0; i < b.N; i++ {
		testEvalDebug(b, "tower", byt)
	}
}

func BenchmarkTowerBytecode(b *testing.B) {
	byt := setUpFile(b, "./tower.js")

	for i := 0; i < b.N; i++ {
		testBytecode(b, "tower", byt)
	}
}

func BenchmarkMandelbrotTree(b *testing.B) {
	src := setUpFile(b, "./mandelbrot.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "mandelbrot", src)
	}
}

func BenchmarkMandelbrotTreeDebug(b *testing.B) {
	src := setUpFile(b, "./mandelbrot.js")

	for i := 0; i < b.N; i++ {
		testEvalDebug(b, "mandelbrot", src)
	}
}

// func BenchmarkMandelbrotBytecode(b *testing.B) {
// 	src := setUpFile(b, "./mandelbrot.js")

// 	for i := 0; i < b.N; i++ {
// 		testBytecode(b, "mandelbrot", src)
// 	}
// }

func BenchmarkPermuteTree(b *testing.B) {
	src := setUpFile(b, "./permute.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "permute", src)
	}
}

func BenchmarkPermuteTreeDebug(b *testing.B) {
	src := setUpFile(b, "./permute.js")

	for i := 0; i < b.N; i++ {
		testEvalDebug(b, "permute", src)
	}
}

func BenchmarkPermuteBytecode(b *testing.B) {
	src := setUpFile(b, "./permute.js")

	for i := 0; i < b.N; i++ {
		testBytecode(b, "permute", src)
	}
}

// func BenchmarkSieveBytecode(b *testing.B) {
// 	src := setUpFile(b, "./sieve.js")

// 	for i := 0; i < b.N; i++ {
// 		testBytecode(b, "sieve", src)
// 	}
// }

func BenchmarkSieveTree(b *testing.B) {
	src := setUpFile(b, "./sieve.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "sieve", src)
	}
}

func BenchmarkSieveTreeDebug(b *testing.B) {
	src := setUpFile(b, "./sieve.js")

	for i := 0; i < b.N; i++ {
		testEvalDebug(b, "sieve", src)
	}
}

func BenchmarkQueensTree(b *testing.B) {
	src := setUpFile(b, "./queens.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "queens", src)
	}
}

func BenchmarkQueensTreeDebug(b *testing.B) {
	src := setUpFile(b, "./queens.js")

	for i := 0; i < b.N; i++ {
		testEvalDebug(b, "queens", src)
	}
}

func BenchmarkQueensBytecode(b *testing.B) {
	src := setUpFile(b, "./queens.js")

	for i := 0; i < b.N; i++ {
		testBytecode(b, "queens", src)
	}
}

func setUpFile(b *testing.B, file string) []byte {
	byt, err := os.ReadFile(file)
	if err != nil {
		b.Fatal("failed to read file", err)
	}
	b.ResetTimer()
	return byt
}

func testEval(b *testing.B, name string, src []byte) {
	main := parser.Parse(name, src)
	obj := evaluator.Eval(main, false)
	result, ok := obj.(*object.Boolean)
	if !ok {
		b.Fatal("result is not type bool")
	}
	if !result.Value {
		b.Error("result incorrect for list")
	}
}

func testEvalDebug(b *testing.B, name string, src []byte) {
	main := parser.Parse(name, src)
	obj := evaluator.Eval(main, true)
	result, ok := obj.(*object.Boolean)
	if !ok {
		b.Fatal("result is not type bool")
	}
	if !result.Value {
		b.Error("result incorrect for list")
	}
}

func testBytecode(b *testing.B, name string, src []byte) {
	main := parser.Parse(name, src)
	com := compiler.New()
	if err := com.Compile(main); err != nil {
		b.Fatalf("compiler error: %s", err)
	}

	virtualMachine := vm.New(com.ByteCode())
	if err := virtualMachine.Run(); err != nil {
		b.Fatalf("vm error: %s", err)
	}
	result := virtualMachine.LastPopStack()

	boo, ok := result.(*object.Boolean)
	if !ok {
		b.Fatalf("end result for benchmark must be boolean")
	}
	if !boo.Value {
		b.Fatalf("wrong result for %s", name)
	}
}
