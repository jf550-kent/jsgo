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

	// Return represent the return node
	// return <expression>;
	ReturnStatement struct {
		Token            token.Token
		ReturnExpression Expression
	}

	// BlockStatement represent statements contian in a block
	BlockStatement struct {
		Token      token.Token
		Statements []Statement
	}

	// ExpressionStatement represent expression in a statement
	ExpressionStatement struct {
		Token      token.Token
		Expression Expression
	}

	// AssignmentStatement represent an assignment to a variable such as a = 10;
	AssignmentStatement struct {
		Token      token.Token
		Identifier *Identifier
		Expression Expression
	}

	// ForStatement represent the for loop in JSGO
	ForStatement struct {
		Token     token.Token
		Init      Statement
		Condition Expression
		Post      Statement
		Body      *BlockStatement
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

func (r *ReturnStatement) statementNode()   {}
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

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) Start() token.Pos     { return bs.Token.Start }
func (bs *BlockStatement) End() token.Pos {
	end := bs.Token.End

	lastStmt := bs.Statements[len(bs.Statements)-1]
	if lastStmt != nil {
		end = lastStmt.End()
	}

	return end
}
func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

func (e *ExpressionStatement) statementNode()       {}
func (e *ExpressionStatement) TokenLiteral() string { return e.Token.Literal }
func (e *ExpressionStatement) Start() token.Pos     { return e.Token.Start }
func (e *ExpressionStatement) End() token.Pos {
	end := e.Token.End

	if e.Expression != nil {
		end = e.Expression.End()
	}

	return end
}
func (e *ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return e.Token.String()
}

func (bs *AssignmentStatement) statementNode()       {}
func (bs *AssignmentStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *AssignmentStatement) Start() token.Pos     { return bs.Token.Start }
func (bs *AssignmentStatement) End() token.Pos {
	if bs.Expression != nil {
		return bs.Expression.End()
	}
	return bs.Token.End
}
func (bs *AssignmentStatement) String() string {
	var out strings.Builder

	out.WriteString(bs.Token.Literal)
	out.WriteString(" = ")
	if bs.Expression != nil {
		out.WriteString(bs.Expression.String())
	}
	return out.String()
}

