package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
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
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
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
		return sy, ok
	}
	return sy, ok
}
