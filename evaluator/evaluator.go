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
	err, ok := obj.(*object.Error)
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
	case *ast.AssignmentStatement:
		val := eval(node.Expression, env)
		if isError(val) {
			return val
		}

		_, env, ok := env.GetIdentifier(node.Identifier.Literal)
		if !ok {
			return newError("failed to set variable")
		}
		env.Set(node.Identifier.Literal, val)
		return val
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
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		if index, ok := node.Function.(*ast.Index); ok {
			right := eval(index.Index, env)
			if isError(right) {
				return right
			}
			if method, ok := right.(*object.String); ok {
				return evalArrayMethodEpression(index, method, args, env)
			}
		}
		function := eval(node.Function, env)
		if isError(function) {
			return function
		}
		return callFunction(function, args)
	case *ast.ForStatement:
		return evalForStatement(node, env)
	case *ast.BlockStatement:
		return evalBlockStatements(node, env)
	case *ast.IFExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := eval(node.ReturnExpression, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.String:
		return &object.String{Value: node.Value}
	case *ast.Array:
		body := evalExpressions(node.Body, env)
		if len(body) == 1 && isError(body[0]) {
			return body[0]
		}
		return &object.Array{Body: body}
	case *ast.Index:
		ident := eval(node.Identifier, env)
		if isError(ident) {
			return ident
		}
		index := eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(ident, index)
	case *ast.Null:
		return NULL
	case *ast.Dictionary:
		return evalDictionary(node, env)
	case *ast.BracketDeclaration:
		return evalDictionaryDeclaration(node, env)
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

func callFunction(fn object.Object, args []object.Object) object.Object {

	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.BuiltIn:
		return fn.Function(args...)
	}
	return newError("not a function: %s", fn.Type())
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
	lType := left.Type()
	rType := right.Type()
	switch {
	case lType == object.NUMBER_OBJECT && rType == object.NUMBER_OBJECT:
		return evalNumberExpression(left, right, op)
	case lType == object.FLOAT_OBJECT && rType == object.FLOAT_OBJECT:
		return evalFloatExpression(left, right, op)
	case (lType == object.FLOAT_OBJECT && rType == object.NUMBER_OBJECT) || (lType == object.NUMBER_OBJECT && rType == object.FLOAT_OBJECT):
		l := convertFloat(left)
		r := convertFloat(right)
		return evalFloatExpression(l, r, op)
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

func evalForStatement(forStmt *ast.ForStatement, env *object.Environment) object.Object {
	if forStmt.Init != nil {
		eval(forStmt.Init, env)
	}

	for {
		// Evaluate the condition
		if forStmt.Condition == nil {
			panic("for loop does not contain condition : is the middle bit :)")
		}

		condition := eval(forStmt.Condition, env)
		if isError(condition) {
			return condition
		}
		if !isTruthy(condition) {
			break
		}
		body := eval(forStmt.Body, env)
		if isError(body) {
			return body
		}
		if _, ok := body.(*object.ReturnValue); ok {
			return body
		}
		if forStmt.Post != nil {
			post := eval(forStmt.Post, env)
			if isError(post) {
				return post
			}
		}
	}

	return NULL
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
		reminder := leftValue.Value % rightValue.Value
		if reminder != 0 {
			return &object.Float{Value: float64(leftValue.Value) / float64(rightValue.Value)}
		}
		return &object.Number{Value: leftValue.Value / rightValue.Value}
	case "<<":
		return &object.Number{Value: leftValue.Value << rightValue.Value}
	case "^":
		return &object.Number{Value: leftValue.Value ^ rightValue.Value}
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
	if val, ok := env.Get(node.Literal); ok {
		return val
	}

	if val, ok := builtin[node.Literal]; ok {
		return val
	}
	return newError("identifier not found: " + node.Literal)
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch left := left.(type) {
	case *object.Array:
		return evalArrayIndexExpression(left, index)
	case *object.Dictionary:
		return evalDictionaryExpression(left, index)
	}
	return newError("index operator not supported: %s", left.Type())
}

func evalArrayIndexExpression(arr *object.Array, index object.Object) object.Object {
	size := int64(len(arr.Body))

	switch right := index.(type) {
	case *object.Number:
		idx := right.Value
		if idx < 0 || idx >= size {
			return NULL
		}
		return arr.Body[idx]
	case *object.String:
		if right.Value == "length" {
			return &object.Number{Value: size}
		}
		return newError("array index unsupported for " + right.Value)
	}
	return newError("array index unsupported for type: " + index.String())
}

func evalArrayMethodEpression(index *ast.Index, method *object.String, args []object.Object, env *object.Environment) object.Object {
	ident := eval(index.Identifier, env)
	if isError(ident) {
		return ident
	}
	arr, ok := ident.(*object.Array)
	if !ok {
		return newError("unsupported array call")
	}
	switch method.Value {
	case "push":
		arr.Body = append(arr.Body, args...)
	}
	return arr
}

func evalDictionary(dic *ast.Dictionary, env *object.Environment) object.Object {

	dicry := &object.Dictionary{Value: make(map[object.Hash]object.KeyValue)}

	for key, val := range dic.Object {
		if _, err := assignDictionaryKey(dicry, key, val, env); err != nil {
			return err
		}
	}
	return dicry
}

func assignDictionaryKey(dic *object.Dictionary, key, value ast.Node, env *object.Environment) (*object.Dictionary, *object.Error) {
	k := eval(key, env)
	if isError(k) {
		return nil, k.(*object.Error) // dangerous man you should know better smh
	}

	h, ok := k.(object.Hasher)
	if !ok {
		return nil, newError("key unable to be hash" + k.String())
	}
	v := eval(value, env)
	if isError(v) {
		return nil, v.(*object.Error)
	}

	keyHash := h.Hash()
	dic.Value[keyHash] = object.KeyValue{Key: k, Value: v}

	return dic, nil
}

func evalDictionaryDeclaration(decl *ast.BracketDeclaration, env *object.Environment) object.Object {
	ident := eval(decl.Identifier, env)
	if isError(ident) {
		return ident
	}

	switch left := ident.(type) {
	case *object.Dictionary:
		if _, err := assignDictionaryKey(left, decl.Key, decl.Value, env); err != nil {
			return err
		}
		return left
	case *object.Array:
		in := eval(decl.Key, env)
		if isError(in) {
			newError("unable to evaluate array index")
		}
		num, ok := in.(*object.Number)
		if !ok {
			return newError("wrong type for array index")
		}
		if int(num.Value) >= len(left.Body) {
			// a more efficient assignment can be done here
			newArr := make([]object.Object, int(num.Value) + 1)
			copy(newArr, left.Body)
			left.Body = newArr
		}
		v := eval(decl.Value, env)
		if isError(v) {
			return newError("unable to eval array assigment")
		}
		left.Body[int(num.Value)] = v
		return left
	}

	msg := fmt.Sprintf("undefined identifier <%s>[] reference for %v:%v", ident.String(),  decl.Start().Line ,decl.Start().Col)
	return newError(msg)
}

func evalDictionaryExpression(dic *object.Dictionary, right object.Object) object.Object {
	key, ok := right.(object.Hasher)
	if !ok {
		return newError("cannot use " + right.String() + "as dictionary index")
	}

	hashKey := key.Hash()
	keyValue, ok := dic.Value[hashKey]
	if !ok {
		return NULL
	}
	return keyValue.Value
}

func convertFloat(node object.Object) *object.Float {
	switch node := node.(type) {
	case *object.Float:
		return node
	case *object.Number:
		return &object.Float{Value: float64(node.Value)}
	}
	panic("unable to convert to float with" + node.String())
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
