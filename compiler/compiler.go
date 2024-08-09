package compiler

import (
	"fmt"

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
	switch node := node.(type) {
	case *ast.Main:
		return c.compileMain(node)
	case *ast.ExpressionStatement:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(bytecode.OpPop)
	case *ast.BinaryExpression:
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(bytecode.OpAdd)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.Number:
		number := &object.Number{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(number))
	}

	return nil
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) emit(op bytecode.Opcode, operands ...int) int {
	instrct := bytecode.Make(op, operands...)
	pos := c.addInstruction(instrct)
	return pos
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.instructions)
	c.instructions = append(c.instructions, ins...)
	return posNewInstruction
}

func (c *Compiler) compileMain(node *ast.Main) error {
	for _, stmt := range node.Statements {
		err := c.Compile(stmt)
		if err != nil {
			return err
		}
	}
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
