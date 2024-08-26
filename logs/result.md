- Write that we can parse utf character
- Undestanding the host language
- Show case your result for partial evaluator
- Demostrate how you can parse all the benchmark 
- In the bytecode it is hard to debug, so you need testing infrature 

Evaluation

# Result

## Interpreter demostration
You can download the interpreter from Github[]. To run the program type:

```
./jsgo <filename> <tree|bytecode> [debug] [-version]
```
The user need to specifiy the filename first and the mode of interpreter to run the program. To enable the checker and partial evaluation you can pass in the optional command `debug`. 

```
./jsgo <filename> <tree|bytecode> [debug] [-version]
```

Below is shows the user can run the `example.js` file with our interperter.
```
./jsgo example.js tree 
./jsgo example.js bytecode
./jsgo example.js bytecode debug
```

The language manage to support the feature mentioned in methodology. 

- string parsing

 ## Correctness
 All the tests created have passed in this project. See Diagram below for an example of running the specified package test.

![Example of running test](image-4.png)

## Present the result of the benchmark
The workflow mentioned in [Engineering Infrastructure] includes a step in the CI process where, upon merging to the main branch, performance benchmarks are automatically recorded. Below, we demonstrate how we managed to evaluate the benchmarks we created. The benchmark setup can be referenced [HERE](https://github.com/jf550-kent/jsgo/blob/main/benchmark/). These benchmarks were conducted on a Linux system with an AMD EPYC 7763 64-Core Processor, targeting the amd64 architecture.

| Benchmark                      | Count | Time (ns/op) | Memory (B/op) | Allocations (allocs/op) |
|--------------------------------|-------|--------------|---------------|-------------------------|
| BenchmarkListTree-4            | 82    | 12434228     | 3377390       | 75797                   |
| BenchmarkListTreeDebug-4       | 93    | 12426883     | 3377546       | 75799                   |
| BenchmarkListBytecode-4        | 122   | 9723327      | 3836623       | 86373                   |
| BenchmarkTowerTree-4           | 39    | 28821471     | 15799438      | 282907                  |
| BenchmarkTowerTreeDebug-4      | 42    | 28807471     | 15799611      | 282908                  |
| BenchmarkTowerBytecode-4       | 121   | 9912466      | 2703388       | 82652                   |
| BenchmarkMandelbrotTree-4      | 1     | 33679514731  | 3261057064    | 407626929               |
| BenchmarkMandelbrotTreeDebug-4 | 1     | 32879780784  | 3261052808    | 407626897               |
| BenchmarkPermuteTree-4         | 68    | 15844331     | 8315216       | 145896                  |
| BenchmarkPermuteTreeDebug-4    | 69    | 15886919     | 8315383       | 145898                  |
| BenchmarkPermuteBytecode-4     | 244   | 4915347      | 1763794       | 45090                   |
| BenchmarkSieveTree-4           | 145   | 8212330      | 1201700       | 108937                  |
| BenchmarkSieveTreeDebug-4      | 145   | 9026527      | 1201843       | 108939                  |
| BenchmarkQueensTree-4          | 96    | 11772284     | 5929917       | 123239                  |
| BenchmarkQueensTreeDebug-4     | 92    | 11792074     | 5930017       | 123169                  |
| BenchmarkQueensBytecode-4      | 307   | 3880841      | 1620185       | 37608                   |
Table: `ns` is nanosecond, `B` is bytes, `allocs` allocations

To understand the table's results, let’s examine the first item, BenchmarkListTree-4. This entry indicates that the benchmark was run with 4 cores (as denoted by -4 in the name), and the test was executed 82 times. Each execution took approximately 12,434,228 nanoseconds, used 3,377,390 bytes of memory, and allocated memory 75,797 times.

We managed to run all the benchmarks specified for our $Inter_{tree}$. Additionally, we enabled partial evaluation for each benchmark, as indicated by the name with ...Debug. As evident, some benchmarks benefit from the optimisation; however, there are cases where performance is actually slower. We will further evaluate these results in later sections.

For our $Inter_{bytecode}$, it is clear that there is a performance improvement over $Inter_{tree}$, even with partial evaluation enabled. However, due to time constraints on the project, we were unable to create additional features to support Mandelbrot and Sieve. We will further elaborate on this point in the evaluation section, discussing the trade-offs of building both a $Inter_{bytecode}$ and a $Inter_{tree}$.

To understand the table's results, let’s examine the first item, BenchmarkListTree-4. This entry indicates that the benchmark was run with 4 cores `-4` in the name, and the test was executed 82 times. Each execution took approximately 12,434,228 nanoseconds, used 3,377,390 bytes of memory, and allocated memory 75,797 times.

We manage to run all the benchmarks we have specifed for our $Inter_{tree}$, additionally we also enable partial evaluation for each of the benchmark the name with `...Debug`. As ecident some of the benchmark benefiot from the optimistation however, there is some that is in fact slower. We will further evaluate the result later in the sections. 

For our $Inter_{bytecode}$, it is obvious that ther is an imrpove performance from the $Inter_{tree}$ even with partial evaluation ebnle. However, due to time constraints of the project we were create additional feature to suport Mandelbrot and Sieve. We will further elborate this point in the evaluations about the tradeoff of buyulding a $Inter_{bytecode}$ and $Inter_{tree}$. 

I just to say is harder to develop and debug in bytecode because is harder to vislua at each stage of the developement because it is in bytecode. For else for the $Inter_{tree}$, it is much easier because the data structure is visulaise a tree structure. Therefore, bytecode interpereter implementer needs to invest in infrasterure to enable debugg better or building tools to vosialus the state at each points.