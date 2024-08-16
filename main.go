package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/jf550-kent/jsgo/compiler"
	"github.com/jf550-kent/jsgo/evaluator"
	"github.com/jf550-kent/jsgo/object"
	"github.com/jf550-kent/jsgo/parser"
	"github.com/jf550-kent/jsgo/vm"
)

const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Yellow  = "\033[33m"
	Green   = "\033[32m"
	WARNING = Yellow
	ERROR   = Red
	RESULT  = Green
	VERSION = "v2.1.0"
)

func main() {
	version := flag.Bool("version", false, "current version of JSGO")
	flag.Parse()
	if *version {
		printOut(VERSION, RESULT)
		return
	}
	if len(os.Args) < 3 {
		printError("Please provide file name as the first argument to be run by jsgo\n")
		printOut("usage ./jsgo <filename> <tree|bytecode> [debug] [-version]", WARNING)
		os.Exit(1)
	}
	fileName := os.Args[1]
	interpreter := os.Args[2]
	if interpreter != "bytecode" && interpreter != "tree" {
		log.Fatalf("flag: interpreter can only be 'tree' or 'bytecode', not: '%s'", interpreter)
	}
	debug := false
	if len(os.Args) > 3 {
		if os.Args[3] == "debug" {
			debug = true
		}
	}

	printOut("Selected interpreter:"+interpreter, RESULT)
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

	if debug {
		isJSGO := "is not"
		if Is(main) {
			isJSGO = "is"
		}
		fmt.Printf("Program: %s %s a JSGO program\n", fileName, isJSGO)
	}

	switch interpreter {
	case "tree":
		result := evaluator.Eval(main, debug)
		if err, ok := result.(*object.Error); ok {
			printError(err.Error())
		}

		if result == nil {
			printError("unable to evalulate file : " + fileName)
		}
		out := fmt.Sprintf("%+v", result)
		printOut(out, RESULT)
	case "bytecode":
		com := compiler.New()
		if debug {
			main = evaluator.Partial(main)
		}
		if err := com.Compile(main); err != nil {
			printError("compiler error: " + err.Error())
		}

		virtualMachine := vm.New(com.ByteCode())
		if err := virtualMachine.Run(); err != nil {
			printError("vm error: " + err.Error())
			break
		}

		result := virtualMachine.LastPopStack()
		out := fmt.Sprintf("%+v", result)
		printOut(out, RESULT)
	}
}

func printError(out string) {
	printOut(out, ERROR)
}
func printOut(out, code string) {
	fmt.Println(code + out + Reset)
}
