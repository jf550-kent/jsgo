package ast

import (
	"strings"

	"github.com/jf550-kent/jsgo/token"
)

// Node is the smallest element in the abstract syntax tree of JSGO.
type Node interface {
	Start() token.Pos // The first position of the [Node]'s first [token.Token]
	End() token.Pos   // The last position of the [Node]'s last character of the [token.Token]
	String() string
}

// Statement is the node in the ast that represent statement in JSGO.
type Statement interface {
	Node
	statementNode()
}

// Expression is the node in the ast that represent expression in JSGO, expression in JSGO produce value.
type Expression interface {
	Node
	expressionNode()
}

// Main is the entry point of the program
type Main struct {
	Name       string
	Statements []Statement
}

// Start return the first [token.Pos] of the first statement
func (m *Main) Start() token.Pos {
	if len(m.Statements) == 0 {
		return token.Pos{Line: 0, Col: 0}
	}
	return m.Statements[0].Start()
}

// End return the last [token.Pos] of the last statement
func (m *Main) End() token.Pos {
	if len(m.Statements) == 0 {
		return token.Pos{Line: 0, Col: 0}
	}
	last := len(m.Statements) - 1
	return m.Statements[last].End()
}

func (m *Main) String() string {
	if len(m.Statements) == 0 {
		return "invalid main file"
	}
	return m.Name
}

// Identifier is the node that represent an identifier.
type Identifier struct {
	Token   token.Token
	Literal string
}

func (i *Identifier) expressionNode()  {}
func (i *Identifier) Start() token.Pos { return i.Token.Start }
func (i *Identifier) End() token.Pos   { return i.Token.End }
func (i *Identifier) String() string   { return i.Literal }

// Statement in the language does not produce value
type (
	// VarStatement represent the var node
	// var <identifier> = <expression>;
	VarStatement struct {
		Token      token.Token
		Variable   *Identifier
		Expression Expression
	}

	// Return represetn the return node
	// return <expression>;
	ReturnStatement struct {
		Token token.Token
		ReturnExpression Expression
	}
)

func (v *VarStatement) statementNode()   {}
func (v *VarStatement) Start() token.Pos { return v.Token.Start }
func (v *VarStatement) End() token.Pos {
	if v.Expression != nil {
		return v.Expression.End()
	}
	return v.Token.End
}
func (v *VarStatement) String() string {
	var s strings.Builder
	s.WriteString("var ")
	if v.Variable != nil {
		s.WriteString(v.Variable.String())
	}
	s.WriteString(" = ")
	if v.Expression != nil {
		s.WriteString(v.Expression.String())
	}
	s.WriteString(";")
	return s.String()
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) Start() token.Pos { return r.Token.Start }
func (r *ReturnStatement) End() token.Pos { 
	if r.ReturnExpression != nil {
		return r.ReturnExpression.End()
	}
	return r.Token.End
}
func (r *ReturnStatement) String() string {
	var s strings.Builder
	s.WriteString(r.Token.String())
	s.WriteString(" ")
	if r.ReturnExpression != nil {
		s.WriteString(r.ReturnExpression.String())
	}
	s.WriteString(";")
	return s.String()
}

// expression
type (
	Number struct {
		Token token.Token
		Value int64
	}

	Float struct {
		Token token.Token
		Value float64
	}
)

func (n *Number) expressionNode()  {}
func (n *Number) Start() token.Pos { return n.Token.Start }
func (n *Number) End() token.Pos   { return n.Token.End }
func (n *Number) String() string   { return n.Token.Literal }

func (f *Float) expressionNode()  {}
func (f *Float) Start() token.Pos { return f.Token.Start }
func (f *Float) End() token.Pos   { return f.Token.End }
func (f *Float) String() string   { return f.Token.Literal }
