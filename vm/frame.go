package vm

import (
	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/object"
)

type Frame struct {
	function *object.BytecodeFunction
	instructionPointer int
	basePointer int
}

func NewFrame(fn *object.BytecodeFunction, basePointer int) *Frame {
	return &Frame{
		function: fn,
		instructionPointer: -1,
		basePointer: basePointer,
	}
}

func (f *Frame) Instruction() bytecode.Instructions {
	return f.function.Instructions
}