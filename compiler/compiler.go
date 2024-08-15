package compiler

import (
	"fmt"
	"sort"

	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/object"
)

const (
	TEMP_POSITION = 9999
)

type EmittedInstruction struct {
	Opcode   bytecode.Opcode
	Position int
}

type CompilationScope struct {
	instructions        bytecode.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Bytecode struct {
	Instructions bytecode.Instructions
	Constants    []object.Object
}

type Compiler struct {
	constants   []object.Object
	symbolTable *SymbolTable

	scopesStack []CompilationScope
	scopeIndex  int
}

// Instructions Example: [OpPop, OpConstant, 0, 3] posNewInstruction = 1

func New() *Compiler {
	globalScope := CompilationScope{
		instructions:        bytecode.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	symbolTable := NewSymbolTable()
	for i, v := range object.Builtins {
		symbolTable.DefineBuiltIn(i, v.Name)
	}
	return &Compiler{
		constants:   []object.Object{},
		symbolTable: symbolTable,
		scopesStack: []CompilationScope{globalScope},
		scopeIndex:  0,
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
		if _, ok := node.Expression.(*ast.BracketDeclaration); ok {
			break
		}
		c.emit(bytecode.OpPop)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			if err := c.Compile(s); err != nil {
				return err
			}
		}

	case *ast.FunctionDeclaration:
		c.enterScope()

		if node.Name != "" {
			c.symbolTable.DefineFunctionName(node.Name)
		}

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Literal)
		}

		if err := c.Compile(node.Body); err != nil {
			return err
		}

		if c.lastInstructionIs(bytecode.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(bytecode.OpReturnValue) {
			c.emit(bytecode.OpReturn)
		}

		freeSym := c.symbolTable.FreeSymbols
		numLocals := c.symbolTable.numberDefinitions
		instructions := c.leaveScope()

		for _, s := range freeSym {
			c.loadSymbol(s)
		}

		compiledFunc := &object.BytecodeFunction{Instructions: instructions, NumLocals: numLocals}
		c.emit(bytecode.OpClosure, c.addConstant(compiledFunc), len(freeSym))

	case *ast.CallExpression:
		if err := c.Compile(node.Function); err != nil {
			return err
		}
		for _, arg := range node.Arguments {
			if err := c.Compile(arg); err != nil {
				return err
			}
		}
		c.emit(bytecode.OpCall, len(node.Arguments))

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

	case *ast.VarStatement:
		symbol := c.symbolTable.Define(node.Variable.Literal)
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		if symbol.Scope == GlobalScope {
			c.emit(bytecode.OpSetGlobal, symbol.Index)
		} else {
			c.emit(bytecode.OpSetLocal, symbol.Index)
		}

	case *ast.Identifier:
		symbl, ok := c.symbolTable.Resolve(node.Literal)
		if !ok {
			return fmt.Errorf("variable is not defined: %s", node.Literal)
		}
		c.loadSymbol(symbl)

	case *ast.AssignmentStatement:
		symbl, ok := c.symbolTable.Resolve(node.Identifier.Literal)
		if !ok {
			return fmt.Errorf("trying to assign an undefined variable: %s", node.Identifier.Literal)
		}
		if err := c.Compile(node.Expression); err != nil {
			return err
		}
		c.emit(bytecode.OpSetGlobal, symbl.Index)

	case *ast.Number:
		number := &object.Number{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(number))

	case *ast.Null:
		c.emit(bytecode.OpNull)

	case *ast.Float:
		number := &object.Float{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(number))

	case *ast.Boolean:
		if node.Value {
			c.emit(bytecode.OpTrue)
		} else {
			c.emit(bytecode.OpFalse)
		}

	case *ast.String:
		str := &object.String{Value: node.Value}
		c.emit(bytecode.OpConstant, c.addConstant(str))

	case *ast.Array:
		for _, a := range node.Body {
			if err := c.Compile(a); err != nil {
				return err
			}
		}
		c.emit(bytecode.OpArray, len(node.Body))

	case *ast.BracketDeclaration:
		if err := c.Compile(node.Identifier); err != nil {
			return err
		}
		if err := c.Compile(node.Key); err != nil {
			return err
		}
		if err := c.Compile(node.Value); err != nil {
			return err
		}
		c.emit(bytecode.OpIndexAssign)

	case *ast.Dictionary:
		keys := []ast.Expression{}
		for k := range node.Object {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			if err := c.Compile(k); err != nil {
				return err
			}
			if err := c.Compile(node.Object[k]); err != nil {
				return err
			}
		}
		size := len(keys) * 2
		c.emit(bytecode.OpDic, size)

	case *ast.Index:
		if err := c.Compile(node.Identifier); err != nil {
			return err
		}
		if err := c.Compile(node.Index); err != nil {
			return err
		}
		c.emit(bytecode.OpIndex)

	case *ast.ReturnStatement:
		if err := c.Compile(node.ReturnExpression); err != nil {
			return err
		}
		c.emit(bytecode.OpReturnValue)

	case *ast.IFExpression:
		if err := c.Compile(node.Condition); err != nil {
			return err
		}
		jumpNotTruePos := c.emit(bytecode.OpJumpNotTrue, TEMP_POSITION)

		if err := c.Compile(node.Body); err != nil {
			return err
		}
		if c.lastInstructionIs(bytecode.OpPop) {
			c.removeLastPop()
		}
		jumpTo := c.emit(bytecode.OpJump, TEMP_POSITION)
		c.changeOperand(jumpNotTruePos, len(c.currentInstructions()))
		if node.Else == nil {
			c.emit(bytecode.OpNull)
		} else {
			if err := c.Compile(node.Else); err != nil {
				return err
			}
			if c.lastInstructionIs(bytecode.OpPop) {
				c.removeLastPop()
			}
		}
		c.changeOperand(jumpTo, len(c.currentInstructions()))
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
	previous := c.scopesStack[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopesStack[c.scopeIndex].previousInstruction = previous
	c.scopesStack[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	c.scopesStack[c.scopeIndex].instructions = append(c.currentInstructions(), ins...)
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
	op := bytecode.Opcode(c.currentInstructions()[opPos])
	c.swapInstruction(opPos, bytecode.Make(op, operand))
}

func (c *Compiler) swapInstruction(pos int, newInstruction []byte) {
	cur := c.currentInstructions()
	for i := 0; i < len(newInstruction); i++ {
		cur[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) ByteCode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) lastInstructionIs(op bytecode.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopesStack[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) currentInstructions() bytecode.Instructions {
	return c.scopesStack[c.scopeIndex].instructions
}

func (c *Compiler) removeLastPop() {
	last := c.scopesStack[c.scopeIndex].lastInstruction
	previous := c.scopesStack[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopesStack[c.scopeIndex].instructions = new
	c.scopesStack[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        bytecode.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopesStack = append(c.scopesStack, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() bytecode.Instructions {
	instructions := c.currentInstructions()

	c.scopesStack = c.scopesStack[:len(c.scopesStack)-1]
	c.scopeIndex--
	c.symbolTable = c.symbolTable.Outer
	return instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopesStack[c.scopeIndex].lastInstruction.Position
	c.swapInstruction(lastPos, bytecode.Make(bytecode.OpReturnValue))

	c.scopesStack[c.scopeIndex].lastInstruction.Opcode = bytecode.OpReturnValue
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(bytecode.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(bytecode.OpGetLocal, s.Index)
	case BuiltInScope:
		c.emit(bytecode.OpGetBuiltIn, s.Index)
	case FreeScope:
		c.emit(bytecode.OpGetFree, s.Index)
	case FunctionScope:
		c.emit(bytecode.OpCurrentClosure)
	}
}
