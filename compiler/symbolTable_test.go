package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"apple":   {Name: "apple", Scope: GlobalScope, Index: 0},
		"banana":  {Name: "banana", Scope: GlobalScope, Index: 1},
		"charlie": {Name: "charlie", Scope: LocalScope, Index: 0},
		"delta":   {Name: "delta", Scope: LocalScope, Index: 1},
		"echo":    {Name: "echo", Scope: LocalScope, Index: 0},
		"fox":     {Name: "fox", Scope: LocalScope, Index: 1},
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

	local := NewEnclosedSymbolTable(global)
	c := local.Define("charlie")
	if c != expected["charlie"] {
		t.Errorf("expected charlie=%+v, got=%+v", expected["charlie"], c)
	}

	d := local.Define("delta")
	if d != expected["delta"] {
		t.Errorf("expected delta=%+v, got=%+v", expected["delta"], d)
	}

	secondLocal := NewEnclosedSymbolTable(local)
	e := secondLocal.Define("echo")
	if e != expected["echo"] {
		t.Errorf("expected echo=%+v, got=%+v", expected["echo"], c)
	}

	f := secondLocal.Define("fox")
	if f != expected["fox"] {
		t.Errorf("expected fox=%+v, got=%+v", expected["fox"], d)
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

func TestResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("apple")
	global.Define("banana")

	local := NewEnclosedSymbolTable(global)
	local.Define("charlie")
	local.Define("delta")

	expected := []Symbol{
		{Name: "apple", Scope: GlobalScope, Index: 0},
		{Name: "banana", Scope: GlobalScope, Index: 1},
		{Name: "charlie", Scope: LocalScope, Index: 0},
		{Name: "delta", Scope: LocalScope, Index: 1},
	}

	for _, symbl := range expected {
		if symbl.Scope == GlobalScope {
			result, ok := global.Resolve(symbl.Name)
			if !ok {
				t.Errorf("name %s not resolvable", symbl.Name)
			}
			if result != symbl {
				t.Errorf("expected %+v to resolve to %+v", symbl, result)
			}
			break
		}

		result, ok := local.Resolve(symbl.Name)
		if !ok {
			t.Errorf("name %s not resolvable", symbl.Name)
		}
		if result != symbl {
			t.Errorf("expected %s to resolve to %+v, got=%+v",
				symbl.Name, symbl, result)
		}
	}
}

func TestResolveNestedLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("apple")
	global.Define("banana")

	local := NewEnclosedSymbolTable(global)
	local.Define("charlie")
	local.Define("delta")

	innerLocal := NewEnclosedSymbolTable(local)
	innerLocal.Define("echo")
	innerLocal.Define("fox")

	tests := []struct {
		table           *SymbolTable
		expectedSymbols []Symbol
	}{
		{
			local,
			[]Symbol{
				{Name: "apple", Scope: GlobalScope, Index: 0},
				{Name: "banana", Scope: GlobalScope, Index: 1},
				{Name: "charlie", Scope: LocalScope, Index: 0},
				{Name: "delta", Scope: LocalScope, Index: 1},
			},
		},
		{
			innerLocal,
			[]Symbol{
				{Name: "apple", Scope: GlobalScope, Index: 0},
				{Name: "banana", Scope: GlobalScope, Index: 1},
				{Name: "echo", Scope: LocalScope, Index: 0},
				{Name: "fox", Scope: LocalScope, Index: 1},
			},
		},
	}

	for _, tt := range tests {
		for _, sym := range tt.expectedSymbols {
			result, ok := tt.table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name not resolvable: %s", sym.Name)
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%v", sym.Name, sym, result)
			}
		}
	}
}

func TestDefineResolveBuiltins(t *testing.T) {
	global := NewSymbolTable()
	firstLocal := NewEnclosedSymbolTable(global)
	secondLocal := NewEnclosedSymbolTable(firstLocal)

	expected := []Symbol{
		{Name: "a", Scope: BuiltInScope, Index: 0},
		{Name: "c", Scope: BuiltInScope, Index: 1},
		{Name: "e", Scope: BuiltInScope, Index: 2},
		{Name: "f", Scope: BuiltInScope, Index: 3},
	}

	for i, v := range expected {
		global.DefineBuiltIn(i, v.Name)
	}

	for _, table := range []*SymbolTable{global, firstLocal, secondLocal} {
		for _, sym := range expected {
			result, ok := table.Resolve(sym.Name)
			if !ok {
				t.Errorf("name %s not resolvable", sym.Name)
				continue
			}
			if result != sym {
				t.Errorf("expected %s to resolve to %+v, got=%+v",
					sym.Name, sym, result)
			}
		}
	}
}
