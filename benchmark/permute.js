var count = 0;
var v = [0, 0, 0, 0, 0, 0];

var swap = function(i, j) {
  var tmp = v[i];
  v[i] = v[j];
  v[j] = tmp;
}

var permute = function(n) {
  count = count + 1;
  if (n != 0) {
    var nOne = n - 1;
    permute(nOne);
    for (var i = nOne; i > -1; i = i - 1) {
      swap(nOne, i);
      permute(nOne);
      swap(nOne, i);
    }
  }
}

permute(6)
var correct = count == 8660;
correct;
