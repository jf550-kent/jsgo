package evaluator

import (
	"fmt"

	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(main *ast.Main) object.Object {
	obj := eval(main, object.NewEnvironment())
	err , ok := obj.(*object.Error)
	if ok {
		panic(err.Error())
	}
	return obj
}

func eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Main:
		return evalMain(node.Statements, env)
	case *ast.VarStatement:
		val := eval(node.Expression, env)
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
		return &object.Number{Value: node.Value}
	case *ast.Float:
		return &object.Float{Value: node.Value}
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Boolean:
		if node.Value {
			return TRUE
		}
		return FALSE
	case *ast.ExpressionStatement:
		return eval(node.Expression, env)
	case *ast.UnaryExpression:
		right := eval(node.Expression, env)
		if isError(right) {
			return right
		}
		return evalUnaryExpression(node.Operator, right)
	case *ast.BinaryExpression:
		left := eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalBinaryExpression(left, right, node.Operator)
	
	case *ast.FunctionDeclaration:
		params := node.Parameters
		body := node.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IFExpression:
		return evalIfExpression(node, env)
	}
	return nil
}

func evalMain(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = eval(stmt, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalExpressions(epxs []ast.Expression, env *object.Environment) []object.Object {
	result := []object.Object{}

	for _, e := range epxs {
		evaluated := eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function) 
	if !ok {
		return newError("not a function: %s", function.Type())
	}

	extendedEnv := extendFunctionEnv(function, args)
	evaluated := eval(function.Body, extendedEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Literal, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalBlockStatements(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJECT || rt == object.ERROR_OBJECT {
				return result
			}
		}
	}
	return result
}

func evalIfExpression(ie *ast.IFExpression, env *object.Environment) object.Object {
	condition := eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}
	if isTruthy(condition) {
		return eval(ie.Body, env)
	} else if ie.Else != nil {
		return eval(ie.Else, env)
	} else {
		return NULL
	}
}

func evalBinaryExpression(left, right object.Object, op string) object.Object {
	switch {
	case left.Type() == object.NUMBER_OBJECT && right.Type() == object.NUMBER_OBJECT:
		return evalNumberExpression(left, right, op)
	case left.Type() == object.FLOAT_OBJECT && right.Type() == object.FLOAT_OBJECT:
		return evalFloatExpression(left, right, op)
	case op == "!=":
		return nativeBoolean(left != right)
	case op == "==":
		return nativeBoolean(left == right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
	}

	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalUnaryExpression(op string, exp object.Object) object.Object {
	switch op {
	case "!":
		return evalBangOperatorExpression(exp)
	case "-":
		return evalNegativeOperatorExpression(exp)
	}
	return newError(fmt.Sprintf("unkonw operator: %s%s", op, exp.Type()))
}

func evalNumberExpression(left, right object.Object, op string) object.Object {
	leftValue, ok := left.(*object.Number)
	if !ok {
		return NULL
	}
	rightValue, ok := right.(*object.Number)
	if !ok {
		return NULL
	}
	switch op {
	case "-":
		return &object.Number{Value: leftValue.Value - rightValue.Value}
	case "*":
		return &object.Number{Value: leftValue.Value * rightValue.Value}
	case "+":
		return &object.Number{Value: leftValue.Value + rightValue.Value}
	case "/":
		return &object.Number{Value: leftValue.Value / rightValue.Value}
	case "<":
		return nativeBoolean(leftValue.Value < rightValue.Value)
	case ">":
		return nativeBoolean(leftValue.Value > rightValue.Value)
	case "==":
		return nativeBoolean(leftValue.Value == rightValue.Value)
	case "!=":
		return nativeBoolean(leftValue.Value != rightValue.Value)
	}
	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalFloatExpression(left, right object.Object, op string) object.Object {
	leftValue, ok := left.(*object.Float)
	if !ok {
		return NULL
	}
	rightValue, ok := right.(*object.Float)
	if !ok {
		return NULL
	}
	switch op {
	case "-":
		return &object.Float{Value: leftValue.Value - rightValue.Value}
	case "*":
		return &object.Float{Value: leftValue.Value * rightValue.Value}
	case "+":
		return &object.Float{Value: leftValue.Value + rightValue.Value}
	case "/":
		return &object.Float{Value: leftValue.Value / rightValue.Value}
	case "<":
		return nativeBoolean(leftValue.Value < rightValue.Value)
	case ">":
		return nativeBoolean(leftValue.Value > rightValue.Value)
	case "==":
		return nativeBoolean(leftValue.Value == rightValue.Value)
	case "!=":
		return nativeBoolean(leftValue.Value != rightValue.Value)
	}
	return newError("unknown operator: %s %s %s", left.Type(), op, right.Type())
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	}
	return FALSE
}

func evalNegativeOperatorExpression(exp object.Object) object.Object {
	switch r := exp.(type) {
	case *object.Number:
		return &object.Number{Value: -r.Value}
	case *object.Float:
		return &object.Float{Value: -r.Value}
	}
	return newError(fmt.Sprintf("unable to minus value: -%s for type: %s", exp.String(), exp.Type()))
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
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

func nativeBoolean(b bool) object.Object {
	if b {
		return TRUE
	}
	return FALSE
}