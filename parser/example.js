var num = 89;
var add = function(a, b) {
  return a + b;
};

var foo = function(a, func) {
  return a + func(a, a);
};

foo(num, add);

var total = num * 90;

