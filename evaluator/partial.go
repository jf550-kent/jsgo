package evaluator

import (
	"fmt"

	"github.com/jf550-kent/jsgo/ast"
)

func Partial(main *ast.Main) *ast.Main {
	stmts := make([]ast.Statement, len(main.Statements))
	for i, stmt := range main.Statements {
		st := partialEvalStatement(stmt)
		stmts[i] = st
	}
	return &ast.Main{Statements: stmts}
}

func partialEvalStatement(stmt ast.Statement) ast.Statement {
	if stmt == nil {
		panic(fmt.Errorf("nil passed into partialEvalStatement"))
	}
	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		e := partialEvalExpression(s.Expression)
		s.Expression = e
		return s
	case *ast.VarStatement:
		e := partialEvalExpression(s.Expression)
		s.Expression = e
		return s
	case *ast.AssignmentStatement:
		e := partialEvalExpression(s.Expression)
		s.Expression = e
		return s
	case *ast.ReturnStatement:
		e := partialEvalExpression(s.ReturnExpression)
		s.ReturnExpression = e
		return s
	}
	return stmt
}

func partialEvalExpression(exp ast.Expression) ast.Expression {
	if exp == nil {
		panic("nil passsed to partialEvalExpression")
	}
	switch e := exp.(type) {
	case *ast.Number, *ast.Float, *ast.Boolean, *ast.String, *ast.Null:
		return e
	case *ast.BinaryExpression:
		return partialEvalBinaryOperation(e)
	case *ast.UnaryExpression:
		return partialEvalUnaryOperation(e)
	case *ast.FunctionDeclaration:
		for i, stmt := range e.Body.Statements {
			st := partialEvalStatement(stmt)
			e.Body.Statements[i] = st
		}
		return e
	case *ast.CallExpression:
		for i, callE := range e.Arguments {
			e.Arguments[i] = partialEvalExpression(callE)
		}
		return e
	case *ast.Array:
		for i, arrEle := range e.Body {
			e.Body[i] = partialEvalExpression(arrEle)
		}
		return e
	case *ast.Index:
		e.Index = partialEvalExpression(e.Index)
		return e
	}

	return exp
}

func partialEvalBinaryOperation(b *ast.BinaryExpression) ast.Expression {
	left := partialEvalExpression(b.Left)
	right := partialEvalExpression(b.Right)

	switch left := left.(type) {
	case *ast.Number:
		if right, ok := right.(*ast.Number); ok {
			return partialNumber(left, right, b)
		}
		if right, ok := right.(*ast.Float); ok {
			left := &ast.Float{Token: left.Token, Value: float64(left.Value)}
			return partialFloat(left, right, b)
		}
	case *ast.Float:
		if right, ok := right.(*ast.Number); ok {
			right := &ast.Float{Token: right.Token, Value: float64(right.Value)}
			return partialFloat(left, right, b)
		}
		if right, ok := right.(*ast.Float); ok {
			return partialFloat(left, right, b)
		}
	}

	b.Left = left
	b.Right = right
	return b
}

func partialNumber(left, right *ast.Number, b *ast.BinaryExpression) ast.Expression {
	switch b.Operator {
	case "-":
		return &ast.Number{Value: left.Value - right.Value}
	case "*":
		return &ast.Number{Value: left.Value * right.Value}
	case "+":
		return &ast.Number{Value: left.Value + right.Value}
	case "/":
		reminder := left.Value % right.Value
		if reminder != 0 {
			return &ast.Float{Value: float64(left.Value) / float64(right.Value)}
		}
		return &ast.Number{Value: left.Value / right.Value}
	case "<<":
		return &ast.Number{Value: left.Value << right.Value}
	case "^":
		return &ast.Number{Value: left.Value ^ right.Value}
	case "<":
		return &ast.Boolean{Value: left.Value < right.Value}
	case ">":
		return &ast.Boolean{Value: left.Value > right.Value}
	case "==":
		return &ast.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &ast.Boolean{Value: left.Value != right.Value}
	}
	return b
}

func partialFloat(left, right *ast.Float, b *ast.BinaryExpression) ast.Expression {
	switch b.Operator {
	case "-":
		return &ast.Float{Value: left.Value - right.Value}
	case "*":
		return &ast.Float{Value: left.Value * right.Value}
	case "+":
		return &ast.Float{Value: left.Value + right.Value}
	case "/":
		return &ast.Float{Value: left.Value / right.Value}
	case "<":
		return &ast.Boolean{Value: left.Value < right.Value}
	case ">":
		return &ast.Boolean{Value: left.Value > right.Value}
	case "==":
		return &ast.Boolean{Value: left.Value == right.Value}
	case "!=":
		return &ast.Boolean{Value: left.Value != right.Value}
	}
	return b
}

func partialEvalUnaryOperation(e *ast.UnaryExpression) ast.Expression {
	expr := partialEvalExpression(e.Expression)

	switch e.Operator {
	case "!":
		return buildBang(expr, e)
	case "-":
		return buildMinus(expr, e)
	}
	return e
}

func buildBang(expr ast.Expression, e *ast.UnaryExpression) ast.Expression {
	switch expr := expr.(type) {
	case *ast.Boolean:
		return &ast.Boolean{Value: !expr.Value}
	case *ast.Number:
		if expr.Value == 0 {
			return &ast.Boolean{Value: true}
		}
		return &ast.Boolean{Value: false}
	case *ast.Float:
		if expr.Value == 0 {
			return &ast.Boolean{Value: true}
		}
		return &ast.Boolean{Value: false}
	case *ast.Null:
		return &ast.Boolean{Value: true}
	case *ast.Dictionary:
		return &ast.Boolean{Value: false}
	case *ast.Array:
		return &ast.Boolean{Value: false}
	case *ast.String:
		if len(expr.Value) == 0 {
			return &ast.Boolean{Value: true}
		}
		return &ast.Boolean{Value: true}
	}
	return e
}

func buildMinus(expr ast.Expression, e *ast.UnaryExpression) ast.Expression {
	switch expr := expr.(type) {
	case *ast.Number:
		return &ast.Number{Value: -expr.Value}
	case *ast.Float:
		return &ast.Float{Value: -expr.Value}
	}
	return e
}
