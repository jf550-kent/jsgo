package vm

import (
	"errors"
	"fmt"

	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/compiler"
	"github.com/jf550-kent/jsgo/object"
)

const (
	STACK_SIZE           = 2048
	MAX_GLOBAL_VARIABLES = 65536
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

type VM struct {
	constants    []object.Object
	instructions bytecode.Instructions
	globals      []object.Object

	stack        []object.Object
	stackPointer int // Must always points to the new value, the object at the top of the stack is stack[stackPointer -1]
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		globals:      make([]object.Object, MAX_GLOBAL_VARIABLES),

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

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := bytecode.Opcode(vm.instructions[ip])

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
		case bytecode.OpEqual, bytecode.OpGreaterThan, bytecode.OpNotEqual:
			if err := vm.runComparison(op); err != nil {
				return err
			}
		case bytecode.OpBang:
			if err := vm.runBang(); err != nil {
				return err
			}
		case bytecode.OpMinus:
			if err := vm.runMinus(); err != nil {
				return err
			}
		case bytecode.OpJumpNotTrue:
			pos := int(bytecode.ReadUint16(vm.instructions[ip+1:]))
			ip += 2
			result, err := vm.pop()
			if err != nil {
				return err
			}
			if !isTruthy(result) {
				ip = pos - 1
			}
		case bytecode.OpJump:
			// [OpJump 0, 3, OpConstant 0, 9]
			pos := int(bytecode.ReadUint16(vm.instructions[ip+1:]))
			ip = pos - 1
		case bytecode.OpNull:
			if err := vm.push(NULL); err != nil {
				return err
			}
		case bytecode.OpSetGlobal:
			globalIndex := bytecode.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			value, err := vm.pop()
			if err != nil {
				return err
			}
			vm.globals[globalIndex] = value
		case bytecode.OpGetGlobal:
			globalIndex := bytecode.ReadUint16(vm.instructions[ip+1:])
			ip += 2
			value := vm.globals[globalIndex]
			if value == nil {
				return fmt.Errorf("variable not defined")
			}

			if err := vm.push(value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) runBinaryOperation(op bytecode.Opcode) error {
	left, right, err := vm.popLeftRight()
	if err != nil {
		return err
	}

	switch {
	case left.Type() == object.NUMBER_OBJECT || left.Type() == object.FLOAT_OBJECT:
		left, right, isNumber, err := vm.checkNumberType(left, right)
		if err != nil {
			break
		}
		if isNumber {
			return vm.runNumberOperation(op, left, right)
		}
		return vm.runFloatOperation(op, left, right)
	}

	return fmt.Errorf("unsupported binary operation: %s %s", left.Type(), right.Type())
}

func (vm *VM) checkNumberType(left, right object.Object) (object.Object, object.Object, bool, error) {
	lType := left.Type()
	rType := right.Type()

	switch {
	case lType == object.NUMBER_OBJECT && rType == object.NUMBER_OBJECT:
		return left, right, true, nil
	case lType == object.FLOAT_OBJECT && rType == object.FLOAT_OBJECT:
		return left, right, false, nil
	case (lType == object.FLOAT_OBJECT && rType == object.NUMBER_OBJECT) || (lType == object.NUMBER_OBJECT && rType == object.FLOAT_OBJECT):
		l := object.ConvertFloat(left)
		r := object.ConvertFloat(right)
		return l, r, false, nil
	}

	return left, right, false, fmt.Errorf("left and right is not a number or float")
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

func (vm *VM) runComparison(op bytecode.Opcode) error {
	left, right, err := vm.popLeftRight()
	if err != nil {
		return err
	}

	if l, r, isNumber, err := vm.checkNumberType(left, right); err == nil {
		if isNumber {
			return vm.compareNumber(op, l, r)
		}
		return vm.compareFloat(op, l, r)
	} else {
		if err.Error() != "left and right is not a number or float" {
			return err
		}
	}

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBool(right == left))
	case bytecode.OpNotEqual:
		return vm.push(nativeBool(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (vm *VM) compareNumber(op bytecode.Opcode, left, right object.Object) error {
	leftValue, ok := left.(*object.Number)
	if !ok {
		return fmt.Errorf("unable to assert number type for compare operation got=%T instead", left)
	}
	rightValue, ok := right.(*object.Number)
	if !ok {
		return fmt.Errorf("unable to assert number type for compare operation got=%T instead", right)
	}

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBool(rightValue.Value == leftValue.Value))
	case bytecode.OpNotEqual:
		return vm.push(nativeBool(rightValue.Value != leftValue.Value))
	case bytecode.OpGreaterThan:
		return vm.push(nativeBool(leftValue.Value > rightValue.Value))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) compareFloat(op bytecode.Opcode, left, right object.Object) error {
	leftValue, ok := left.(*object.Float)
	if !ok {
		return fmt.Errorf("unable to assert float type for binary operation got=%T instead", left)
	}
	rightValue, ok := right.(*object.Float)
	if !ok {
		return fmt.Errorf("unable to assert float type for binary operation got=%T instead", left)
	}

	switch op {
	case bytecode.OpEqual:
		return vm.push(nativeBool(rightValue.Value == leftValue.Value))
	case bytecode.OpNotEqual:
		return vm.push(nativeBool(rightValue.Value != leftValue.Value))
	case bytecode.OpGreaterThan:
		return vm.push(nativeBool(leftValue.Value > rightValue.Value))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func (vm *VM) runBang() error {
	operand, err := vm.pop()
	if err != nil {
		return err
	}

	switch operand {
	case TRUE:
		return vm.push(FALSE)
	case FALSE:
		return vm.push(TRUE)
	case NULL:
		return vm.push(TRUE)
	}
	return vm.push(FALSE)
}

func (vm *VM) runMinus() error {
	operand, err := vm.pop()
	if err != nil {
		return err
	}

	switch r := operand.(type) {
	case *object.Number:
		return vm.push(&object.Number{Value: -r.Value})
	case *object.Float:
		return vm.push(&object.Float{Value: -r.Value})
	}

	return fmt.Errorf("unsupported type for negation: %s", operand.Type())
}

func (vm *VM) popLeftRight() (object.Object, object.Object, error) {
	right, err := vm.pop()
	if err != nil {
		return nil, nil, err
	}
	left, err := vm.pop()
	if err != nil {
		return nil, nil, err
	}

	return left, right, nil
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

func nativeBool(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {

	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}
