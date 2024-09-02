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