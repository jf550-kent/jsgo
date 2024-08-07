var manBot = function(size) {
  var sum = 0;
  var byteAcc = 0;
  var bitNum = 0;

  for (var y = 0; y < size; y = y + 1) {
    var ci = ((2.0 * y) / size) - 1.0;
    for (var x = 0; x < size; x = x +1) {
      var zrzr = 0.0;
      var zi = 0.0;
      var zizi = 0.0;
      var cr = ((2.0 * x) / size) - 1.5;

      var done = false;
      var escape = 0;
      for (var z = 0; z < 50; z = z +1) {
        if (!done) {
          var zr = zrzr - zizi + cr;
          zi = 2.0 * zr * zi + ci;
  
          zrzr = zr * zr;
          zizi = zi * zi;
  
          if (zrzr + zizi > 4.0) {
            done = true;
            escape = 1;
          }
        }
      }

      byteAcc = (byteAcc << 1) + escape
      bitNum = bitNum + 1;
      
      if (bitNum == 8) {
        sum = sum ^ byteAcc
        byteAcc = 0
        bitNum = 0
      } 

      if (x == size - 1) {
        byteAcc = byteAcc << (8 - bitNum)
        sum = sum ^ byteAcc
        byteAcc = 0
        bitNum = 0
      }
    }
  }

  return sum;
}

var correct = false;
if (manBot(500) == 191) {
  if (manBot(750) == 50) {
    if (manBot(1) == 128) {
      correct = true
    }
  }
}

correct;