package main

import (
	"fmt"
)

func man(size int) int {
	sum := 0
	byteAcc := 0
	bitNum := 0

	for y := 0; y < size; y = y + 1 {
		var ci = ((2.0 * float64(y)) / float64(size)) - 1.0

		for x := 0; x < size; x = x + 1 {
			zrzr := 0.0
			zi := 0.0
			zizi := 0.0
			cr := ((2.0 * float64(x)) / float64(size)) - 1.5

			done := false
			escape := 0
			for z := 0; !done && z < 50; z = z + 1 {
				zr := zrzr - zizi + cr
				zi = 2.0*zr*zi + ci

				zrzr = zr * zr
				zizi = zi * zi
				if zrzr+zizi > 4.0 {
					done = true
					escape = 1
				}
			}

			byteAcc = (byteAcc << 1) + escape
			bitNum = bitNum + 1

			if bitNum == 8 {
				sum = sum ^ byteAcc
				byteAcc = 0
				bitNum = 0
			}

			if x == size-1 {
				byteAcc = byteAcc << (8 - bitNum)
				sum = sum ^ byteAcc
				byteAcc = 0
				bitNum = 0
			}
		}
	}
	return sum
}

func main() {
	fmt.Println(man(5))
}
