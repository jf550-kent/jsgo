package compiler

type SymbolScope string

const (
	GlobalScope   SymbolScope = "GLOBAL"
	LocalScope    SymbolScope = "LOCAL"
	BuiltInScope  SymbolScope = "BUILT_IN"
	FreeScope     SymbolScope = "FREE"
	FunctionScope SymbolScope = "FUNCTION"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer             *SymbolTable
	store             map[string]Symbol
	numberDefinitions int
	FreeSymbols       []Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol), FreeSymbols: []Symbol{}}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (st *SymbolTable) Define(s string) Symbol {
	symbol := Symbol{Name: s, Index: st.numberDefinitions, Scope: GlobalScope}
	if st.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	st.numberDefinitions++
	st.store[s] = symbol
	return symbol
}

func (st *SymbolTable) Resolve(s string) (Symbol, bool) {
	sy, ok := st.store[s]
	if !ok && st.Outer != nil {
		sy, ok = st.Outer.Resolve(s)
		if !ok {
			return sy, ok
		}
		if sy.Scope == GlobalScope || sy.Scope == BuiltInScope {
			return sy, ok
		}
		free := st.defineFree(sy)
		return free, ok
	}
	return sy, ok
}

func (s *SymbolTable) DefineBuiltIn(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltInScope}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) defineFree(sym Symbol) Symbol {
	s.FreeSymbols = append(s.FreeSymbols, sym)

	symbol := Symbol{Name: sym.Name, Index: len(s.FreeSymbols) - 1}
	symbol.Scope = FreeScope

	s.store[sym.Name] = symbol
	return symbol
}

func (s *SymbolTable) DefineFunctionName(name string) Symbol {
	symbol := Symbol{Name: name, Index: 0, Scope: FunctionScope}
	s.store[name] = symbol
	return symbol
}
