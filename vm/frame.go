package vm

import (
	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/object"
)

type Frame struct {
	function    *object.BytecodeFunction
	ip          int
	basePointer int
}

func NewFrame(fn *object.BytecodeFunction, basePointer int) *Frame {
	return &Frame{
		function:    fn,
		ip:          -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instructions() bytecode.Instructions {
	return f.function.Instructions
}
