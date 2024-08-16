# JSGO 🔥🔥🔥
JavaScript compiler built in go

## **Installation**

In this directory, you will find the zip `jsgo_<machine spec>.tar.gz` that you can download based on your machine specification. Download the correct one and unzip the file, it will contain the binary that you can use to run the interpreter.

If you have the go compiler, you can clone the repo: https://github.com/jf550-kent/jsgo.git and build from source. 

Alternatively, you can visit the release page to download the latest version of this project: https://github.com/jf550-kent/jsgo/releases

## **Usage**
```
./jsgo <filename> <tree|bytecode> [debug] [-version]
```
The order of the argument provided matters.

In the command you must specific the filename first with an extension of `.js`. 

Secondly, which mode of the interpreter you want to use. This project supports only `tree` or `bytecode`. 

Optional: you can pass in `debug`. This mode, will pass in a abstract syntax tree optimization and checker that checks if the program is well formed. Both mode are performed in a best effort, will default to action the action that does not crash the program.

Getting the version of the interpreter
```
./jsgo --version
```

Here is the example command correct command that works:
```
./jsgo ./benchmark/list.js tree debug
./jsgo ./benchmark/mandelbrot.js tree debug
./jsgo ./benchmark/permute.js tree debug
./jsgo ./benchmark/queens.js tree debug
./jsgo ./benchmark/sieve.js tree debug
./jsgo ./benchmark/tower.js tree debug


./jsgo ./benchmark/list.js bytecode debug
./jsgo ./benchmark/permute.js bytecode debug
./jsgo ./benchmark/queens.js bytecode debug
./jsgo ./benchmark/tower.js bytecode debug
```

## **Software enginering pratice**
This project uses the issues tracker: https://github.com/jf550-kent/jsgo/issues for feature development and fixing errors.
To contribute to the project, you need to submit PR: https://github.com/jf550-kent/jsgo/pulls. This is also how the project handle development.

## **Project walk through**
- This `performance/` directory stores the result of the benchmark
- `token/` & `lexer/` & `ast/` & `parser/` this four directory is the code for the frontend of the interpreter.
- `evaluator/` is the tree walking interpreter that evalualate the result after the parser builts a ast.
- `vm/` similar to the `evaluator/` vm is the directory that execute the stack-based bytecode intructions.
- `compiler/` compiles a AST into bytecode intructions for the `vm/` to run.
- `benchmark/` stores the benchmarking files and script to run for benchmarking
- `.github/` is used for Continuous integration 
- `.goreleaser.yaml` is config used for  Continuous deployment
- `main.go` is the entry file for this project
- `specification.md` is the specifiction of the JSGO language
