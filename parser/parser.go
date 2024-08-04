package parser

import (
	"fmt"
	"strconv"

	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/lexer"
	"github.com/jf550-kent/jsgo/token"
)

// every parse<> needs to make the current parser to point at ; then is up to the line:37 to skip ;

type (
	unaryExpressionFunc  func() ast.Expression
	binaryExpressionFunc func(ast.Expression) ast.Expression
)

const (
	SYNTAX_ERROR   = "SyntaxError"
	TYPE_ERROR     = "TypeError"
	INTERNAL_ERROR = "InternalError"
	ILLEGAL_TOKEN  = "IllegalToken"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
	CALL        // function()
)

var precedences = map[token.TokenType]int{
	token.EQUAL:     EQUALS,
	token.NOT_EQUAL: EQUALS,
	token.GTR:       LESSGREATER,
	token.LSS:       LESSGREATER,
	token.ADD:       SUM,
	token.MINUS:     SUM,
	token.DIVIDE:    PRODUCT,
	token.MUL:       PRODUCT,
	token.LPAREN:    CALL,
}

func Parse(filename string, src []byte) *ast.Main {
	l := lexer.New(src)

	p := new(filename, l)

	// fill up the first 2 token in parser
	p.next()
	p.next()

	main := &ast.Main{Name: filename, Statements: []ast.Statement{}}

	for p.currentToken.TokenType != token.EOF {
		stmt := p.parse()
		if stmt != nil {
			main.Statements = append(main.Statements, stmt)
		}
		p.next()
	}

	return main
}

type parser struct {
	l    *lexer.Lexer
	name string

	currentToken token.Token
	nextToken    token.Token

	unaryExpressionFuncs map[token.TokenType]unaryExpressionFunc
	binaryExpressionFunc map[token.TokenType]binaryExpressionFunc
}

func new(filename string, l *lexer.Lexer) *parser {
	p := &parser{l: l, name: filename}

	p.binaryExpressionFunc = map[token.TokenType]binaryExpressionFunc{}

	p.unaryExpressionFuncs = map[token.TokenType]unaryExpressionFunc{
		token.IDENT:    p.parseIdent,
		token.NUMBER:   p.parseNumber,
		token.FLOAT:    p.parseFloat,
		token.TRUE:     p.parseBoolean,
		token.FALSE:    p.parseBoolean,
		token.IF:       p.parseIFExpression,
		token.MINUS:    p.parseUnaryExpression,
		token.BANG:     p.parseUnaryExpression,
		token.FUNCTION: p.parseFunctionDeclaration,
		token.LPAREN:   p.parseGroupedExpression,
	}

	p.binaryExpressionFunc = map[token.TokenType]binaryExpressionFunc{
		token.ADD:       p.parseBinaryExpression,
		token.MINUS:     p.parseBinaryExpression,
		token.MUL:       p.parseBinaryExpression,
		token.DIVIDE:    p.parseBinaryExpression,
		token.LSS:       p.parseBinaryExpression,
		token.GTR:       p.parseBinaryExpression,
		token.NOT_EQUAL: p.parseBinaryExpression,
		token.EQUAL:     p.parseBinaryExpression,
		token.LPAREN:    p.parseCallExpression,
	}

	return p
}

func (p *parser) parse() ast.Statement {
	switch p.currentToken.TokenType {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT:
		if !p.peekExpect(token.ASSIGN) {
			break
		}
		return p.parseAssignmentStatement()
	case token.FOR:
		return p.parseForStatement()
	}
	return p.parseExpressionStatement()
}

func (p *parser) parseVarStatement() ast.Statement {
	varStmt := &ast.VarStatement{Token: p.currentToken}

	if !p.peekExpect(token.IDENT) {
		errMsg := varStmt.String() + "expect variable when declaring var"
		p.panicError(errMsg, SYNTAX_ERROR, varStmt.Start())
		return nil
	}
	p.next()

	ident, ok := p.parseIdent().(*ast.Identifier)
	if !ok {
		p.panicError("failed to parse identifier", INTERNAL_ERROR, varStmt.End())
		return nil
	}
	varStmt.Variable = ident
	p.next()

	if !p.expect(token.ASSIGN) {
		errMsg := varStmt.String() + " :expect = after identifier when declaring var"
		p.panicError(errMsg, SYNTAX_ERROR, varStmt.Start())
		return nil
	}
	p.next()

	varStmt.Expression = p.parseExpression(1)

	if !p.peekExpect(token.SEMICOLON) {
		errMsg := varStmt.String() + " :expect ; after expression when declaring var"
		p.panicError(errMsg, SYNTAX_ERROR, varStmt.Start())
		return nil
	}
	p.next()

	return varStmt
}

func (p *parser) parseReturnStatement() ast.Statement {
	st := &ast.ReturnStatement{Token: p.currentToken}
	p.next()

	st.ReturnExpression = p.parseExpression(1)

	if !p.peekExpect(token.SEMICOLON) {
		errMsg := st.String() + " :expect ; after expression in return statement"
		p.panicError(errMsg, SYNTAX_ERROR, st.End())
		return nil
	}
	p.next()

	return st
}

