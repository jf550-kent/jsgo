## Methodology
Every file will need to have 2 methods, benchmark and verify.

Calling the benchmark function will start executing the operation. Calling verify needs to return a bool, to determine if the benchmark 
produce the correct result.

To measure the performance of JSGO compiler, I would need to write a script to automating run all the benchmark and record the time.

We can further break down the benchmarking into running all the file. And each file.

Running all files, I will use need to write a script that runs all files and check the result.

Running indivisual files, will have similar steps compared with running all files.

every file needs to ends with console.log(p.verify(result)); if there is false then we can record that there is an excution errors.

Conditions of testing:
I want to run a file, and the last stadout needs to be a boolean, and true means that the program successfully run and produce a correct results. 

NodeJs vs JSGO vs Bun
The benchmark setup allows us to analayse the performance runtime of different compiler. With this method of development, it allows us to start obsovus how is our compiler performing compared to other compilers.

To ensure the correctness of our compiler we ensure that the result should be true. 