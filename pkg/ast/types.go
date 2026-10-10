package ast

import "sync"

type SymbolType string

const (
	SymbolFunction SymbolType = "FUNCTION"
	SymbolStruct   SymbolType = "STRUCT"
	SymbolClass    SymbolType = "CLASS"
)

type Symbol struct {
	Name       string
	Type       SymbolType
	SourceFile string
	LineNumber int
}

type SymbolTable struct {
	mu sync.RWMutex
	Symbols map[string]Symbol
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		Symbols: make(map[string]Symbol),
	}
}


func (st *SymbolTable) AddSymbol(sym Symbol) {
	st.mu.Lock()
	defer st.mu.Unlock()
	key := sym.Name + "_" + string(sym.Type)
	st.Symbols[key] = sym
}


func (st* SymbolTable) HasSymbol(name string) bool{
	st.mu.RLock()
	defer st.mu.RUnlock()
	_, exists := st.Symbols[name] 
	return exists
}