<div style="text-align: justify">

# 5. Result
We successfully built $Inter_{bytecode}$ and $Inter_{tree}$ for our $L_{JSGO}$. The user can access the interpreter from [Github](https://github.com/jf550-kent/jsgo/releases). **Code 1** shows the command for the interpreter. The user need to specifiy the file name first and then the mode of interpreter to run the program. To enable the checker and partial evaluation you can pass in the optional command `debug`. 
**5 Code 2** shows how you can run the `example.js` file in different options.

```
./jsgo <filename> <tree|bytecode> [debug] [-version]
```
**5 Code 1**
```
./jsgo example.js tree 
./jsgo example.js bytecode
./jsgo example.js bytecode debug
```
**Code 2**

The language can support all the features mentioned in methodology section 3.1. See Code 3, a list of our langauge capability. In Line 1, we can declare a var statement and Lines 2 - 4 we can assign a value to an identifier. The data type our langauge supports are `null`, string literal, number, float, boolean. We can perform operation between number and float see Line 4. We support bitwise operation `^` and `&` for number see Lines 5 - 6. Lines 5 - 14, illustrates how our language supports binary operations and unary operations with the correct precedence. Interestingly, our string literal supports `/` escaping see Line 35 where we log an emoji. For composite data structure we support array, and objects. See Lines 26 - 29, for array declaration, index and the built in push method and accessing the array length. Lines 31 to 32 shows how to handle objects. Our language can supports recursion and closure see Lines 16 - 22. Lastly, we can define a for loop in Lines 24 and if/else in Lines 37 - 39.

```
var apple = 10;
apple = "Yellow";
apple = null;
apple = 80 + 9.0
8 ^ 7; 
8 & 9; 
!true; 
!false;
9 + 9;
9 - 1.0;
9 * 1.0 + 19;
9 / 10;
90 != 90;
9 == 3;

var recur = function(a) {
  if (a == 1) {
    return 1;
  }
  var b = a
  return b + recur(b - 1);
}

for (var i = 0; i < 10; i = i + 1) {}

var arr = [9, 9.0, true]
arr[0] // 9
arr["push"](10) // arr -> [9, 9.0, true, 10]
arr["length"] // 4

var obj = { "value": 9, "next": 9 }
obj["value"] // 9

var smile = "\u2603" 
console.log(smile) -> ☃

if (true) {
} else {
}
```
Code 3

Furthermore, all of the benchmarks mentioned in section 3.5 can be run by $Inter_{tree}$. Our partial evaluator can produce correct result when used with our interpreter. For $Inter_{bytecode}$ it can run all the benchmarks execpt for Mandelbrot and Sieve. This is because of constraints on the project, we were unable to create additional features to support Mandelbrot and Sieve. See Table for the result of our interpreter performance. 

Our benchmarks have rewritten and the implementation were tested. We end each benchmarks's with a boolean variable to identicate if the bencmark evaluate to pass or fail, so we log all the benchmarks's last line by running the node command. **Diagram 1** show the result of running all the benchmarks which logged all true. This means that our JS implementations is correct. 

![alt text](image-6.png)
**Diagram 1**

To check if a source code is part of $L_{JSGO}$ we created a definition interpreter to check the correctness of our AST. You can enable it by passing in the `debug` in the command. Diagaram show the command and result when you use `debug` and it will log out if the JS file you have provided is part of our language.
![alt text](image-7.png)

All the tests created have passed in this project. **Diagram 2** shows result of all passed tests when we run the tests.

![Example of running test](image-4.png)
**Diagram 2**

</div>