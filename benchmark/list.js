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
