## 1. Introduction
An interpreter may seem like a magic box until one understands how it works. Anyone who writes code might want to learn how to build one to uncover the inner workings of this "magic box". Understanding the basic mechanics of an interpreter is beneficial, as it provides valuable insights into a tool used daily. However, when we have decided to create a programming language, we are often motivated to quickly build out features during the initial phases, focusing solely on adding new capabilities. As a result, testing and measuring performance often take a backseat. Without adopting proper engineering principles and setting up robust engineering infrastructures, the project can become exponentially complex as new features are added. Ultimately, this leads to things breaking down without proper warning, leaving engineers unsure where to start troubleshooting the problem. Additionally, without adequate testing, engineers will have minimal confidence in making changes, fearing that their updates can potentially cause new issues. Consequently slowing down the engineering process and possibly renders the project difficult to contribute to.

This dissertation aims to apply a structured engineering approach to language implementations focusing on correctness, performance, and robust engineering infrastructures. Through adopting these principles, we can make better decisions throughout the interpreter development process. Thus, ensuring that every aspect of the language design and implementation is grounded in sound engineering practices. By applying structured engineering principles, we identify and analyze the trade-offs involved in different approaches, providing valuable insights for more efficient language development in the future. By evaluating the engineering efforts, challenges and benefits in implementing this approaches in interpreter contructions. 

Additionally, this paper intend to offer readers with no prior experience in language implementations a introduction to the fundamental concepts of building a programming language. By providing a broad overview of the process, we aim to enhance understanding and offer practical insights into the complexities of language design and interpreter construction.

Correctness means that the program produces the correct result for a given input, an essential principle in software development. In building an interpreter, correctness involves multiple critical aspects, such as properly tokenizing the source text, ensuring the parser constructs a well-formed abstract syntax tree (AST), and more. To guarantee these properties, various mechanisms will be employed, ranging from writing unit tests for specific cases to creating integration tests and developing a definition interpreter to verify that the AST is accurately constructed. In this paper, we will apply these mechanisms at various stages of building an interpreter.

Beyond ensuring that our interpreter meets its functional requirements, we must also focus on its performance. Efficient interpreters are essential because they minimise overhead and optimise resource usage such as memory, enhancing overall system performance. We focus on monitoring the performance of each component of the interperter to allow us to assess how code changes impact efficiency and identify performance issues. By understanding the performance of each component, we can pinpoint bottlenecks and prioritise their optimisation. For instance, we can build a partial evaluator to reduce the size of the AST tree.

One critical aspect of building our interpreter is establishing a engineering infrastructure. A specified engineering process provides a defined approach to programming language implementation, ensuring the project remains organised. By defining engineering processes in advance, such as code deployment, we reduce the effort needed to manage these tasks, allowing us to concentrate on building the project. This process also allows us to track code changes effectively, which helps us switch between features. Additionally, creating workflows for repetitive tasks to automate the tasks helps minimise human error and saves time, further allowing us to focus on building the interpreter.

Therefore, we aim to develop a programming language called JSGO $L_{\text{JSGO}}$. This language is a variant of JavaScript that can be executed using the NodeJS runtime. Thus, allowing us to verify the correctness of our interpreter with NodeJS. We will build the interpreters from scratch in Go, with no third-party dependencies. This means that we will use Go as the host language for constructing the interpreters, relying solely on the Go standard library. 

The components we are building include a lexer, parser, and abstract syntax tree, which will be used to construct our tree-walking interpreter. To enhance performance, we will also create a bytecode interpreter alongside our tree-walking interpreter to improve speed and memory efficiency.

## 2. Literature review
How can we build an interpreter with zero prior knowledge? Given that the author had no prior experience with language implementations, one of the contributions of this paper is to create a guide that provides an overview of language implementation.

