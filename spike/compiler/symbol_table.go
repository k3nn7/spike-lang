package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name        string
	SymbolScope SymbolScope
	Index       int
}

type SymbolTable struct {
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func (symbolTable *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:        name,
		SymbolScope: GlobalScope,
		Index:       symbolTable.numDefinitions,
	}
	symbolTable.store[name] = symbol
	symbolTable.numDefinitions++

	return symbol
}

func (symbolTable *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := symbolTable.store[name]
	return symbol, ok
}
