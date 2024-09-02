## 4.2 Syntactic analysis
Building an AST allows us to represent the source text in a data structure that is easier for the interpreter to work with. In this phase, we created two main packages: `parser` and `ast`. The parser is responsible for taking the tokens and transforming them into an AST. The `ast` package specifies all the abstract syntax nodes available in $L_{JSGO}$. In $L_{JSGO}$, an expression evaluates to a value, whereas a statement executes an action without necessarily producing a value; depending on the context, it can have side effects. The statement types include ExpressionStatement, VarStatement, AssignmentStatement, ReturnStatement, ForStatement, and BlockStatement. For expressions, we have Number, Float, Identifier, Boolean, Function, Null, String, Array, Dictionary, BinaryExpression, IfExpression, UnaryExpression, CallExpression, Index, and BracketDeclaration. Each type of expression and statement represents a different kind of node. Some nodes have other nodes as their children. For instance, VarStatement has child nodes of Identifier and Expression.

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

The parser also needs to report any syntax errors, including the precise location of the error in the source text. We have implemented simple error handling that reports the first instance of a grammatical error. One of an additional features from Writing an interpreter.

Building on the var example (Code 3), the parser will need to take in the stream of token and build an var statement node and expression statement node. [Need to create diagram]
```
Token(VAR, "var", {1, 1}, {1, 3})
Token(IDENT, "apple", {1, 5}, {1, 9})
Token(ASSIGN, "=", {1, 11}, {1, 11}) ->  
Token(NUMBER, "10", {1, 13}, {1, 14})
Token(SEMICOLON, ";", {1,15}, {1, 15})
``` 

| Node                    | Description |
|-------------------------|-------------|
| Main                    | Main is node represent the $L_{JSGO}$ program |
| VarStatement            | var declaration node |
| ReturnStatement         | return statement node |
| BlockStatement          | block statement node, contains slice of statements |
| ExpressionStatement     | expression wrap into a statement |
| AssignmentStatement     | assign statement node |
| ForStatement            | for statement node |
| Number                  | Number expression node |
| Identifier              | Identifier expression node |
| Float                   | Float expression node |
| Boolean                 | Boolean expression node |
| IFExpression            | If expression node |
| BinaryExpression        | Binary expression node |
| UnaryExpression         | Unary expression node |
| FunctionDeclaration     | Function declaration node |
| CallExpression          | Call expression node |
| String                  | String expression node |
| Array                   | Array expression node |
| Index                   | Index expression node |
| Null                    | Null expression node |
| Dictionary              | Dictionary expression node|


#### 4.2.1 Correctness
The parser receive tokens from the lexer and transform them into an abstract syntax tree (AST). Therefore, we need to ensure the correctness the result of our parser as it needs to accurately represent the source text. To validate this, we write tests for every language features to verify that the parser produces the correct AST.

For example, the `TestVar` test checks if the parser can correctly build a VarStatement node. The first test case in `var apple = "Hello world"` at line 7, checks if the parser can creates the correct identifier with the name `apple` Line 19 and checks if the expression is string "Hello world".

```
func TestVar(t *testing.T) {
	tests := []struct{
		input              string
		expectedVariable   string
		expectedExpression any
	}{
		{"var apple = 10;", "apple", 10},
		...
		{"var apple = \"Hello world\";", "apple", "Hello world"},
	}

	for _, tt := range tests {
		main := Parse("", []byte(tt.input))
		if len(main.Statements) != 1 {
			t.Errorf("main should have 1 statement. got=%d", len(main.Statements))
		}

		varStmt := checkStatement[*ast.VarStatement](t, main.Statements[0])
		if varStmt.Variable.Literal != tt.expectedVariable {
			t.Errorf("wrong variable. expected=%s, got=%s", tt.expectedVariable, varStmt.Variable.Literal)
		}
		testValueExpression(t, varStmt.Expression, tt.expectedExpression)
	}
}
```

Additional tests have been created to examine various nodes, located [here](https://github.com/jf550-kent/jsgo/blob/main/parser/parser_test.go). The full list of tests are TestParsingDictionaryDecl, TestParsingDictionaryWithExpressions, TestParsingDictionaryIntegerKeys, TestParsingDictionaryBooleanKeys, 
TestParsingDictionarysStringKeys, TestParsingEmptyDictionary, TestParsingIndexString, TestParsingIndex, TestParsingArray, TestParsingEmptyArray, TestForExpression, TestAssignmentStatement, TestExpressionStatement, TestIfElseExpression,  TestIfExpression, TestCallExpressionNoArgument, TestCallExpressionParsing, TestOperatorPrecedenceParsing, TestFunctionDeclaration, TestUnaryExpression, TestBinaryExpression, TestReturn. The subword after Test of each test names shows the specifc feeature being tested. 

#### 4.2.2 Definition interpreter
Additionally, we will build a definitional interpreter based on the proposed language features. This interpreter will be used to check the correctness of any AST to confirm that it conforms to the $L_{JSGO}$'s grammar. In our case, we applied this concept to build a checker that validates whether an Abstract Syntax Tree (AST) belongs to $L_{JSGO}$. 

The checker uses a recursive function to evaluate the AST. For each node in the AST, the function checks if it adheres to the grammar of $L_{JSGO}$. For instance, for a binary expression node we check if the operator is a valid set of operator Line 7, only if it is part of the permmited operator we can check if the child lefe and right node, else we conclude that the AST is not part of $L_{JSGO}$ by return false. See the code implementation [here](https://github.com/jf550-kent/jsgo/blob/main/is.go)

```
func check(node ast.Node) bool {
	switch node := node.(type) {
	case *ast.Number, *ast.Float, *ast.Boolean, *ast.Null, *ast.String, *ast.Array, *ast.Dictionary, *ast.Identifier:
		return true
	case *ast.BinaryExpression:
		switch node.Operator {
		case "+", "-", "*", "/", "<<", "^", "<", ">", "==", "!=":
			return check(node.Left) && check(node.Right)
		default:
			return false
		}
		...
	}
	...
}
```
This feature is enabled in debug mode, where we can check if the parser builds an AST that belongs to $L_{JSGO}$.

#### 4.2.2 Performance
We wrote a benchmark that has a similar implementation as the lexer benchmark named BenchmarkExample. The function is located [here](https://github.com/jf550-kent/jsgo/blob/25b7d13e5763e36168dbf25d0d044d6dcaad88f7/parser/parser_test.go#L12).