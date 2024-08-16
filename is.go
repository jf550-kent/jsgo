package main

import (
	"fmt"

	"github.com/jf550-kent/jsgo/ast"
)

func Is(node *ast.Main) bool {
	for _, stmt := range node.Statements {
		if !check(stmt) {
			return false
		}
	}
	return true
}

func check(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.Number, *ast.Float, *ast.Boolean, *ast.Null, *ast.String, *ast.Array, *ast.Dictionary, *ast.Identifier:
		return true
	case *ast.BinaryExpression:
		switch node.Operator {
		case "+", "-", "*", "/", "<<", "^", "<", ">", "==", "!=":
			return check(node.Left) && check(node.Right)
		default:
			return false
		}
	case *ast.IFExpression:
		correct := checkBlockStatements(node.Body)
		correct = check(node.Condition) && correct
		if node.Else != nil { correct = correct && checkBlockStatements(node.Else)}
		return correct
	case *ast.UnaryExpression:
		if node.Operator != "!" && node.Operator != "-" {
			return false
		}
		return check(node.Expression)
	case *ast.FunctionDeclaration:
		return checkBlockStatements(node.Body)
	case *ast.Index:
		return check(node.Identifier) && check(node.Index)
	case *ast.CallExpression:
		if !check(node.Function) {
			return false
		}

		if len(node.Arguments) != 0 {
			for _, arg := range node.Arguments {
				if !check(arg) {
					return false
				}
			}
		}

		return true
	case *ast.BlockStatement:
		return checkBlockStatements(node)
	case *ast.ExpressionStatement:
		return check(node.Expression)
	case *ast.ReturnStatement:
		return check(node.ReturnExpression)
	case *ast.ForStatement:
		return check(node.Condition)
	case *ast.VarStatement:
		return check(node.Expression)
	case *ast.AssignmentStatement:
		return check(node.Identifier) && check(node.Expression)
	}
	return true
}

func checkBlockStatements(stmts *ast.BlockStatement) bool {
	for _, stmt := range stmts.Statements {
		if !check(stmt) {
			fmt.Println(stmt)
			return false
		}
	}

	return true
}

// Binary operation allows:
// number [+|-|*|/|<<|^|<|>|==|!=] number 
// float [+|-|*|/|<|>|==|!=] float
// float [+|-|*|/|<|>|==|!=] number = float [+|-|*|/|<|>|==|!=] float
// <expression> [!= | == ] <expression>