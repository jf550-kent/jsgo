package main

import (
	"os"

	"github.com/jf550-kent/jsgo/evaluator"
	"github.com/jf550-kent/jsgo/parser"
)


func main() {
	b, _ := os.ReadFile("example.js")
	main := parser.Parse("", b)
	v := evaluator.Eval(main)
	print(v.String())
}