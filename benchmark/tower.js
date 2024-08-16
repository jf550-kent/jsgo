var this_movesDone = 0;
var this_piles = [null, null, null];

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
buildTowerAt(0, 13)
moveDisks(13, 0, 1)
var correct = this_movesDone == 8191
correct;