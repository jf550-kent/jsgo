package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jf550-kent/jsgo/evaluator"
	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	WARNING = Yellow
	ERROR   = Red
	RESULT  = Green
)

func main() {
	if len(os.Args) < 2 {
		printError("Please provide file name as the first argument to be run by jsgo\n")
		printOut("usage ./jsgo <filename> -debug=<false | true >", WARNING)
		os.Exit(1)
	}
	debug := *flag.Bool("debug", false, "enable debug mode")
	fileName := os.Args[1]
	defer func() {
		if r := recover(); r != nil {
			panicMsg := fmt.Sprintf("%v\n in file : %s", r, fileName)
			printError(panicMsg)
		}
	}()

	printOut("\nWelcome to JSGO 🔥🔥🔥 🚀🚀 🔥🔥🔥 \n", RESULT)
	content, err := os.ReadFile(fileName)
	if err != nil {
		printError(err.Error())
	}

	main := parser.Parse(fileName, content)
	if main == nil {
		printError("unable to parse file : " + fileName)
	}

	result := evaluator.Eval(main)
	if err, ok := result.(*object.Error); ok {
		printError(err.Error())
	}

	if result == nil {
		printError("unable to evalulate file : " + fileName)
	}

	out := fmt.Sprintf("%+v", result)
	printOut(out, RESULT)

	if debug {
		print("in debug mode")
	}
}

func printError(out string) {
	printOut(out, ERROR)
}
func printOut(out, code string) {
	fmt.Println(code + out + Reset)
}