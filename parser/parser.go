package parser

import (
	"fmt"
	"strconv"

	"github.com/jf550-kent/jsgo/ast"
	"github.com/jf550-kent/jsgo/lexer"
	"github.com/jf550-kent/jsgo/token"
)

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
		token.IDENT:  p.parseIdent,
		token.NUMBER: p.parseNumber,
		token.FLOAT:  p.parseFloat,
		token.TRUE: p.parseBoolean,
		token.FALSE: p.parseBoolean,
		token.IF: p.parseIFExpression, 
	}

	return p
}

func (p *parser) parse() ast.Statement {
	switch p.currentToken.TokenType {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
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
	return nil
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

func (p *parser) parseIFExpression() ast.Expression {
	exp := &ast.IFExpression{Token: p.currentToken}

	if !p.peekExpect(token.LPAREN) {
		errMsg := exp.String() + " : missing ( " 
		p.panicError(errMsg, SYNTAX_ERROR, exp.End())
	}
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

	// if p.expect(token.ELSE) {

	// }

	if !p.peekExpect(token.SEMICOLON) {
		err := exp.String() + " : expected ; after if expression"
		p.panicError(err, SYNTAX_ERROR, exp.End())
	}
	p.next()

	return exp
}

func (p *parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.currentToken}
	block.Statements = []ast.Statement{}

	for !p.expect(token.SEMICOLON) && !p.expect(token.EOF) {
		stmt := p.parse()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.next()
	}
	return block
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
	return p.nextToken.Precedence()
}

func (p *parser) pred() int {
	return p.currentToken.Precedence()
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