// PLEASE change
func (n *ForStatement) statementNode()   {}
func (n *ForStatement) Start() token.Pos { return n.Token.Start }
func (n *ForStatement) End() token.Pos   { return n.Token.End }
func (n *ForStatement) String() string   { return n.Token.Literal }

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

	Boolean struct {
		Token token.Token
		Value bool
	}

	IFExpression struct {
		Token     token.Token
		Condition Expression
		Body      *BlockStatement
		Else      *BlockStatement
	}

	BinaryExpression struct {
		Token    token.Token
		Left     Expression
		Operator string
		Right    Expression
	}

	UnaryExpression struct {
		Token      token.Token
		Operator   string
		Expression Expression
	}

	FunctionDeclaration struct {
		Token      token.Token
		Parameters []*Identifier
		Body       *BlockStatement
	}

	CallExpression struct {
		Token     token.Token
		Function  Expression
		Arguments []Expression
	}

	String struct {
		Token token.Token
		Value string
	}

	Array struct {
		Token token.Token
		Body  []Expression
	}

	Index struct {
		Token      token.Token
		Identifier Expression
		Index      Expression
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

func (f *Boolean) expressionNode()  {}
func (f *Boolean) Start() token.Pos { return f.Token.Start }
func (f *Boolean) End() token.Pos   { return f.Token.End }
func (f *Boolean) String() string   { return f.Token.Literal }

func (i *IFExpression) expressionNode()  {}
func (i *IFExpression) Start() token.Pos { return i.Token.Start }
func (i *IFExpression) End() token.Pos {
	end := i.Token.End

	if i.Condition != nil {
		end = i.Condition.End()
	}

	if i.Else != nil {
		return i.Else.End()
	}

	if i.Body != nil {
		return i.Body.End()
	}

	return end
}
func (i *IFExpression) String() string {
	var st strings.Builder

	st.WriteString(i.Token.Literal)
	st.WriteString(" (")
	if i.Condition != nil {
		st.WriteString(i.Condition.String())
	}
	st.WriteString(") {")

	if i.Body != nil {
		st.WriteString(i.Body.String())
	}

	st.WriteString(" }")
	if i.Else != nil {
		st.WriteString(" else {")
		st.WriteString(i.Else.String())
		st.WriteString(" }")
	}
	st.WriteString(";")
	return st.String()
}

func (b *BinaryExpression) expressionNode() {}
func (b *BinaryExpression) Start() token.Pos {
	if b.Left != nil {
		return b.Left.Start()
	}
	return b.Token.Start
}
func (b *BinaryExpression) End() token.Pos {
	if b.Right != nil {
		return b.Right.End()
	}
	return b.Token.End
}
func (b *BinaryExpression) String() string {
	var s strings.Builder
	s.WriteString("(")
	if b.Left != nil {
		s.WriteString(b.Left.String())
	}

	s.WriteString(" " + b.Operator + " (")
	if b.Right != nil {
		s.WriteString(b.Right.String())
	}
	s.WriteString(")")
	return s.String()
}

func (u *UnaryExpression) expressionNode()  {}
func (u *UnaryExpression) Start() token.Pos { return u.Token.Start }
func (u *UnaryExpression) End() token.Pos {
	if u.Expression != nil {
		return u.Expression.End()
	}
	return u.Token.End
}
func (u *UnaryExpression) String() string {
	var s strings.Builder
	s.WriteString("(")
	s.WriteString(u.Operator)
	if u.Expression != nil {
		s.WriteString(u.Expression.String())
	}
	s.WriteString(")")
	return s.String()
}

func (f *FunctionDeclaration) expressionNode()  {}
func (f *FunctionDeclaration) Start() token.Pos { return f.Token.Start }
func (f *FunctionDeclaration) End() token.Pos {
	if f.Body != nil {
		return f.Body.End()
	}

	if len(f.Parameters) != 0 {
		return f.Parameters[len(f.Parameters)-1].End()
	}
	return f.Token.End
}
func (f *FunctionDeclaration) String() string {
	var s strings.Builder

	s.WriteString("function (")

	if f.Parameters != nil {
		for _, p := range f.Parameters {
			s.WriteString(p.String())
			s.WriteString(", ")
		}
	}

	s.WriteString(") {")
	if f.Body != nil {
		s.WriteString(f.Body.String())
	}
	s.WriteString("};")
	return s.String()
}

func (c *CallExpression) expressionNode() {}
func (c *CallExpression) Start() token.Pos {
	if c.Function != nil {
		return c.Function.Start()
	}

	return c.Token.Start
}
func (c *CallExpression) End() token.Pos {
	if len(c.Arguments) != 0 {
		return c.Arguments[len(c.Arguments)-1].End()
	}

	return c.Token.End
}
func (c *CallExpression) String() string {
	var out strings.Builder

	args := []string{}
	for _, a := range c.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(c.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

func (n *String) expressionNode()  {}
func (n *String) Start() token.Pos { return n.Token.Start }
func (n *String) End() token.Pos   { return n.Token.End }
func (n *String) String() string   { return n.Token.Literal }

func (al *Array) expressionNode() {}
func (n *Array) Start() token.Pos { return n.Token.Start }
func (n *Array) End() token.Pos   { return n.Token.End }
func (al *Array) String() string {
	var out strings.Builder

	elements := []string{}
	for _, el := range al.Body {
		elements = append(elements, el.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

func (ie *Index) expressionNode() {}
func (n *Index) Start() token.Pos { return n.Token.Start }
func (n *Index) End() token.Pos   { return n.Token.End }
func (ie *Index) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Identifier.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")

	return out.String()
}
