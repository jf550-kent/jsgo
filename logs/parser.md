<div style="text-align: justify">

# Syntactic analysis
Building an AST allows us to represent the source text in a data structure that is easier for the interpreter to work with. In this phase, we created two main packages: `parser` and `ast`. The parser is responsible for taking the tokens and transforming them into an AST. The `ast` package specifies all the abstract syntax nodes available in $L_{JSGO}$. In the grammar of $L_{JSGO}$, $L_{JSGO}$ is essentially an array of statements. In $L_{JSGO}$, an expression evaluates to a value, whereas a statement executes an action without necessarily producing a value; depending on the context, it can have side effects. The statement types include ExpressionStatement, VarStatement, AssignmentStatement, ReturnStatement, ForStatement, and BlockStatement. For expressions, we have Number, Float, Identifier, Boolean, Function, Null, String, Array, Dictionary, BinaryExpression, IfExpression, UnaryExpression, CallExpression, Index, and BracketDeclaration. Each type of expression and statement represents a different kind of node. Some nodes have other nodes as their children. For instance, VarStatement has child nodes of Identifier and Expression.


We use a recursive descent parser to build the AST. In the parser package, we have one exported function `Parse`[^]. This function takes in the file name and source code in the form of byte slice. The `Parse` is a function that sets the lexer and parser, which are to recursivly build an abstract syntax tree with the root node as Main node (Code). A Main node represent the a JSOG program. In this package, we have the `parser`[^] struct. We highlight the `parser` struct's key method. Which are `p.parse()` (Code) [^] and `p.parseExpression(precedence int)`(Code) [^]. The `parser` struct's method `parse` parse a single statement which is used in the `Parse` function in a for to be called until there is a EOF token (Code). The `p.parseExpression(precedence int)` (Thuston) function is the main function that recursivly build any expression node.

```
func Parse(filename string, src []byte) *ast.Main {
  ... // lexer and parser setup code

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
func (p *parser) parse() ast.Statement {
	switch p.currentToken.TokenType {
	case token.VAR:
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IDENT:
		switch p.nextToken.TokenType {
		case token.ASSIGN:
			return p.parseAssignmentStatement()
		}
	case token.FOR:
		return p.parseForStatement()
	}
	return p.parseExpressionStatement()
}

func (p *parser) parseExpression(precedence int) ast.Expression {
	unaryFunc, ok := p.unaryExpressionFuncs[p.currentToken.TokenType]
  ...
	result := unaryFunc()
	for !p.peekExpect(token.SEMICOLON) && precedence < p.peekPred() {
		binaryFunc, ok := p.binaryExpressionFunc[p.nextToken.TokenType]
    ...
		p.next()
		result = binaryFunc(result)
	}
	return result
}
```

Building on the var example (Code 3), the parser will need to take in the stream of token and build an var statement node and expression statement node. 
```
Token(VAR, "var", {1, 1}, {1, 3})
Token(IDENT, "apple", {1, 5}, {1, 9})
Token(ASSIGN, "=", {1, 11}, {1, 11}) ->  
Token(NUMBER, "10", {1, 13}, {1, 14})
Token(SEMICOLON, ";", {1,15}, {1, 15})
``` 

## Result

## Correctness
-> AST important for string no test
-> Parser write alot of tests
Create a table on the 

### Cost

### Benefit

## Performance


### Cost
### Benefit

## Evaluation
- Should I talk about making a clear distinction of between expression and statement on you language implementations
  - It reduce the confusion of further down the road

</div>

```
# JSGO definition
Expression ::= Number | Float | Identifier | Boolean | Function | Null | String | Array | Dictionary |  BinaryExpression | IfExpression | UnaryExpression  | CallExpression | Index | BracketDeclaration

Number = 891 (int64)
Float = 89.1 (float64)
Identifier = foo
Boolean = true | false
Function = function ([Identifier]) { [BlockStatement] }
Null = null
String = "hello world"
Array = []Expression
Dictionary = {<Expression>: <Expression>} Key only support string
BinaryExpression = <Expression> (+ | - | * | / | != | == | << | ^)  <Expression>, The << | ^ only support for number
IfExpression = if (<Expression>) {[BlockStatement]} [esle {[BlockStatement]}]
UnaryExpression =  (! <Expression>) | (- <Number | Float>)
CallExpression = <Expression>([]Expression) // In this case the parentheness is the literal, not grouping
Index = <Expression>[<Expression>] // In this case the [] is the literal, not optional
BracketDeclaration = <Expression>[<Expression>] = <Expression> // In this case the [] is the literal, not optional

Statement ::= ExpressionStatement | VarStatement | AssignmentStatement | ReturnStatement | ForStatement | BlockStatement

ExpressionStatement = <Expression> // ExpressionStatement is a statement that wraps an expression. 
VarStatement = var Identifier = <Expression>[;]
ReturnStatement = return <Expression>;
AssignmentStatement = Identifier = <Expression>[;]
ForStatement = for ([<VarStatement>]; <Expression>; [<AssignmentStatement>]) { [BlockStatement] }
BlockStatement = []Statement

JSGO ::= Main == []Statement
```
Grammar 1 
