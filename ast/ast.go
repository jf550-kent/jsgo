package ast

import (
	"strings"
	"github.com/jun-hf/jsgo/token"
)

type Node interface {
	Start() token.Pos // The first position of the [Node]'s first [token.Token]
	End() token.Pos   // The last position of the [Node]'s last character of the [token.Token]
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Main is the entry point of the program
type Main struct {
	Name string
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

type Identifier struct {
	Token token.Token
	Literal string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) Start() token.Pos { return i.Token.Start }
func (i *Identifier) End() token.Pos { return i.Token.End }
func (i *Identifier) String() string { return i.Literal }

// Statement in the language does not produce value
type (
	// VarStatement represent the var node
	// var <identifier> = <expression>;
	VarStatement struct {
		Token token.Token
		Variable *Identifier
		Expression Expression
	}
)

func (v *VarStatement) statementNode() {}
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

func (n *Number) expressionNode() {}
func (n *Number) Start() token.Pos { return n.Token.Start } 
func (n *Number) End() token.Pos   { return n.Token.End } 
func (n *Number) String() string { return n.Token.Literal }

func (f *Float) expressionNode() {}
func (f *Float) Start() token.Pos { return f.Token.Start } 
func (f *Float) End() token.Pos   { return f.Token.End } 
func (f *Float) String() string { return f.Token.Literal }