package object

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jf550-kent/jsgo/ast"
)

type ObjectType string

const (
	NUMBER_OBJECT       ObjectType = "NUMBER"
	FLOAT_OBJECT        ObjectType = "FLOAT"
	BOOLEAN_OBJECT      ObjectType = "BOOLEAN"
	NULL_OBJECT         ObjectType = "NULL"
	RETURN_VALUE_OBJECT ObjectType = "RETURN_VALUE"
	ERROR_OBJECT        ObjectType = "ERROR"
	FUNCTION_OBJECT     ObjectType = "FUNCTION"
	STRING_OBJECT       ObjectType = "STRING"
	ARRAY_OBJECT        ObjectType = "ARRAY"
)

// Object is used in the evaluator to represent value in when evaluating the AST of JSGO.
type Object interface {
	Type() ObjectType
	String() string
}

// Number is the object representing the number type in JSGO
type Number struct {
	Value int64
}

func (n *Number) String() string   { return fmt.Sprintf("%d", n.Value) }
func (n *Number) Type() ObjectType { return NUMBER_OBJECT }

type Float struct {
	Value float64
}

func (f *Float) String() string   { return strconv.FormatFloat(f.Value, 'f', -1, 64) }
func (f *Float) Type() ObjectType { return FLOAT_OBJECT }

// Boolean represent the boolean value in the language when evaluating the ast
type Boolean struct {
	Value bool
}

func (b *Boolean) String() string   { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return ObjectType(BOOLEAN_OBJECT) }

// Null represent the NULL value in the language, it means that there is no value
type Null struct{}

func (n *Null) String() string   { return "NULL" }
func (n *Null) Type() ObjectType { return NULL_OBJECT }

// Function represent the Function declaration.
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) String() string {
	var out strings.Builder

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) Type() ObjectType { return FUNCTION_OBJECT }

type String struct {
	Value string
}

func (s *String) String() string   { return s.Value }
func (s *String) Type() ObjectType { return STRING_OBJECT }

type Array struct {
	Body []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJECT }
func (a *Array) String() string {
	var out strings.Builder

	elements := []string{}
	for _, el := range a.Body {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// ReturnValue represent the value that is being returned
type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) String() string   { return rv.Value.String() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJECT }

// Error represent the error object in when evaluating the AST.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJECT }
func (e *Error) String() string   { return "error: " + e.Message }
func (e *Error) Error() string    { return e.Message }
