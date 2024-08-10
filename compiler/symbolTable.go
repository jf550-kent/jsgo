package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store             map[string]Symbol
	numberDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{store: make(map[string]Symbol)}
}

func (st *SymbolTable) Define(s string) Symbol {
	symbol := Symbol{Name: s, Index: st.numberDefinitions, Scope: GlobalScope}
	st.numberDefinitions++
	st.store[s] = symbol
	return symbol
}

func (st *SymbolTable) Resolve(s string) (Symbol, bool) {
	sy, ok := st.store[s]
	return sy, ok
}
