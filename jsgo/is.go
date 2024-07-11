package main

import "github.com/jun-hf/jsgo/ast"


// if I design Is such that it reports errors
func Is(main ast.Main) bool {
	if len(main.Statements) != 0 {
		for _, stmt := range main.Statements {
			if ok := isStatement(stmt); !ok {
				return false
			}
		}
	}
	return true 
}

func isStatement(stmt ast.Statement) bool {
	switch st := stmt.(type) {
	case *ast.VarStatement:
		print(st)
		return true
	}
	return false
}

func main() {
}