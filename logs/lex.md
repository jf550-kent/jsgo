<div style="text-align: justify">

# Lexical analysis 
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

## Result

## Correctness
In the token package, the important feature we tested for is if a letter is a keyword or identifier [^]. In the Lexer, we created test to check for if the token reports the correct line and column in the source text [^]. Each operator and keyword has a test to correctly assert the correct token creation[^]. Due to the complexity of tokenising string literal we have a range of tests to check our lexer can correct convert the source text into the correct UTF-8 code points[^]. Lastly, there is integration test to check the whole JSGO program [^]. 

### Cost
Manual work needed to manually specify the correct test cases.

### Benefit

## Performance
We created a benchmark to measure the performance of the Lexer. It measure the speed and memory usage of the Lexer tokenistion the whole file [^].


### Cost
### Benefit

## Evaluation
One of the challenges when building the lexer is supporting UTF-8 representations for string literals. In our language, string literals are declared inside double quotes "". Initially, we assumed that supporting escape sequences with `\` would be a trivial task. Our goal was to support features such as:

|source|result|
|------|------|
|"\n"  | newline|
|"\t"  | tab |
|"\U0001F600" | 😀|
|"\u2603"|☃|

However, implementing this feature proved to be more complex than anticipated. We first needed to understand how our host language represents strings and then write additional code (l.convertString) to correctly convert these to UTF-8 code points. This added complexity includes maintaining the code and writing additional tests. At runtime, it increases the lexer's workload, as every string literal requires additional processing to convert escape sequences.

For new language implementers, it may be practical to omit this feature if your language does not require it.


|                     name                 | char       | line       |
|:-----------------------------------------|:----------:|:-----------|
| ../.goreleaser.yaml                      |       1133 |         46 |
| ../ast/ast.go                            |      12023 |        538 |
| ../ast/ast_test.go                       |       1421 |         26 |
| ../benchmark/benchmark.go                |        176 |          4 |
| ../benchmark/benchmark_test.go           |       4003 |        212 |
| ../bytecode/bytecode.go                  |       4242 |        194 |
| ../bytecode/bytecode_test.go             |       1981 |         86 |
| ../compiler/compiler.go                  |      10596 |        469 |
| ../compiler/compiler_test.go             |      27864 |       1009 |
| ../compiler/symbolTable.go               |       1792 |         84 |
| ../compiler/symbolTable_test.go          |       7138 |        297 |
| ../evaluator/builtin.go                  |        287 |         19 |
| ../evaluator/evaluator.go                |      13628 |        569 |
| ../evaluator/evaluator_test.go           |      11600 |        524 |
| ../evaluator/partial.go                  |       5069 |        207 |
| ../is.go                                 |       2015 |         90 |
| ../lexer/lexer.go                        |       6923 |        314 |
| ../lexer/lexer_test.go                   |      14890 |        294 |
| ../main.go                               |       2417 |        116 |
| ../object/builtin.go                     |        508 |         32 |
| ../object/environment.go                 |        915 |         48 |
| ../object/object.go                      |       5314 |        220 |
| ../parser/parser.go                      |      14616 |        649 |
| ../parser/parser_test.go                 |      20032 |        830 |
| ../run/main.go                           |       1183 |         65 |
| ../token/token.go                        |       3798 |        224 |
| ../token/token_test.go                   |       2775 |        133 |
| ../vm/frame.go                           |        402 |         24 |
| ../vm/vm.go                              |      17947 |        747 |
| ../vm/vm_test.go                         |       9999 |        425 |

</div>
