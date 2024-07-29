package evaluator

import (
	"fmt"

	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/object"
)

// func Eval(m ast.Main) (object.Object, error) {
// 	var result object.Object

// 	for _, stmt := range m.Statements {
// 		result = eval(stmt)

// 		switch r := result.(type) {
// 		case *object.ReturnValue:
// 			return r.Value, nil
// 		case *object.Error:
// 			return nil, r
// 		}
// 	}
// }

// func eval (node ast.Node) object.Object {

// }

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Main:
		return Eval(node.Statements[0], env)
	case *ast.VarStatement:
		val := Eval(node.Expression, env)
		if isError(val) {
			return val
		}
		env.Set(node.Variable.Literal, val)

		v, ok := env.Get(node.Variable.Literal)
		if !ok {
			return newError("failed to set variable")
		}
		return v
	case *ast.Number:
		num := &object.Number{Value: node.Value}
		return num
	case *ast.Identifier:
		return evalIdentifier(node, env)
	}
	return nil
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Literal)
	if !ok {
		return newError("identifier not found: " + node.Literal)
	}
	return val
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJECT
	}
	return false
}