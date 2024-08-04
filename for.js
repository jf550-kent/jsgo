for (var i = 0; i < 5; i++) {

}

// Without initialization
var i = 0;
for (; i < 5; i++) {}

// 3. Loop Condition

// True condition
for (var i = 0; i < 5; i++) {}

// False condition (no execution)
for (var i = 5; i < 0; i++) {}

// 4. Loop Increment/Decrement



// 5. Nested Loops

// Simple Nested Loop
for (var i = 0; i < 3; i++) {
    for (var j = 0; j < 2; j++) {
        console.log(i, j);
    }
}

// Deeply Nested Loops
for (var i = 0; i < 3; i++) {
    for (var j = 0; j < 2; j++) {
        for (var k = 0; k < 2; k++) {
            console.log(i, j, k);
        }
    }
}

// `return` statement in a loop inside a function
function test() {
    for (var i = 0; i < 5; i++) {
        if (i === 3) return i;
    }
}

// 7. Infinite Loops

// Empty condition
for (;;) {}

// Always true condition
for (var i = 0; true; i++) {}

// Misconfigured loop (potential infinite loop)
for (var i = 0; i != 10; i += 2) {}

// 8. Invalid Syntax and Edge Cases

// Invalid Initialization
for (1 = 0; i < 5; i++) {} // Should throw an error

// Invalid Condition
for (var i = 0; i++; i < 5) {} // Should throw an error

// Invalid Increment
for (var i = 0; i < 5; i += ) {} // Should throw an error

// Variable Scope Outside the Loop
for (var i = 0; i < 5; i++) {}
console.log(i); // Should throw an error since 'i' is not defined

// 9. Loop with Non-Standard Iterables

// Loop over strings
for (var i = 0; i < "abc".length; i++) {}

// Loop over arrays
var arr = [1, 2, 3];
for (var i = 0; i < arr.length; i++) {}

// Loop over objects (should not work)
var obj = {a: 1, b: 2};
for (var i = 0; i < obj.length; i++) {} // Should throw an error or be invalid

// 10. Edge Cases

// Loop with no iterations
for (var i = 0; i < 0; i++) {}

// Loop where condition depends on external variable
var flag = false;
for (var i = 0; !flag; i++) {
  if (i === 3) {
    flag = true;
  }
}