To address this, we examined materials that teach language implementation from the ground up. We found useful resources, such as Crafting Interpreters [reference], Compilation Essentials [reference], and the Writing an Interpreter and Writing a Compiler series by Thorsten Ball. These materials were evaluated for their comprehensiveness and beginner-friendly approach. Ultimately, we chose Thorsten's series, which offers a step-by-step guide to building a language called Monkey. Although Monkey is designed for educational purposes and lacks certain features like for-loops and uses put for printing instead of console.log, it provided a solid foundation. We leveraged Thorsen's series to build our interpreter and extended it to meet our specific requirements. We also used part of the series's coding logic for our interperter when appropriate.

Additionally, inspired by Compilation Essentials [reference], we developed a definitional interpreter to verify the well-formedness of the AST and a partial evaluator for optimising the AST. The materials in Compilation Essentials provided a Python-based approach to building a definitional interpreter and partial evaluator, which we adapted to our needs.

## 3. Methodology

### 3.1 Goals
In this paper, we focus on correctness, performance, and engineering infrastructure in building an interpreter for $L_{\text{JSGO}}$. See below for $L_{\text{JSGO}}$'s specification. $L_{\text{JSGO}}$ is a variation of JavaScript that can be excuted by NodeJS. The end goal of the interpreter is to be able to parse and evaluate a subset of the benchmarks from "Are We Fast Yet" [Reference]. Successfully passing the benchmarks will serve as strong evidence that our interpreter functions correctly. Futhermore, we will leverage the NodeJS runtime to run the sets of benchmarks in conjunction of with our interpreters to validate the correctness of the results.

**Language features**
```
var apple = "apple";

```
Above is the specification for JSGO.

Given that the NodeJS runtime benefits from years of rigorous engineering, aiming to build a fully supported JavaScript interpreter is beyond the scope of this project. Therefore, we carefully selected the essential features to implement. This approach also necessitates the modification of the benchmark from "Are We Fast Yet"[REFERENCE] to suit our use cases [Reference the benchmark]. However, the core methodology for proving the correctness of our interpreter remains unchanged: our interpreter must run the same files as NodeJS and produce identical results.

We will build the interpreters using Go, leveraging insights from "Writing an Interpreter" [REFERENCE] and "Writing a Compiler" [REFERENCE]to guide our development process. We will also reference the Go compiler's source code, drawing on its design to inform how Go constructs its own compiler for inspirations of our interperter.

We have chosen Go as the host language for the implementation because it is a statically typed language that promotes productivity, even for those new to the language. Additionally, Go's standard library offers a rich set of tools, enabling engineers to accomplish tasks without relying on third-party dependencies. Especially, Go's built-in support for testing and benchmarking is highly convenient, allowing us to easily write unit tests and benchmark our codes. This makes Go an suitable choice for this project.

### 3.2 Strategic for Correctness
To ensure correctness in every feature developed in our interpreters, we plan to enforce unit tests at each step as features are merged into the main codebase. Additionally, at a higher level, we will ensure that every major feature undergoes integration testing. The table below details the steps for ensuring correctness. This approach provides confidence that new additions to the project work as expected and our existing implementations are not broken.

Additionally, we will build a definition interpreter to to check if the AST is well formed.

#### 3.2.1 Lexical analysis
In this phase, we will build a lexer capable of correctly tokenizing concrete text. Therefore, to ensure **correctness** of the lexer unit test we need to set up unit test on every character the interpreter plan to supports. For example, the source text "var" needs to correctly identify as a keyword token with the relevant informations. 

```
var = Token(Type: Keyword, Literal: "var", Start: Line: 1, Col: )
"hello world" = Token(Type: string, Literal: "hello world", Start: {Line: 1, Col: 1}, End: {Line:1, Col: 11})
```

#### 3.2.2 Syntactic analysis
The parser will be constructed based on the proposed language features described earlier. It will receive tokens from the lexer and transform them into an abstract syntax tree (AST). Ensuring the correctness of the parser’s output is crucial, as it must accurately represent the source text. To validate this, we will write tests for each supported language feature to verify that the parser produces the correct AST.