func (p *parser) parseExpressionStatement() ast.Statement {
	stmt := &ast.ExpressionStatement{Token: p.currentToken}

	stmt.Expression = p.parseExpression(1)
	if stmt.Expression == nil {
		return nil
	}

	if !p.peekExpect(token.SEMICOLON) {
		err := fmt.Sprintf("%s: missing ;", stmt.String())
		p.panicError(err, SYNTAX_ERROR, stmt.End())
	}

	p.next()
	return stmt
}

func (p *parser) parseAssignmentStatement() ast.Statement {
	p.check(token.IDENT)
	assign := &ast.AssignmentStatement{Token: p.currentToken}
	ident := p.parseIdent()
	if i, ok := ident.(*ast.Identifier); ok {
		assign.Identifier = i
	}

	if assign.Identifier == nil {
		p.panicError("compiler unable to parse identifier", INTERNAL_ERROR, p.currentToken.Start)
	}
	p.next()
	if !p.expect(token.ASSIGN) {
		p.panicError("for assignment expected to have = after identifier", SYNTAX_ERROR, assign.End())
	}
	p.next()
	assign.Expression = p.parseExpression(1)

	if p.peekExpect(token.SEMICOLON) {
		p.next()
	}
	return assign
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	unaryFunc, ok := p.unaryExpressionFuncs[p.currentToken.TokenType]
	if !ok {
		errMsg := fmt.Sprintf("unary expression not found for %s", p.currentToken)
		p.panicError(errMsg, SYNTAX_ERROR, p.currentToken.Start)
		return nil
	}

	result := unaryFunc()

	for !p.peekExpect(token.SEMICOLON) && precedence < p.peekPred() {
		binaryFunc, ok := p.binaryExpressionFunc[p.nextToken.TokenType]
		if !ok {
			return result
		}
		p.next()
		result = binaryFunc(result)
	}

	return result
}

func (p *parser) parseBinaryExpression(left ast.Expression) ast.Expression {
	expr := &ast.BinaryExpression{
		Token:    p.currentToken,
		Left:     left,
		Operator: p.currentToken.Literal,
	}

	pred := p.pred()
	p.next()
	expr.Right = p.parseExpression(pred)
	return expr
}

func (p *parser) parseCallExpression(left ast.Expression) ast.Expression {
	c := &ast.CallExpression{Token: p.currentToken, Function: left}

	c.Arguments = []ast.Expression{}

	if p.peekExpect(token.RPAREN) {
		p.next()
		return c
	}

	p.next()

	c.Arguments = append(c.Arguments, p.parseExpression(LOWEST))

	for p.peekExpect(token.COMMA) {
		p.next()
		p.next()
		c.Arguments = append(c.Arguments, p.parseExpression(LOWEST))
	}

	if !p.peekExpect(token.RPAREN) {
		err := c.String() + " : missing ) "
		p.panicError(err, SYNTAX_ERROR, c.End())
	}
	p.next()
	return c
}

func (p *parser) parseUnaryExpression() ast.Expression {
	ury := &ast.UnaryExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
	}
	p.next()
	ury.Expression = p.parseExpression(6)

	return ury
}

func (p *parser) parseIFExpression() ast.Expression {
	exp := &ast.IFExpression{Token: p.currentToken}

	if !p.peekExpect(token.LPAREN) {
		errMsg := exp.String() + " : missing ( "
		p.panicError(errMsg, SYNTAX_ERROR, exp.End())
	}
	p.next()
	p.next()

	exp.Condition = p.parseExpression(1)

	if !p.peekExpect(token.RPAREN) {
		errMsg := exp.String() + " missing )"
		p.panicError(errMsg, SYNTAX_ERROR, exp.End())
	}
	p.next()

	if !p.peekExpect(token.LBRACE) {
		errMsg := exp.String() + " missing {"
		p.panicError(errMsg, SYNTAX_ERROR, exp.End())
	}
	p.next()

	exp.Body = p.parseBlockStatement()

	if p.peekExpect(token.ELSE) {
		p.next()

		if !p.peekExpect(token.LBRACE) {
			errMsg := exp.String() + " missing {"
			p.panicError(errMsg, SYNTAX_ERROR, exp.End())
		}
		p.next()
		exp.Else = p.parseBlockStatement()
	}

	if !p.peekExpect(token.SEMICOLON) {
		err := exp.String() + " : expected ; after if expression"
		p.panicError(err, SYNTAX_ERROR, exp.End())
	}

	return exp
}

