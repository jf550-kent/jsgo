function createMandelbrot() {
  function verifyResult(result, innerIterations) {
    if (innerIterations === 500) { return result === 191; }
    if (innerIterations === 750) { return result === 50; }
    if (innerIterations === 1) { return result === 128; }

    process.stdout.write('No verification result for ' + innerIterations + ' found\n');
    process.stdout.write('Result is: ' + result + '\n');
    return false;
  }

  // Define mandelbrot function
  function mandelbrot(size) {
    var sum = 0;
    var byteAcc = 0;
    var bitNum = 0;

    var y = 0;

    while (y < size) {
      var ci = ((2.0 * y) / size) - 1.0;
      var x = 0;

      while (x < size) {
        var zrzr = 0.0;
        var zi = 0.0;
        var zizi = 0.0;
        var cr = ((2.0 * x) / size) - 1.5;

        var z = 0;
        var notDone = true;
        var escape = 0;
        while (notDone && z < 50) {
          var zr = zrzr - zizi + cr;
          zi = 2.0 * zr * zi + ci;

          // preserve recalculation
          zrzr = zr * zr;
          zizi = zi * zi;

          if (zrzr + zizi > 4.0) {
            notDone = false;
            escape = 1;
          }
          z += 1;
        }

        byteAcc = (byteAcc << 1) + escape;
        bitNum += 1;

        // Code is very similar for these cases, but using separate blocks
        // ensures we skip the shifting when it's unnecessary, which is most cases.
        if (bitNum === 8) {
          sum ^= byteAcc;
          byteAcc = 0;
          bitNum = 0;
        } else if (x === size - 1) {
          byteAcc <<= (8 - bitNum);
          sum ^= byteAcc;
          byteAcc = 0;
          bitNum = 0;
        }
        x += 1;
      }
      y += 1;
    }
    return sum;
  }

  // Define innerBenchmarkLoop function
  function innerBenchmarkLoop(innerIterations) {
    return verifyResult(mandelbrot(innerIterations), innerIterations);
  }

  // Return object with methods
  return {
    innerBenchmarkLoop: innerBenchmarkLoop
  };
}