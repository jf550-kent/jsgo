package main

import (
	"testing"
	"encoding/binary"
)


func TestInstructionsString(t *testing.T) {
	var n byte = byte(9)
	res := binary.BigEndian.Uint16([]byte{1, n, n})
	print(res)
}