package compiler

import (
	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/object"
)

type Compiler struct {
	instructions bytecode.Instructions
	constants    []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constants:    []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

type Bytecode struct {
	Instructions bytecode.Instructions
	Constants    []object.Object
}

func (c *Compiler) ByteCode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}
