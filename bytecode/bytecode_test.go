package bytecode

import "testing"

func TestMake(t *testing.T) {
	tests := []struct {
		op       Opcode
		operands []int
		expected []byte
	}{
		{OpConstant, []int{65534}, []byte{byte(OpConstant), 255, 254}},
		{OpGetLocal, []int{255}, []byte{byte(OpGetLocal), 255}},
		{OpClosure, []int{65534, 255}, []byte{byte(OpClosure), 255, 254, 255}},
	}

	for _, tt := range tests {
		instruction := Make(tt.op, tt.operands...)

		if len(instruction) != len(tt.expected) {
			t.Errorf("instruction has wrong length. want=%d, got=%d", len(tt.expected), len(instruction))
		}

		for i, b := range tt.expected {
			if instruction[i] != tt.expected[i] {
				t.Errorf("wrong byte at pos %d. want=%d, got=%d",
					i, b, instruction[i])
			}
		}
	}
}

func TestInstructionsString(t *testing.T) {
	intrustions := []Instructions{
		Make(OpConstant, 1),
		Make(OpConstant, 2),
		Make(OpConstant, 65535),
		Make(OpAdd),
		Make(OpGetLocal, 1),
		Make(OpClosure, 65535, 255),
	}

	expected := "0000 OpConstant 1\n0003 OpConstant 2\n0006 OpConstant 65535\n0009 OpAdd\n0010 OpGetLocal 1\n0012 OpClosure 65535 255\n"
	merged := Instructions{}
	for _, ins := range intrustions {
		merged = append(merged, ins...)
	}

	if merged.String() != expected {
		t.Errorf("instructions wrongly formatted\ngot:\n%s \nexpected:\n%s", merged.String(), expected)
	}
}

func TestReadOperands(t *testing.T) {
	tests := []struct {
		op        Opcode
		operands  []int
		bytesRead int
	}{
		{OpConstant, []int{65535}, 2},
		{OpGetLocal, []int{255}, 1},
		{OpClosure, []int{65535, 255}, 3},
	}

	for _, tt := range tests {
		instruct := Make(tt.op, tt.operands...)

		def, err := Lookup(byte(tt.op))
		if err != nil {
			t.Fatalf("definition not found: %s", err)
		}

		operandsRead, n := ReadOperands(def, instruct[1:])
		if n != tt.bytesRead {
			t.Fatalf("n wrong. got=%d, expected=%d", n, tt.bytesRead)
		}

		for i, want := range tt.operands {
			if operandsRead[i] != want {
				t.Errorf("operand wrong: got=%d, expected=%d", operandsRead[i], want)
			}
		}
	}
}
