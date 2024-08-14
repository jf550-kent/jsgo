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
	MAX_FRAMES           = 1024
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

type VM struct {
	constants []object.Object
	globals   []object.Object

	stack        []object.Object
	stackPointer int // Must always points to the new value, the object at the top of the stack is stack[stackPointer -1]

	frames      []*Frame
	framesIndex int
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.BytecodeFunction{Instructions: bytecode.Instructions}
	mainClosure := &object.Closure{Fn: mainFn}
	mainFrame := NewFrame(mainClosure, 0)

	frames := make([]*Frame, MAX_FRAMES)
	frames[0] = mainFrame

	return &VM{
		constants: bytecode.Constants,
		globals:   make([]object.Object, MAX_GLOBAL_VARIABLES),

		stack:        make([]object.Object, STACK_SIZE),
		stackPointer: 0,

		frames:      frames,
		framesIndex: 1,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.stackPointer == 0 {
		return nil
	}

	return vm.stack[vm.stackPointer-1]
}

func (vm *VM) Run() error {
	var ip int
	var ins bytecode.Instructions
	var op bytecode.Opcode

	for vm.currentFrame().ip < len(vm.currentFrame().Instructions())-1 {
		vm.currentFrame().ip++

		ip = vm.currentFrame().ip
		ins = vm.currentFrame().Instructions()
		op = bytecode.Opcode(ins[ip])

		switch op {
		case bytecode.OpConstant:
			constantIndex := bytecode.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
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
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			result, err := vm.pop()
			if err != nil {
				return err
			}
			if !isTruthy(result) {
				vm.currentFrame().ip = pos - 1
			}
		case bytecode.OpJump:
			// [OpJump 0, 3, OpConstant 0, 9]
			pos := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip = pos - 1
		case bytecode.OpNull:
			if err := vm.push(NULL); err != nil {
				return err
			}
		case bytecode.OpSetGlobal:
			globalIndex := bytecode.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			value, err := vm.pop()
			if err != nil {
				return err
			}
			vm.globals[globalIndex] = value
		case bytecode.OpArray:
			size := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			start := vm.stackPointer - size
			array := vm.makeArray(start, vm.stackPointer, size)
			vm.stackPointer = start

			if err := vm.push(array); err != nil {
				return err
			}
		case bytecode.OpIndexAssign:
			ident := vm.stack[vm.stackPointer-3]
			index, ok := vm.stack[vm.stackPointer-2].(*object.Number)
			if !ok {
				return fmt.Errorf("wrong type for array index")
			}
			expr := vm.stack[vm.stackPointer-1]

			switch val := ident.(type) {
			case *object.Array:
				if int(index.Value) >= len(val.Body) {
					newArr := make([]object.Object, int(index.Value)+1)
					copy(newArr, val.Body)
					val.Body = newArr
				}
				val.Body[index.Value] = expr
			default:
				return fmt.Errorf("cannot index with type=%v", ident)
			}
		case bytecode.OpDic:
			size := int(bytecode.ReadUint16(ins[ip+1:]))
			vm.currentFrame().ip += 2
			start := vm.stackPointer - size
			dictionary, err := vm.makeDictionary(start, vm.stackPointer)
			if err != nil {
				return err
			}
			vm.stackPointer = start
			if err := vm.push(dictionary); err != nil {
				return err
			}
		case bytecode.OpIndex:
			index, err := vm.pop()
			if err != nil {
				return err
			}
			identifier, err := vm.pop()
			if err != nil {
				return err
			}

			if err := vm.runIndexExpression(identifier, index); err != nil {
				return err
			}

		case bytecode.OpGetGlobal:
			globalIndex := bytecode.ReadUint16(ins[ip+1:])
			vm.currentFrame().ip += 2
			value := vm.globals[globalIndex]
			if value == nil {
				return fmt.Errorf("variable not defined")
			}

			if err := vm.push(value); err != nil {
				return err
			}
		case bytecode.OpCall:
			args := int(bytecode.ReadUnit8(ins[ip+1:]))
			vm.currentFrame().ip += 1

			if err := vm.runCall(args); err != nil {
				return err
			}
		case bytecode.OpSetLocal:
			localIndex := bytecode.ReadUnit8(ins[ip+1:])
			vm.currentFrame().ip += 1
			frame := vm.currentFrame()
			val, err := vm.pop()
			if err != nil {
				return err
			}
			vm.stack[frame.basePointer+int(localIndex)] = val
		case bytecode.OpGetLocal:
			localIndex := bytecode.ReadUnit8(ins[ip+1:])
			vm.currentFrame().ip += 1

			frame := vm.currentFrame()

			if err := vm.push(vm.stack[frame.basePointer+int(localIndex)]); err != nil {
				return err
			}
		case bytecode.OpReturnValue:
			returnValue, err := vm.pop()
			if err != nil {
				return err
			}
			frame := vm.popFrame()
			vm.stackPointer = frame.basePointer
			vm.pop()

			if err := vm.push(returnValue); err != nil {
				return err
			}
		case bytecode.OpReturn:
			frame := vm.popFrame()
			vm.stackPointer = frame.basePointer
			vm.pop()

			if err := vm.push(NULL); err != nil {
				return err
			}
		case bytecode.OpGetBuiltIn:
			index := bytecode.ReadUnit8(ins[ip+1:])
			vm.currentFrame().ip += 1

			def := object.Builtins[index]
			err := vm.push(def)
			if err != nil {
				return err
			}
		case bytecode.OpClosure:
			index := int(bytecode.ReadUint16(ins[ip+1:]))
			free := int(bytecode.ReadUnit8(ins[ip+3:]))

			vm.currentFrame().ip += 3
			if err := vm.pushClosure(index, free); err != nil {
				return err
			}
		case bytecode.OpGetFree:
			freeIndex := bytecode.ReadUnit8(ins[ip+1:])
			vm.currentFrame().ip += 1

			currentClosure := vm.currentFrame().function
			if err := vm.push(currentClosure.Free[freeIndex]); err != nil {
				return err
			}
		case bytecode.OpCurrentClosure:
			currClosure := vm.currentFrame().function
			if err := vm.push(currClosure); err != nil {
				return err
			}
		}
	}
	return nil
}

func (vm *VM) currentFrame() *Frame {
	return vm.frames[vm.framesIndex-1]
}

func (vm *VM) pushClosure(index, free int) error {
	constant := vm.constants[index]
	function, ok := constant.(*object.BytecodeFunction)
	if !ok {
		return fmt.Errorf("not a function: %+v", constant)
	}

	freeVar := make([]object.Object, free)
	for i := 0; i < free; i++ {
		freeVar[i] = vm.stack[vm.stackPointer-free+i]
	}
	vm.stackPointer = vm.stackPointer - free

	closure := &object.Closure{Fn: function, Free: freeVar}
	return vm.push(closure)
}

func (vm *VM) pushFrame(f *Frame) {
	if vm.framesIndex > MAX_FRAMES {
		panic("maximum frame reached")
	}
	vm.frames[vm.framesIndex] = f
	vm.framesIndex++
}

func (vm *VM) popFrame() *Frame {
	vm.framesIndex--
	return vm.frames[vm.framesIndex]
}

func (vm *VM) runIndexExpression(identifier, index object.Object) error {
	identifierType := identifier.Type()
	indexType := index.Type()
	switch {
	case identifierType == object.ARRAY_OBJECT && indexType == object.NUMBER_OBJECT:
		return vm.runArrayIndex(identifier, index)
	case identifierType == object.DICTIONARY_OBJECT && indexType == object.NUMBER_OBJECT:
		return vm.runDictionaryIndex(identifier, index)
	}

	return fmt.Errorf("index operation not supported for %s[%s]", identifierType, indexType)
}

func (vm *VM) runArrayIndex(identifier, index object.Object) error {
	arrayObj, ok := identifier.(*object.Array)
	if !ok {
		return fmt.Errorf("not array object passed to index array")
	}
	num, ok := index.(*object.Number)
	if !ok {
		return fmt.Errorf("non number is used to index array")
	}
	numIdex := num.Value
	max := int64(len(arrayObj.Body) - 1)

	if numIdex < 0 || numIdex > max {
		return vm.push(NULL)
	}
	return vm.push(arrayObj.Body[numIdex])
}

func (vm *VM) runCall(numArgs int) error {
	caller := vm.stack[vm.stackPointer-1-numArgs]
	switch caller := caller.(type) {
	case *object.Closure:
		return vm.callClosure(caller, numArgs)
	case *object.BuiltIn:
		return vm.callBuiltin(caller, numArgs)
	}
	return fmt.Errorf("calling non-function and non-built-in")
}

func (vm *VM) callClosure(fn *object.Closure, numArgs int) error {

	frame := NewFrame(fn, vm.stackPointer-numArgs)
	vm.pushFrame(frame)

	vm.stackPointer = frame.basePointer + fn.Fn.NumLocals

	return nil
}

func (vm *VM) callBuiltin(builtin *object.BuiltIn, numArgs int) error {
	args := vm.stack[vm.stackPointer-numArgs : vm.stackPointer]

	result := builtin.Function(args...)
	vm.stackPointer = vm.stackPointer - numArgs - 1

	if result != nil {
		vm.push(result)
	} else {
		vm.push(NULL)
	}

	return nil
}

func (vm *VM) runDictionaryIndex(identifier, index object.Object) error {
	dic, ok := identifier.(*object.Dictionary)
	if !ok {
		return fmt.Errorf("identifier is not an dictionary for indexing")
	}

	key, ok := index.(object.Hasher)
	if !ok {
		return fmt.Errorf("unable to hash key for indexing dictionary")
	}

	keyValue, ok := dic.Value[key.Hash()]
	if !ok {
		return vm.push(NULL)
	}
	return vm.push(keyValue.Value)
}

func (vm *VM) makeDictionary(start, end int) (object.Object, error) {
	dic := make(map[object.Hash]object.KeyValue)

	for i := start; i < end; i += 2 {
		key := vm.stack[i]
		value := vm.stack[i+1]
		hash, ok := key.(object.Hasher)

		if !ok {
			return nil, fmt.Errorf("dictionary key unhashbale: %s", key.String())
		}

		keyHash := hash.Hash()
		dic[keyHash] = object.KeyValue{Key: key, Value: value}
	}
	return &object.Dictionary{Value: dic}, nil
}

func (vm *VM) makeArray(start, end, size int) object.Object {
	arr := make([]object.Object, size)

	for i := start; i < end; i++ {
		arr[i-start] = vm.stack[i]
	}
	return &object.Array{Body: arr}
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