func (p *parser) parseFunctionDeclaration() ast.Expression {
	f := &ast.FunctionDeclaration{Token: p.currentToken}
	p.next()

	// function () function (a, b, t) {}
	if !p.expect(token.LPAREN) {
		err := f.String() + " : expected ( for function declaration"
		p.panicError(err, SYNTAX_ERROR, f.End())
	}

	f.Parameters = p.parseFunctionParameters()
	p.next()
	if !p.peekExpect(token.LBRACE) {
		err := f.String() + " : missing { for function declaration"
		p.panicError(err, SYNTAX_ERROR, f.End())
	}
	p.next()

	f.Body = p.parseBlockStatement()

	if !p.peekExpect(token.SEMICOLON) {
		err := f.String() + " : missing ; for function declaration"
		p.panicError(err, SYNTAX_ERROR, f.End())
	}
	return f
}

func (p *parser) parseFunctionParameters() []*ast.Identifier {
	r := []*ast.Identifier{}

	if p.peekExpect(token.RPAREN) {
		return r
	}
	p.next()
	if !p.expect(token.IDENT) {
		err := "only identifier allowed in funciton parameters"
		p.panicError(err, SYNTAX_ERROR, p.currentToken.End)
	}

	for p.expect(token.IDENT) {
		id := &ast.Identifier{Token: p.currentToken, Literal: p.currentToken.Literal}
		r = append(r, id)

		if p.peekExpect(token.RPAREN) {
			break
		}
		if !p.peekExpect(token.COMMA) {
			err := "missing , in function parameters"
			p.panicError(err, SYNTAX_ERROR, id.End())
		}
		p.next()
		p.next()
	}
	return r
}

// parseBlockStatement always start the when the current token in the parser is { and ends at }
func (p *parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}
	p.next()

	for !p.expect(token.RBRACE) && !p.expect(token.EOF) {
		stmt := p.parse()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.next()
	}
	return block
}

func (p *parser) parseGroupedExpression() ast.Expression {
	p.next()
	exp := p.parseExpression(LOWEST)
	if !p.peekExpect(token.RPAREN) {
		return nil
	}
	p.next()
	return exp
}

func (p *parser) parseForStatement() ast.Statement {
	p.check(token.FOR)
	forStmt := &ast.ForStatement{Token: p.currentToken}

	p.next()
	if !p.expect(token.LPAREN) {
		p.panicError(fmt.Sprintf("%s : expecting (", p.currentToken), SYNTAX_ERROR, p.currentToken.End)
	}
	p.next()

	forStmt.Init = p.parseVarStatement()
	p.next()

	forStmt.Condition = p.parseExpression(1)
	p.next()
	if !p.expect(token.SEMICOLON) {
		p.panicError(fmt.Sprintf("%s : expecting ; after initialisation", forStmt), SYNTAX_ERROR, p.currentToken.End)
	}
	p.next()
	forStmt.Post = p.parseAssignmentStatement()
	p.next()
	if !p.expect(token.RPAREN) {
		p.panicError(fmt.Sprintf("%s : expecting ; after initialisation", forStmt), SYNTAX_ERROR, p.currentToken.End)
	}
	p.next()
	forStmt.Body = p.parseBlockStatement()
	return forStmt
}

func (p *parser) parseIdent() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Literal: p.currentToken.Literal}
}

func (p *parser) parseNumber() ast.Expression {
	v, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		p.panicError("unable to convert number", INTERNAL_ERROR, p.currentToken.Start)
		return nil
	}

	return &ast.Number{Token: p.currentToken, Value: v}
}

func (p *parser) parseFloat() ast.Expression {
	f, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.panicError("unable to convert float", INTERNAL_ERROR, p.currentToken.Start)
		return nil
	}
	return &ast.Float{Token: p.currentToken, Value: f}
}

func (p *parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: p.currentToken, Value: p.expect(token.TRUE)}
}

func (p *parser) peekPred() int {
	pr, ok := precedences[p.nextToken.TokenType]
	if !ok {
		return LOWEST
	}
	return pr
}

func (p *parser) pred() int {
	pr, ok := precedences[p.currentToken.TokenType]
	if !ok {
		return LOWEST
	}
	return pr
}

func (p *parser) expect(t token.TokenType) bool {
	return p.currentToken.TokenType == t
}

func (p *parser) peekExpect(t token.TokenType) bool {
	return p.nextToken.TokenType == t
}

func (p *parser) next() {
	p.currentToken = p.nextToken
	ntTok, err := p.l.Lex()
	if err != nil {
		p.panicError(err.Error(), ILLEGAL_TOKEN, ntTok.Start)
	}
	p.nextToken = ntTok
}

func (p *parser) panicError(msg, errorType string, pos token.Pos) {
	panic(fmt.Errorf("%s: %s %s:%d:%d", errorType, msg, p.name, pos.Line, pos.Col))
}

// check should be used to check if tok is the parser current token.
// if not it will panic. In the parser, it is up to the developer to advance the token
// therefore check() acts as an check for before creating the ast to verify the expected token
// ONLY use for developer errors, such as moving token wrongly.
func (p *parser) check(tok token.TokenType) {
	if p.currentToken.TokenType != tok {
		p.panicError(INTERNAL_ERROR, "expecting "+tok.String(), p.currentToken.Start)
	}
}
