function createPermute() {
  var count = 0;
  var v = null;

  var permute = function(n) {
    count += 1;
    if (n !== 0) {
      var n1 = n - 1;
      permute(n1);
      for (var i = n1; i >= 0; i -= 1) {
        swap(n1, i);
        permute(n1);
        swap(n1, i);
      }
    }
  }

  var swap = function(i, j) {
    var tmp = v[i];
    v[i] = v[j];
    v[j] = tmp;
  }

  // Define benchmark function
  var benchmark = function() {
    count = 0;
    v = [];
    for (var i = 0; i < 6; i += 1) {
      v[i] = 0;
    }
    permute(6);
    return count;
  }

  // Define verifyResult function
  var verify = function(result) {
    return result === 8660;
  }

  // Return object with methods
  return {
    benchmark: benchmark,
    verify: verify
  };
}

var p = createPermute();
var result = p.benchmark();
p.verify(result);