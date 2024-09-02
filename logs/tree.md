## 4.3 Tree walking interpreter
The source code defined by the developer is first transformed into an AST by the parser, which is used by $Inter_{tree}$ to evaluate the program. In this case, $Inter_{tree}$ can treat the AST as a standard tree traversal problem, evaluating the program by visiting each node at runtime. During runtime, we need a way to represent values. Therefore, we created an object system to represent these values, keeping them separate from the AST nodes. The object system are our runtime value, the definition is in the `object` package. This approach helps maintain a clean separation between the AST node and object representation, with the object system being more lightweight compared to an AST node that contains syntactic information.

Our $Inter_{tree}$ is an `Eval` function that takes in an AST node `Main` and a boolean debug flag to opt-in for partial evaluation. Within `Eval`, the top-level environment is created, which is then passed to the private eval function. An environment in our interpreter is a crucial structure that maps variable names to their corresponding values or functions during runtime, ensuring correct evaluation within the appropriate scope. Thus in our code, we implemented the [`Environment`](https://github.com/jf550-kent/jsgo/blob/5415802df0edaffac116917f7d912354a860edee/object/environment.go#L5) [FOOTNOTE] struct with a map field. The eval function is the main recursive function that traverses the AST. In Go, we distinguish whether a declaration is exported by capitalising the first character of its name. This is why Eval is exported in the package, whereas eval is not. Refer to the [code](https://github.com/jf550-kent/jsgo/blob/main/evaluator/evaluator.go) [FOOT] for eval to examine how the function evaluates the abstract syntax tree. [Footnote]

We will highlight one example of how does the function evaluate the tree. In the diagram, the interpreter evaluates the Abstract Syntax Tree (AST) as follows: In Step 1, it encounters the vsar statement node and evaluates the binary expression 5 + 5 first to get 10 (Step 2). This result is then bind to apple in Step 4. Next, the interpreter moves to Step 5 the var basket = 8 * apple; statement. In Step 6, it evaluates the binary expression 8 * apple, where apple is looked up (value 10), and the calculation 8 * 10 results in 80, which is then bind to basket. Therefore, by the end, apple holds the value 10, and basket holds the value 80.

```
var apple = 5 + 5;
var basket = 8 * apple;
```


![Tree walking](image.png)

### 4.3.1 Correctness
To ensure the correctness of our tree walking interpreter, we wrote tests that check if the interpreter can evaluate to the correct results.

| Test Name                     | Description                                |
|:------------------------------|:------------------------------------------|
| TestVarStatement               | Check var statment is correctly declared |
| TestUnaryOperation             | Check the prefix ! and -                  |
| TestEvalNumberExpression       | Check numbers operation is evaluate to the correct value |
| TestEvalBooleanExpression      | Check the interpreter can boolean condition such as <, !=, ==, >  |
| TestBangOperator               | Check !<expression> operation                 |
| TestIfElseCondition            | Check evaluation of if else                |
| TestReturnStatements           | Check the interpreter can return the correct value   |
| TestErrorHandling              | Check if the interpreter can report runtime error such as type mistach |
| TestFunctionObject             | Check if function can be declared and convert to the runtime representation of function |
| TestEnclosingEnvironments      | Check if function closure                  |
| TestFunctionApplication        | Check function call can evaluate to the correct value |
| TestAssignment                 | Check if value to correct assign           |
| TestFor                        | Check for loop                |
| TestArrayLiterals              | Check if array is correctly declared and converted to the runtime representations of array               |
| TestArrayIndexExpressions      | Check if the interpreter can access its elements |
| TestArrayLength                | Check if the array can report the correct size |
| TestArrayPush                  | Check if the push method works for array                |
| TestArrayFunctionCall          | Check if the interpreter can call an element of an array that is a function                 |
| TestDictionaryExpressions      | Check if the interpreter can access object's element         |
| TestDictionaryDeclaration      | Check if object declaration is correctly converted to runtime representation of object|
| TestClosure                    | Test closure             |

All the benchmarks created in section 3.5 ends with a boolean variable that check if the program is correct. We ran the $Inter_{tree}$ with all the benchmarks, and the boolean evaluate to true. This shows that our $Inter_{tree}$ is able to a diverse set of programs.

### 4.3.2 Performance
After the interpreter is fully built according to the specifications. We used the benchmarks in section 3.5 to benchmark the performance of our $Inter_{tree}$. 

To further improve performance we also built a partial evaluator [^]. The goal of the partial evaluator is to reduce the size of the tree. It evaluate the program with the static data availbale to perform some optimistion up front in order to reduce the amount of operations at runtime.
```
// user defined add function
var add = function (a) {
	var b = 8 + ((8 - 1) * 2);
	return b + a;
};

// Transformed add function
var add = function (a) {
	var b = 22; 
	return b + a;
};
```
In the above by performing the operation 8 + ((8 -1) * 2) up front and transforming it into a single AST node with value of 22 will efficiently saves memory and number of operations.  Imagine the add function is called 1000 times, the second transformed add function will efficiently saves the memory for storing 3 binary node and skipped to perform the binary operations at runtime.