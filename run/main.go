package main

import (
	"os"
)

func main() {
	byt, err := os.ReadFile("./trial.txt")
	if err != nil {
		panic(err)
	}
	print(byt)
}