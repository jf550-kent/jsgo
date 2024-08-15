package bytecode

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Opcode byte

type Instructions []byte

func (ins Instructions) String() string {
	var res strings.Builder

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&res, "ERROR: %s\n", err)
			continue
		}

		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&res, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		i += 1 + read
	}

	return res.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := def.OperandNum

	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not match defined %d\n",
			len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

const (
	OpConstant Opcode = iota
	OpAdd
	OpPop
	OpSub
	OpMul
	OpDiv
	OpSHL
	OpXOR
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTrue
	OpJump
	OpNull
	OpGetGlobal
	OpSetGlobal
	OpArray
	OpDic
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpGetLocal
	OpSetLocal
	OpGetBuiltIn
	OpClosure
	OpGetFree
	OpCurrentClosure
	OpIndexAssign
	OpFor
)

type Definition struct {
	Name         string
	OperandWidth []int
	ByteSize     int
	OperandNum   int
}

var definitions = map[Opcode]*Definition{
	OpConstant:       {"OpConstant", []int{2}, 2, 1},
	OpAdd:            {"OpAdd", []int{}, 0, 0},
	OpPop:            {"OpPop", []int{}, 0, 0},
	OpSub:            {"OpSub", []int{}, 0, 0},
	OpMul:            {"OpMul", []int{}, 0, 0},
	OpDiv:            {"OpDiv", []int{}, 0, 0},
	OpSHL:            {"OpSHL", []int{}, 0, 0},
	OpXOR:            {"OpXOR", []int{}, 0, 0},
	OpTrue:           {"OpTrue", []int{}, 0, 0},
	OpFalse:          {"OpFalse", []int{}, 0, 0},
	OpEqual:          {"OpEqual", []int{}, 0, 0},
	OpNotEqual:       {"OpNotEqual", []int{}, 0, 0},
	OpGreaterThan:    {"OpGreaterThan", []int{}, 0, 0},
	OpMinus:          {"OpMinus", []int{}, 0, 0},
	OpBang:           {"OpBang", []int{}, 0, 0},
	OpJumpNotTrue:    {"OpJumpNotTrue", []int{2}, 2, 1},
	OpJump:           {"OpJump", []int{2}, 2, 1},
	OpNull:           {"OpNull", []int{}, 0, 0},
	OpGetGlobal:      {"OpGetGlobal", []int{2}, 2, 1},
	OpSetGlobal:      {"OpSetGlobal", []int{2}, 2, 1},
	OpArray:          {"OpArray", []int{2}, 2, 1},
	OpDic:            {"OpDic", []int{2}, 2, 1},
	OpIndex:          {"OpIndex", []int{}, 0, 0},
	OpCall:           {"OpCall", []int{1}, 1, 1},
	OpReturnValue:    {"OpReturnValue", []int{}, 0, 0},
	OpReturn:         {"OpReturn", []int{}, 0, 0},
	OpGetLocal:       {"OpGetLocal", []int{1}, 1, 1},
	OpSetLocal:       {"OpSetLocal", []int{1}, 1, 1},
	OpGetBuiltIn:     {"OpGetBuiltIn", []int{1}, 1, 1},
	OpClosure:        {"OpClosure", []int{2, 1}, 3, 2},
	OpGetFree:        {"OpGetFree", []int{1}, 1, 1},
	OpCurrentClosure: {"OpCurrentClosure", []int{}, 0, 0},
	OpIndexAssign:    {"OpIndexAssign", []int{}, 0, 0},
	OpFor:            {"OpFor", []int{1, 1, 1}, 3, 3},
}

func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}

	return def, nil
}

// Make takes an opcode and abitrary number of operand and return an Instruction
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1 + def.ByteSize
	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		size := def.OperandWidth[i]
		switch size {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += size
	}

	return instruction
}

// [0000 0001, 0100 0101] -> [1, 69]
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, def.OperandNum)
	offset := 0

	for i, size := range def.OperandWidth {
		switch size {
		case 2:
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1:
			operands[i] = int(ReadUnit8(ins[offset:]))
		}

		offset += size
	}

	return operands, offset
}

func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

func ReadUnit8(ins Instructions) uint8 {
	return uint8(ins[0])
}
