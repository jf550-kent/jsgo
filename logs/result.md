- Write that we can parse utf character
- Undestanding the host language
- Show case your result for partial evaluator
- Demostrate how you can parse all the benchmark 
- In the bytecode it is hard to debug, so you need testing infrature 

Evaluation

# Result

## Interpreter demostration
You can download the interpreter from []. To run the program type []. The user need to specifiy the filename first and the mode of interpreter to run the program. To enable the checker and partial evaluation you can pass in the optional command `debug`. 

```
./jsgo example.js tree 
./jsgo example.js bytecode
./jsgo example.js bytecode debug
```

```
./jsgo <filename> <tree|bytecode> [debug] [-version]
```

```
// var statement
var apple = 90;

// function & object declaration
var create = function(value) {
  var apple = 90; // closure
  return {"value": value}
};

// supported operation
8 == -9;
9 != -9;
9.0 > 1.0;
9 < 10;
9.9 + 1;
10 - 1;
19 * 9;
10 /10;
90 << 9
10 ^ 9;

!false;
!true;

var arr = [3, 2, 9]
arr[1]; // array access
var obj = {"Hello": 90}
obj["Hello"] // object access

arr[1] = 90
arr["push"](90) // built in array function

null; // null value
create(90) // function call
console.log("Hello world") // built in printing to console

// if else
if (true) {
  90;
} else {
  9.10;
}
```


Above is the example file that our interpreter can parse, 

## Partial evaluation

## Definition interperter

## Present the result of the benchmark

Presnt benchmark result and past records


## Presenting the benchmark
The benchmarks were conducted on a Linux system with an AMD EPYC 7763 64-Core Processor, targeting the amd64 architecture.

Representing the result
Time line