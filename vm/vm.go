package vm

import (
	"errors"
	"fmt"

	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/compiler"
	"github.com/jf550-kent/jsgo/object"
)

const STACK_SIZE = 2048
var (
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type VM struct {
	constants    []object.Object
	instructions bytecode.Instructions

	stack        []object.Object
	stackPointer int // Must always points to the new value, the object at the top of the stack is stack[stackPointer -1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,

		stack:        make([]object.Object, STACK_SIZE),
		stackPointer: 0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}

	return vm.stack[vm.stackPointer-1]
}

// [9 6 2]
func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := bytecode.Opcode(vm.instructions[ip])

		// PUSH 0000 0001 PUSH 0000 0001
		switch op {
		case bytecode.OpConstant:
			constantIndex := bytecode.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			if err := vm.push(vm.constants[constantIndex]); err != nil {
				return err
			}

		case bytecode.OpAdd, bytecode.OpSub, bytecode.OpMul, bytecode.OpDiv, bytecode.OpSHL, bytecode.OpXOR:
			if err := vm.runBinaryOperation(op); err != nil {
				return err
			}
		case bytecode.OpPop:
			vm.pop()
		case bytecode.OpTrue:
			if err := vm.push(TRUE); err != nil {
				return err
			}
		case bytecode.OpFalse:
			if err := vm.push(FALSE); err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) runBinaryOperation(op bytecode.Opcode) error {
	right, err := vm.pop()
	if err != nil {
		return err
	}
	left, err := vm.pop()
	if err != nil {
		return err
	}

	lType := left.Type()
	rType := right.Type()

	switch {
	case lType == object.NUMBER_OBJECT && rType == object.NUMBER_OBJECT:
		return vm.runNumberOperation(op, left, right)
	case lType == object.FLOAT_OBJECT && rType == object.FLOAT_OBJECT:
		return vm.runFloatOperation(op, left, right)
	case (lType == object.FLOAT_OBJECT && rType == object.NUMBER_OBJECT) || (lType == object.NUMBER_OBJECT && rType == object.FLOAT_OBJECT):
		l := object.ConvertFloat(left)
		r := object.ConvertFloat(right)
		return vm.runFloatOperation(op, l, r)
	case lType != rType:
		return fmt.Errorf("type mismatch: %s %s", left.Type(), right.Type())
	}

	return fmt.Errorf("unsupported binary operation: %s %s", left.Type(), right.Type())
}

func (vm *VM) runNumberOperation(op bytecode.Opcode, left, right object.Object) error {
	leftValue, ok := left.(*object.Number)
	if !ok {
		return fmt.Errorf("unable to assert number type for binary operation got=%T instead", left)
	}
	rightValue, ok := right.(*object.Number)
	if !ok {
		return fmt.Errorf("unable to assert number type for binary operation got=%T instead", right)
	}

	var result int64
	switch op {
	case bytecode.OpSub:
		result = leftValue.Value - rightValue.Value
	case bytecode.OpMul:
		result = leftValue.Value * rightValue.Value
	case bytecode.OpAdd:
		result = leftValue.Value + rightValue.Value
	case bytecode.OpDiv:
		if (leftValue.Value % rightValue.Value) != 0 {
			l := float64(leftValue.Value)
			r := float64(rightValue.Value)
			return vm.push(&object.Float{Value: l / r})
		}
		result = leftValue.Value / rightValue.Value
	case bytecode.OpSHL:
		result = leftValue.Value << rightValue.Value
	case bytecode.OpXOR:
		result = leftValue.Value ^ rightValue.Value
	default:
		return fmt.Errorf("unknown number operator: %d", op)
	}

	return vm.push(&object.Number{Value: result})
}

func (vm *VM) runFloatOperation(op bytecode.Opcode, left, right object.Object) error {
	leftValue, ok := left.(*object.Float)
	if !ok {
		return fmt.Errorf("unable to assert float type for binary operation got=%T instead", left)
	}
	rightValue, ok := right.(*object.Float)
	if !ok {
		return fmt.Errorf("unable to assert float type for binary operation got=%T instead", left)
	}

	var result float64
	switch op {
	case bytecode.OpSub:
		result = leftValue.Value - rightValue.Value
	case bytecode.OpMul:
		result = leftValue.Value * rightValue.Value
	case bytecode.OpAdd:
		result = leftValue.Value + rightValue.Value
	case bytecode.OpDiv:
		result = leftValue.Value / rightValue.Value
	default:
		return fmt.Errorf("unknown float operator: %d", op)
	}
	return vm.push(&object.Float{Value: result})
}

func (vm *VM) push(ob object.Object) error {
	if vm.stackPointer >= STACK_SIZE {
		return errors.New("stack overflow")
	}

	vm.stack[vm.stackPointer] = ob
	vm.stackPointer++

	return nil
}

func (vm *VM) pop() (object.Object, error) {
	if vm.stackPointer == 0 {
		return nil, errors.New("trying to pop an empty stack")
	}
	vm.stackPointer--
	return vm.stack[vm.stackPointer], nil
}

// TEST only function
func (vm *VM) lastPopStack() object.Object {
	if vm.stackPointer < 0 {
		panic("vm stack pointer is non zero")
	}
	return vm.stack[vm.stackPointer]
}
