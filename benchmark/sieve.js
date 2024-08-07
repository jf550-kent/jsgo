var sieve = function (flags, size) {
  var primeCount = 0;

  for (var i = 2; i < size + 1; i = i + 1) {
    if (flags[i - 1]) {
      primeCount = primeCount + 1;
      for (var k = i + i; k < size + 1; k = k + i) {
        flags[k - 1] = false;
      }
    }
  }
  return primeCount;
};

var flags = [];
var size = 5000;
for (var ar = 0; ar < size; ar = ar + 1) {
  flags["push"](true)
}

var correct = sieve(flags, size) == 669;
correct;