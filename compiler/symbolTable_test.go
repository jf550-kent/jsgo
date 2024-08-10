package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"apple":  {Name: "apple", Scope: GlobalScope, Index: 0},
		"banana": {Name: "banana", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	apple := global.Define("apple")
	if apple != expected["apple"] {
		t.Errorf("symbol apple incorrect got=%+v", apple)
	}

	banana := global.Define("banana")
	if banana != expected["banana"] {
		t.Errorf("symbol banana incorrect got=%+v", banana)
	}
}

func TestResolve(t *testing.T) {
	global := NewSymbolTable()
	global.Define("apple")
	global.Define("banana")

	expected := []Symbol{
		{Name: "apple", Scope: GlobalScope, Index: 0},
		{Name: "banana", Scope: GlobalScope, Index: 1},
	}

	for _, symbl := range expected {
		result, ok := global.Resolve(symbl.Name)
		if !ok {
			t.Errorf("name %s not resolvable", symbl.Name)
		}
		if result != symbl {
			t.Errorf("expected %+v to resolve to %+v", symbl, result)
		}
	}
}
