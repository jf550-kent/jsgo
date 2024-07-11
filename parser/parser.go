package parser

import (
	"fmt"
	"strconv"

	"github.com/jun-hf/jsgo/ast"
	"github.com/jun-hf/jsgo/lexer"
	"github.com/jun-hf/jsgo/token"
)

type (
	unaryExpressionFunc  func() ast.Expression
	binaryExpressionFunc func(ast.Expression) ast.Expression
)

const (
	SYNTAX_ERROR = "SyntaxError"
	TYPE_ERROR = "TypeError"
	INTERNAL_ERROR = "InternalError"
)

type parser struct {
	l    *lexer.Lexer
	name string

	currentToken token.Token
	nextToken    token.Token

	errorList []error

	unaryExpressionFuncs map[token.TokenType]unaryExpressionFunc
	binaryExpressionFunc map[token.TokenType]binaryExpressionFunc
}

func new(filename string, l *lexer.Lexer) *parser {
	p := &parser{l: l, errorList: []error{}, name: filename}

	p.binaryExpressionFunc = map[token.TokenType]binaryExpressionFunc{
	}

	p.unaryExpressionFuncs = map[token.TokenType]unaryExpressionFunc{
		token.IDENT: p.parseIdent,
		token.NUMBER: p.parseNumber,
		token.FLOAT: p.parseFloat,
	}

	return p
}

func Parse(filename string, src []byte) (*ast.Main, error) {
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

	return main, processError(p.errorList)
}

func (p *parser) parse() ast.Statement {
	switch p.currentToken.TokenType {
	case token.VAR:
		return p.parseVarStatement()
	}
	return p.parseExprStmt()
}

func (p *parser) parseVarStatement() ast.Statement {
	varStmt := &ast.VarStatement{Token: p.currentToken}

	if !p.peekExpect(token.IDENT) {
		errMsg := varStmt.String() + "expect variable when declaring var"
		p.addError(errMsg, SYNTAX_ERROR, varStmt.Start())
		return nil
	}

	ident, ok := p.parseIdent().(*ast.Identifier)
	if !ok {
		p.addError("", INTERNAL_ERROR, varStmt.End())
		return nil
	}
	varStmt.Variable = ident

	p.next()

	varStmt.Expression = p.parseExpression(1)

	if p.expect(token.SEMICOLON) {
		p.next()
	}
	
	return varStmt
}

func (p *parser) parseExprStmt() ast.Statement {
	return nil
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	unaryFunc, ok := p.unaryExpressionFuncs[p.currentToken.TokenType]
	if !ok {
		errMsg := fmt.Sprintf("unary expression not found for %s", p.currentToken)
		p.addError(errMsg, SYNTAX_ERROR, p.currentToken.Start)
		return nil
	}
	result := unaryFunc()

	for !p.expect(token.SEMICOLON) && precedence < p.peekPred() {
		binaryFunc, ok := p.binaryExpressionFunc[p.nextToken.TokenType]
		if !ok {
			return result
		}
		p.next()
		result = binaryFunc(result)
	}

	return result
}

func (p *parser) parseIdent() ast.Expression {
	return &ast.Identifier{Token: p.currentToken, Literal: p.currentToken.Literal}
}

func (p *parser) parseNumber() ast.Expression {
	v, err := strconv.ParseInt(p.currentToken.Literal, 10, 64)
	if err != nil {
		p.addError("unable to convert number", INTERNAL_ERROR, p.currentToken.Start )
		return nil
	}

	return &ast.Number{ Token: p.currentToken, Value: v}
}

func (p *parser) parseFloat() ast.Expression {
	f, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		p.addError("unable to convert float", INTERNAL_ERROR, p.currentToken.Start )
		return nil
	}
	return &ast.Float{ Token: p.currentToken, Value: f}
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
		p.errorList = append(p.errorList, err)
	}
	p.nextToken = ntTok
}

func (p *parser) addError(msg, errorType string, pos token.Pos) {
	err := fmt.Errorf("SyntaxError: %s %s:%d:%d", msg, p.name, pos.Line, pos.Col)
	p.errorList = append(p.errorList, err)
}

// TODO: THINK HAVE A BETTER WAY IN REPRESENTING THIS FUNCTIONS
func processError(e []error) error {
	return nil
}
