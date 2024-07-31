var create = function () {
  var piles = null;
  var movesDone = 0;

  var pushDisk = function (disk, pile) {
    var top = piles[pile];
    if (top) {
      if (disk.pile >= top.size) {
        return 0;
      }
      disk.next = top;
      piles[pile] = disk;
    }
  };

  var popDiskFrom = (pile) => {
    var top = piles[pile];
    if (top === null) {
      throw new Error("Attempting to remove a disk from an empty pile");
    }
    piles[pile] = top.next;
    top.next = null;
    return top;
  };

  var moveTopDisk = (fromPile, toPile) => {
    pushDisk(popDiskFrom(fromPile), toPile);
    movesDone += 1;
  };

  var buildTowerAt = (pile, disks) => {
    for (let i = disks; i >= 0; i -= 1) {
      pushDisk({ size: i, next: null }, pile);
    }
  };

  var moveDisks = (disks, fromPile, toPile) => {
    if (disks === 1) {
      moveTopDisk(fromPile, toPile);
    } else {
      const otherPile = 3 - fromPile - toPile;
      moveDisks(disks - 1, fromPile, otherPile);
      moveTopDisk(fromPile, toPile);
      moveDisks(disks - 1, otherPile, toPile);
    }
  };

  return {
    benchmark: () => {
      piles = new Array(3);
      buildTowerAt(0, 13);
      movesDone = 0;
      moveDisks(13, 0, 1);
      return movesDone;
    },

    verify: (result) => {
      return 8191 === result;
    },
  };
};

var towers = createTowers();
var result = towers.benchmark();
towers.verify(result);