In this example, the parser should correctly build the var statement into VarStatementNode with the correct reference to the identifier and string expression.

```
var message = "hello world"; -> VarStatementNode: {Identifier: Identifier{message}, Expression: String{"hello word"}}
```

Additionally, we will build a definitional interpreter based on the proposed language features. This interpreter will be used to check any AST to confirm that it conforms to the $L_{JSGO}$'s grammar.

During the syntactic analysis phase, the parser must verify that the syntax of the source text is correct and report any syntax errors, including the precise location of the error in the source text. To achieve this, we will implement simple error handling that reports the first instance of a grammatical error. Therefore, we are required to write tests to verify that the parser accurately reports errors to users. This approach also distinguishes our implementation from Writing an interpreter.

#### 3.2.3 Tree walking interpreter
Once we have a complete AST, we can begin constructing the interpreter. We will build a tree-walking interpreter $Inter_{text{tree}}$, which will take the AST and evaluate the program.

To ensure the correctness of $Inter_{text{tree}}$, we will run tests for each new feature introduced by the interpreter. 
```
var apple = 9; 
console.log(apple); // The interpreter needs to print out 9.
```

For example, when evaluating a var statement, we need to verify that the identifier is correctly declared. We will write test cases to check if the interpreter can correctly evaluate such statements.

Furthermore, we will assess the behavior of the interpreter with various features. This includes checking how the interpreter handles function and variable declarations. For instance, the following code snippet should correctly declare the identifier apple and assign it the value 9, declare a function eat, and correctly print the result of eat(apple).

Example integration test:
```
var apple = 9;
var eat = function(a) {
    return a - 3;
};

console.log(eat(apple)); // The interpreter needs to print out 6.
```

#### 3.2.4 Bytecode interpreter
Compared to a tree-walking interpreter, a bytecode interpreter converts the program into a sequence of instructions encoded in bytecode, hence the name. As some argue, bytecode interpreters can be faster due to their optimized execution of bytecode instructions (Author, Year). To enhance our language, we plan to build a bytecode interpreter, $Inter_{bytecode}$. This interpreter will take the AST we previously constructed and transform it into bytecode instructions for execution.

Similar to the check for correctness with the $Inter_{tree}$, we will use the same strategy to test the interpreter's features. For instance, we test if the compiler creates the correct bytecode intrustions for a var statement and also checks if the virtual machine can correctly execute the instructions to access variable.

Finally, we will cross-examine the results of both interpreters by comparing their outputs. Furthermore, we will validate their correctness by comparing the results with those produced by the NodeJS runtime.

### 3.3 Strategic for Performance
In this dissertation, we focus on measuring the performance of every feature developed. Specifically, we will monitor both the memory usage and execution speed of each features. Throughout all phases of interpreter construction, we emphasise benchmarking key features to assess their performance effectively.

Furthermore, we propose to build a partial evaluator and a bytecode interpreter for improving the performance of our language implementations.

#### 3.3.1 Lexical analysis
In this phase, we should benchmark the performance of the lexer of tokenisising the entire source file.

#### 3.3.2 Syntactic analysis
Based of the language specification, we can benchmark the parser with all the supported features of the language. This will allows us to have an awareness of the performance of our parser. In terms of improving the performance of the parser, we will be building a partial evaluator (AST optimisation) to efficiently reduce the size of the AST.

#### 3.3.3 Tree walking interpreter
After the interpreter is fully built according to the specifications. We can start using the test suite we have specify to be ran with our $Inter_{tree}$  to check if the interpreter produces correct results. Furthermore, we will be using the Nodejs runtime to excute to ensure that our interpreter produce the same results as Nodejs.

With a complete interpreter according to our language specification we had set up, we can then benchmark our interpreter to check for the performance of our interpreter. This allows us to understand how does our tree walking interpreter is doing.

