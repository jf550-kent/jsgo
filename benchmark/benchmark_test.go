package benchmark

import (
	"os"
	"testing"

	"github.com/jf550-kent/jsgo/evaluator"
	"github.com/jf550-kent/jsgo/parser"
	"github.com/jf550-kent/jsgo/object"
)

func BenchmarkList(b *testing.B) {
	byt := setUpFile(b, "./list.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "list", byt)
	}
}

func BenchmarkTower(b *testing.B) {
	byt := setUpFile(b, "./tower.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "tower", byt)
	}
}

func BenchmarkMandelbrot(b *testing.B) {
	src := setUpFile(b, "./mandelbrot.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "mandelbrot", src)
	}
}

func BenchmarkPermute(b *testing.B) {
	src := setUpFile(b, "./permute.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "permute", src)
	}
}

func BenchmarkSieve(b *testing.B) {
	src := setUpFile(b, "./sieve.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "sieve", src)
	}
}

func BenchmarkQueens(b *testing.B) {
	src := setUpFile(b, "./queens.js")

	for i := 0; i < b.N; i++ {
		testEval(b, "queens", src)
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
	obj := evaluator.Eval(main)
	result, ok := obj.(*object.Boolean)
	if !ok {
		b.Fatal("result is not type bool")
	}
	if !result.Value {
		b.Error("result incorrect for list")
	}
}