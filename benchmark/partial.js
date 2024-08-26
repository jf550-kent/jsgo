var add = function (a) {
	var b = 8 + ((8 - 1) * 2);
	return b + a;
};

for (var i = 0; i < 10000; i = i + 1) {
  add(9)
}
true;