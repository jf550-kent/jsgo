package bytecode

import "encoding/binary"

type Opcode byte

type Instructions []byte

const (
	OpConstant Opcode = iota
)

type Definition struct {
	Name        string
	OperandSize []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstant", []int{2}},
}

// OpConstant : [8, 8] -> operand can be 16 bits

func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	// optimization -> use a field value instead of iterations
	for _, w := range def.OperandSize {
		instructionLen += w
	}

	instruction := make([]byte, instructionLen)
	instruction[0] = byte(op)

	offset := 1
	for i, o := range operands {
		size := def.OperandSize[i]
		switch size {
			case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		}
		offset += size
	}

	return instruction
}
