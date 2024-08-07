var freeMaxs = null;
var freeRows = null;
var freeMins = null;
var queenRows = null;

var result = true;

var queens = function() {
  freeRows = [true, true, true, true, true, true, true, true]
  freeMaxs = [true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true]
  freeMins = [true, true, true, true, true, true, true, true, true, true, true, true, true, true, true, true]
  queenRows = [-1, -1, -1, -1, -1, -1, -1, -1]

  return placeQueen(0);
}

var placeQueen = function(c) {
  for (var r = 0; r < 8; r = r + 1) {
    if (getRowColumn(r, c)) {
      queenRows[r] = c
      setRowColumn(r, c, false)

      if (c == 7) {
        return true;
      }

      if (placeQueen(c + 1)) {
        return true;
      }
      setRowColumn(r, c, true)
    }
  }
  return false;
}

var getRowColumn = function(r, c) {
  if (freeRows[r]) {
    if (freeMaxs[c + r]) {
      if (freeMins[c - r + 7]) {
        return true;
      }
    }
  }

  return false;
}

var setRowColumn = function(r, c, v) {
  freeRows[r] = v;
  freeMaxs[c + r] = v;
  freeMins[c - r + 7] = v;
}

for (var i = 0; i < 10; i = i + 1) {
  if (result) {
    if (queens()) {
      result = true;
    } else {
      result = false
    }
  } else {
    result = false
  }
}

result;