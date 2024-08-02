var a = 100;

if (a == 100) {
  var a = "ahdi";
  console.log(a); // "inside if"
};

var printA = function() {
  var a = "inside function";
  console.log(a); // "inside function"
};

console.log(a); 