In the aspect of improving the performance of the $L_{JSGO}$ we plan to build a bytecode interpreter.


### 3.4 Engineering infrastrure
The feature development process begins with a feature proposal in GitHub Issues. A pull request is then created based on the feature proposal. Once the pull request is reviewed, it will be merged into the main branch. We will use Git to track code changes. For continuous integration, we will define GitHub Actions workflows to perform checks such as tests, linting, and static analysis. Additionally, we will create a workflow to automate the benchmarking process and store the results. For continuous deployment, we will leverage the tool Goreleaser to automate the deployment of our project.

#### 3.5 Benchmarks 
Talk about your benchmark and what new test you have added. I am not too sure how I want to word this part, that belongs to methodology by I will revisit this part after I have written the implementation of benchmarks. Basically I want to saidf that i have to change the syntax to allow for interpreter ot work AND i HAVE also added new test cases. 

## Word
concrete text == source text == source codes == human readable tex

## Implementations

## Lexer
In this phase, we need to tokensise the user defined source code. We have specific a huge range of tokens and JavaScript keywords. The rationale for this approach is that including a wide range of keywords incurs minimal additional code management, even if we later decide not to use some of them. The primary goals of the lexer are to tokenize numbers, floats, keywords, identifiers, strings, and relevant operators such as `+ - ! == !=`. Additionally, the lexer must accurately report the position of each token in the source text. To do that we have created 2 main components. 

First is `Token` (Code 1) that represent the smallest element of the language, the fields `Start` and `End` represent the position of the token at the file. While `Literal` is the actual string representation of the token. `TokenType` is the key field for the rest of the interpreter's component to identify different token type. See the actual [TokenType](https://github.com/jf550-kent/jsgo/blob/5415802df0edaffac116917f7d912354a860edee/token/token.go#L23C1-L86C2) definition. 

Secondly, we created `Lexer` (Code 1), which is the program that transform our user defined source text into tokens. The `Lexer` has field `src` that store the source text and the `src`'s index of the current pointer as `position` and next pointer as `nextPosition`. The field `line` and `col` represent the current character's position at the source file. The fields are the current character's position at the source file. 

```
// Token is the smallest element of the JSGO
type Token struct {
	TokenType TokenType // Type of the token
	Literal   string // Actual string representation in the source code
	Start     Pos // Start position of the token
	End       Pos // End position of the token
}

// Lexer tokenization the source text for the language
type Lexer struct {
	src          []byte // source text for tokenization
	position     int    // current position at [Lexer.src]
	nextPosition int    // next position to be lex at [Lexer.src]

	line int // the current line at the source text
	col  int // the current column at the source text

	ch byte // the current byte at [Lexer.src]
}
```
<center>Code 1</center>

There are 2 key method in the `Lexer` that do the tokenisation. In Go this is denoted as `func (l *Lexer)`, where the `*` in `*Lexer` is the pointer to the `Lexer` struct. Which are `Lex`[^1] (Code 2) and `next`[^2] (Code 2). 

The `next` method is responsible to correctly advance the Lexer pointer to the `l.src` and update the correct `col` and `line` position.

At the `Lex` method we skip the whitespace before we step into the switch statement. In the switch statement we check for the byte `l.ch` which store the current byte for the `Lexer` struct. We check for single byte by specific different case statement, if `l.ch` matches to `+` it will create a token with the `TokenType` of `ADD`. 

We will reach the default case when `l.ch` did not match any of our specified case statement. Here in the default case, we use the method `l.isLetter` to check if the current byte is a letter. If it is a letter the lexer will extract the letter then check if it is a keyword which then create a keyword token if it is a keyword if not an identifier token. 

If the `l.ch` is not a letter, we will then check if it is a digit. If it is a digit we will use the `l.getDigitToken()` [^3] to create an Number token or a Float token. We return an Illegal token if `l.ch` does not match an cases we have specify in our switch statement. 

