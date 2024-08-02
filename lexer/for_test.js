for (var i = 0; i < 5; i++) {

}

// Without initialization
let i = 0;
for (; i < 5; i++) {}

// 3. Loop Condition

// True condition
for (let i = 0; i < 5; i++) {}

// False condition (no execution)
for (let i = 5; i < 0; i++) {}

// 4. Loop Increment/Decrement

// Increment
for (let i = 0; i < 5; i++) {}

// Decrement
for (let i = 5; i > 0; i--) {}

// Custom increment
for (let i = 0; i < 5; i += 2) {}

// 5. Nested Loops

// Simple Nested Loop
for (let i = 0; i < 3; i++) {
    for (let j = 0; j < 2; j++) {
        console.log(i, j);
    }
}

// Deeply Nested Loops
for (let i = 0; i < 3; i++) {
    for (let j = 0; j < 2; j++) {
        for (let k = 0; k < 2; k++) {
            console.log(i, j, k);
        }
    }
}

// 6. Loop Control Statements

// `break` statement
for (let i = 0; i < 5; i++) {
    if (i === 3) break;
}

// `continue` statement
for (let i = 0; i < 5; i++) {
    if (i === 3) continue;
}

// `return` statement in a loop inside a function
function test() {
    for (let i = 0; i < 5; i++) {
        if (i === 3) return i;
    }
}

// 7. Infinite Loops

// Empty condition
for (;;) {}

// Always true condition
for (let i = 0; true; i++) {}

// Misconfigured loop (potential infinite loop)
for (let i = 0; i != 10; i += 2) {}

// 8. Invalid Syntax and Edge Cases

// Invalid Initialization
for (1 = 0; i < 5; i++) {} // Should throw an error

// Invalid Condition
for (let i = 0; i++; i < 5) {} // Should throw an error

// Invalid Increment
for (let i = 0; i < 5; i += ) {} // Should throw an error

// Variable Scope Outside the Loop
for (let i = 0; i < 5; i++) {}
console.log(i); // Should throw an error since 'i' is not defined

// 9. Loop with Non-Standard Iterables

// Loop over strings
for (let i = 0; i < "abc".length; i++) {}

// Loop over arrays
let arr = [1, 2, 3];
for (let i = 0; i < arr.length; i++) {}

// Loop over objects (should not work)
let obj = {a: 1, b: 2};
for (let i = 0; i < obj.length; i++) {} // Should throw an error or be invalid

// 10. Edge Cases

// Loop with no iterations
for (let i = 0; i < 0; i++) {}

// Loop where condition depends on external variable
let flag = false;
for (let i = 0; !flag; i++) {
    if (i === 3) flag = true;
}

// 11. Loop with `try/catch`

// Exception handling within the loop
for (let i = 0; i < 5; i++) {
    try {
        if (i === 3) throw new Error("test error");
    } catch (e) {
        console.log(e);
    }
}

// 12. Loop with asynchronous functions

// Async function in the loop
async function asyncFunc() {
    return new Promise((resolve) => setTimeout(resolve, 100));
}

for (let i = 0; i < 5; i++) {
    await asyncFunc();
}