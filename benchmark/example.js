var apple = 9;
var appleExample = 10;

apple + appleExample;

var add = function(a) {
  return a + a;
};

add(apple, appleExample);

var createNode = function(v, n) {
  return { "value": v, "next": n };
}

var length = function(node) {
  if (node["next"] == null) {
    return 1;
  }
  return 1 + length(node["next"]);
}

var createList = function(size) {
  if (size == 0) { return null; }
  return createNode(size, createList(size - 1));
}

var tail = function(x, y, z) {
  if (isShorterThan(y, x)) {
    return tail(tail(x["next"], y, z), tail(y["next"], z, x), tail(z["next"], x, y));
  }
  return z;
}

var isShorterThan = function(x, y) {
  var xTail = x;
  var yTail = y;

  for (;yTail != null;) {
    if (xTail == null) { return true; }
    xTail = xTail["next"];
    yTail = yTail["next"];
  }
  return false;
}

var result = length(tail(createList(15), createList(10), createList(6)))
var correct = result == 10;
correct;
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

correct;var count = 0;
var v = [0, 0, 0, 0, 0, 0];

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

var swap = function(i, j) {
  var tmp = v[i];
  v[i] = v[j];
  v[j] = tmp;
}

permute(6)
var correct = count == 8660;
correct;
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

result;var sieve = function (flags, size) {
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
correct;var this_movesDone = 0;

var pushDisk = function(disk, pile) {
  var top = this_piles[pile]
  if (top) {
    if (disk["size"] > top["size"] + 1) {
      return "Cannot put a big disk on a smaller one";
    }
  }

  disk["next"] = top
  this_piles[pile] = disk
}

var createTowerDisk = function (size) {
  return { "size": size, "next": null };
}

var buildTowerAt = function (pile, disks) {
  for (var d = disks; d > -1; d = d -1) {
    pushDisk(createTowerDisk(d), pile)
  }
}

var popDiskFrom = function (pile) {
  var top = this_piles[pile]
  if (top == null) {
    return "Trying to remove a empty pile";
  }
  this_piles[pile] = top["next"]
  top["next"] = null
  return top;
}

var moveTopDisk = function (fromPile, toPile) {
  pushDisk(popDiskFrom(fromPile), toPile)
  this_movesDone = this_movesDone + 1
}

var moveDisks = function (disks, fromPile, toPile) {
  if (disks == 1) {
    moveTopDisk(fromPile, toPile);
  } else {
    var otherPile = (3 - fromPile) - toPile;
    moveDisks(disks - 1, fromPile, otherPile);
    moveTopDisk(fromPile, toPile);
    moveDisks(disks - 1, otherPile, toPile);
  }
}
var this_piles = [null, null, null];
buildTowerAt(0, 13)
moveDisks(13, 0, 1)
var correct = this_movesDone == 8191
correct;
