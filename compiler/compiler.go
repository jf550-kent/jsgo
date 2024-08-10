package compiler

import (
	"fmt"

	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/object"
)

const (
	TEMP_POSITION = 9999
)

type EmittedInstruction struct {
	Opcode bytecode.Opcode
	Position int
}

type Compiler struct {
	instructions bytecode.Instructions
	constants    []object.Object

	lastInstruction EmittedInstruction
	previousInstruction EmittedInstruction
}

// Instructions Example: [OpPop, OpConstant, 0, 3] posNewInstruction = 1

func New() *Compiler {
	return &Compiler{
		instructions: bytecode.Instructions{},
		constants:    []object.Object{},
		lastInstruction: EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
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
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
					return err
			}
		}
	case *ast.BinaryExpression:
		if node.Operator == "<" {
			return c.compileLessThan(node)
		}
		if err := c.Compile(node.Left); err != nil {
			return err
		}
		if err := c.Compile(node.Right); err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(bytecode.OpAdd)
		case "-":
			c.emit(bytecode.OpSub)
		case "*":
			c.emit(bytecode.OpMul)
		case "/":
			c.emit(bytecode.OpDiv)
		case "<<":
			c.emit(bytecode.OpSHL)
		case "^":
			c.emit(bytecode.OpXOR)
		case ">":
			c.emit(bytecode.OpGreaterThan)
		case "==":
			c.emit(bytecode.OpEqual)
		case "!=":
			c.emit(bytecode.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.UnaryExpression:
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(bytecode.OpBang)
		case "-":
			c.emit(bytecode.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.Number:
		number := &object.Number{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(number))
	case *ast.Float:
		number := &object.Float{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(number))
	case *ast.Boolean:
		if node.Value {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}
	case *ast.IFExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jumpNotTruePos := c.emit(bytecode.OpJumpNotTrue, TEMP_POSITION)

		if err := c.Compile(node.Body); err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}
		if node.Else == nil {
			return c.changeOpJumpNotTrueOperand(jumpNotTruePos)
		}
		jumpTo := c.emit(bytecode.OpJump, TEMP_POSITION)
		c.changeOpJumpNotTrueOperand(jumpNotTruePos)
		if err := c.Compile(node.Else); err != nil {
			return err
		}
		if c.lastInstructionIsPop() {
			c.removeLastPop()
		}
		c.changeOperand(jumpTo, len(c.instructions))
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

	c.setLastInstruction(op, pos)
	return pos
}

func (c *Compiler) setLastInstruction(op bytecode.Opcode, pos int) {
	previous := c.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.previousInstruction = previous
	c.lastInstruction = last
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

func (c *Compiler) compileLessThan(node *ast.BinaryExpression) error {
	if node.Operator != "<" {
		panic("only use compileLSS for < operator")
	}

	if err := c.Compile(node.Right); err != nil {
		return err
	}

	if err := c.Compile(node.Left); err != nil {
		return err
	}

	c.emit(bytecode.OpGreaterThan)
	return nil
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := bytecode.Opcode(c.instructions[opPos])
	c.swapInstruction(opPos, bytecode.Make(op, operand))
}

func (c *Compiler) swapInstruction(pos int, newInstruction []byte) {
	for i := 0; i < len(newInstruction); i++ {
		c.instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOpJumpNotTrueOperand(pos int) error {
	last := len(c.instructions)
	instrt := bytecode.Make(bytecode.OpJumpNotTrue, last)

	for _, b := range instrt {
		c.instructions[pos] = b
		pos++
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

func (c *Compiler) lastInstructionIsPop() bool {
	return c.lastInstruction.Opcode == bytecode.OpPop
}

func (c *Compiler) removeLastPop() {
	c.instructions = c.instructions[:c.lastInstruction.Position]
	c.lastInstruction = c.previousInstruction
}