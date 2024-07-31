function createElement(v, next) {
  return {
    val: v,
    next: next
  };
}

function getLength(element) {
  if (element.next === null) {
    return 1;
  }
  return 1 + getLength(element.next);
}

function createList() {
  function makeList(length) {
    if (length === 0) {
      return null;
    }
    var e = createElement(length, makeList(length - 1));
    return e;
  }

  function isShorterThan(x, y) {
    var xTail = x;
    var yTail = y;

    while (yTail !== null) {
      if (xTail === null) { return true; }
      xTail = xTail.next;
      yTail = yTail.next;
    }
    return false;
  }

  function tail(x, y, z) {
    if (isShorterThan(y, x)) {
      return tail(
        tail(x.next, y, z),
        tail(y.next, z, x),
        tail(z.next, x, y)
      );
    }
    return z;
  }

  function benchmark() {
    var result = tail(
      makeList(15),
      makeList(10),
      makeList(6)
    );
    return getLength(result);
  }

  function verifyResult(result) {
    return result === 10;
  }

  return {
    benchmark: benchmark,
    verifyResult: verifyResult
  };
}

var l = createList()
var result = l.benchmark()
console.log(l.verifyResult(result))
