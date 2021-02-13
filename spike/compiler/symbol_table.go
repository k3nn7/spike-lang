package compiler

type SymbolScope string

const (
	BuiltinScope SymbolScope = "BUILTIN"
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
)

type Symbol struct {
	Name        string
	SymbolScope SymbolScope
	Index       int
}

type SymbolTable struct {
	Outer          *SymbolTable
	store          map[string]Symbol
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	return &SymbolTable{
		Outer:          outer,
		store:          make(map[string]Symbol),
		numDefinitions: 0,
	}
}

func (symbolTable *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: symbolTable.numDefinitions}
	if symbolTable.Outer == nil {
		symbol.SymbolScope = GlobalScope
	} else {
		symbol.SymbolScope = LocalScope
	}
	symbolTable.store[name] = symbol
	symbolTable.numDefinitions++

	return symbol
}

func (symbolTable *SymbolTable) DefineBuiltin(index int, name string) {
	symbol := Symbol{Name: name, Index: index, SymbolScope: BuiltinScope}
	symbolTable.store[name] = symbol
}

func (symbolTable *SymbolTable) Resolve(name string) (Symbol, bool) {
	symbol, ok := symbolTable.store[name]

	if !ok && symbolTable.Outer != nil {
		return symbolTable.Outer.Resolve(name)
	}

	return symbol, ok
}
