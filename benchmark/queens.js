function createQueens() {
  // Internal state
  var freeMaxs = null;
  var freeRows = null;
  var freeMins = null;
  var queenRows = null;

  var queens = function() {
    freeRows = [];
    freeMaxs = [];
    freeMins = [];
    queenRows = [];
    
    for (var i = 0; i < 8; i += 1) {
      freeRows[i] = true;
      freeMaxs[i] = true;
      freeMins[i] = true;
      queenRows[i] = -1;
    }

    return placeQueen(0);
  }

  var placeQueen = function(c) {
    for (var r = 0; r < 8; r += 1) {
      if (getRowColumn(r, c)) {
        queenRows[r] = c;
        setRowColumn(r, c, false);

        if (c === 7) {
          return true;
        }

        if (placeQueen(c + 1)) {
          return true;
        }
        setRowColumn(r, c, true);
      }
    }
    return false;
  }

  var getRowColumn = function(r, c) {
    return freeRows[r] && freeMaxs[c + r] && freeMins[c - r + 7];
  }

  var setRowColumn = function(r, c, v) {
    freeRows[r] = v;
    freeMaxs[c + r] = v;
    freeMins[c - r + 7] = v;
  }

  var benchmark = function() {
    var result = true;
    for (var i = 0; i < 10; i += 1) {
      result = result && queens();
    }
    return result;
  }

  var verify = function(result) {
    return result;
  }

  // Return object with methods
  return {
    benchmark: benchmark,
    verifyResult: verify
  };
}

