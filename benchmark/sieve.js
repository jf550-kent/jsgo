function createSieve() {
  var flags = null;
  var size = 5000;

  var sieve = function (flags, size) {
    var primeCount = 0;

    for (var i = 2; i <= size; i += 1) {
      if (flags[i - 1]) {
        primeCount += 1;
        var k = i + i;
        while (k <= size) {
          flags[k - 1] = false;
          k += i;
        }
      }
    }
    return primeCount;
  };

  var benchmark = function () {
    flags = new Array(size);
    for (var i = 0; i < flags.length; i += 1) {
      flags[i] = true;
    }
    return sieve(flags, size);
  };

  function verify(result) {
    return result === 669;
  }

  return {
    benchmark: benchmark,
    verifyResult: verifyResult,
  };
}

var s = createSieve();
var result = s.benchmark();
s.verify(result);
