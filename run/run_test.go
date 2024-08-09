package main

import (
	"testing"
)


func TestInstructionsString(t *testing.T) {
	var shl = 20 << 10
	var xor = 99 ^ 8


	print(shl, xor)
}