```
func (l *Lexer) next() {
	if l.nextPosition == len(l.src) {
		l.ch = 0
		l.position = l.nextPosition
		return
	}
	l.position = l.nextPosition
	l.ch = l.src[l.position]
	switch l.ch {
	case '\n':
		l.col = 0
		l.line++
	default:
		l.col++
	}
	l.nextPosition++
}

func (l *Lexer) Lex() (token.Token, error) {
	var tok token.Token
	l.skipWhitespace()
	switch l.ch {
	case '+':
		pos := l.currentPos()
		tok = newToken(token.ADD, "+", pos, pos)
	
	...

	default:
		if l.isLetter() {
			start := l.currentPos()
			s, end := l.getLetter()
			if ty, ok := token.Keyword(s); ok {
				return newToken(ty, s, start, end), nil
			}
			return newToken(token.IDENT, s, start, end), nil
		}
		if l.isDigit() {
			return l.getDigitToken()
		}
		return newToken(token.ILLEGAL, "ILLEGAL", l.currentPos(), l.currentPos()), errors.New("ILLEGAL token")
	}
	l.next()
	return tok, nil
}
```
<center>Code 2: The implementation only show the important parts of the code.</center>

Therefore, the Lexer deals with the byte represntations of our source code. It checks each byte to transform it into tokens based on the `Lex` switch case algorithm we have created. Code 3 show the result of Lexer tokenising the the var statement.

```
var apple = 10;

// Result
Token(VAR, "var", {1, 1}, {1, 3})
Token(IDENT, "apple", {1, 5}, {1, 9})
Token(ASSIGN, "=", {1, 11}, {1, 11})
Token(NUMBER, "10", {1, 13}, {1, 14})
Token(SEMICOLON, ";", {1,15}, {1, 15})
```
<center>Code 3</center>



- [^1] `Lex` actual code implementation: (Here)[]
- [^2] `next` actual code implementation: (Here)[]

## Syntactic analysis
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

## Partial evaluator
The task of the partial evaluator is to evaluate the program with the static data trying to perform some optimistion up front in order to reduce the amount of operations at runtime.
var add = function (a) {
	var b = 8 + ((8 - 1) * 2); §§§§
	return b + a;
};

var add = function (a) {
	var b = 22; 
	return b + a;
};

In the above by performing the operation 8 + ((8 -1) * 2) up front and transforming it into a single AST node with value of 22 will efficiently saves memory and number of operations.  Imagine the add function is called 1000 times, the second will efficiently saves the memory for storing 3 binary node and skipped to perform the binary operations at runtime. One might argue this can help improve performance of the parser.

## Tree interpreter
Essentially, the source code defined by the developer is first transformed into an AST by the parser, which $Inter_{text{tree}}$ then uses to evaluate the program. In this case, $Inter_{text{tree}}$ can treat the AST as a standard tree traversal problem, evaluating the program by visiting each node at runtime.

During runtime, we need a way to represent values. Therefore, we plan to create an object system to represent these values, keeping them separate from the AST nodes. This approach helps maintain a clean separation between the AST node and object representation, with the object system being more lightweight compared to an AST node that contains syntactic information.

## Bytecode interpreter
The $Inter_{text{bytecode}}$ will consist of a compiler and a virtual machine. In this context, our compiler is a program that converts the AST into bytecode instructions specifically for use by the virtual machine at runtime. Unlike traditional compilers that produce artifacts such as executables, our compiler generates bytecode instructions on-the-fly without producing permanent files. Once the bytecode instructions are generated by the compiler, the virtual machine will execute them following the fetch-decode-execute cycle.

## Result
- Show the working of the interpereter
- code coverage
- 

## Evaluation
- code coverage
- Performance

In the code snippet ... is abrevaited. and foot node contain the lint to the actual code implementation hosted in github