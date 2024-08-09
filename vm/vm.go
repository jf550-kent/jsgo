package vm

import (
	"errors"

	"github.com/jf550-kent/jsgo/bytecode"
	"github.com/jf550-kent/jsgo/compiler"
	"github.com/jf550-kent/jsgo/object"
)

const STACK_SIZE = 2048

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

		case bytecode.OpAdd:
			r, err := vm.pop()
			if err != nil {
				return err
			}
			l, err := vm.pop()
			if err != nil {
				return err
			}

			// better error handling
			left := l.(*object.Number).Value
			right := r.(*object.Number).Value

			result := left + right
			vm.push(&object.Number{Value: result})
		}
	}
	return nil
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
