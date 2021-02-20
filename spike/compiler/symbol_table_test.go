package compiler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SymbolTable_Define(t *testing.T) {
	expectedA := Symbol{
		Name:        "a",
		SymbolScope: GlobalScope,
		Index:       0,
	}
	expectedB := Symbol{
		Name:        "b",
		SymbolScope: GlobalScope,
		Index:       1,
	}

	symbolTable := NewSymbolTable()

	a := symbolTable.Define("a")
	assert.Equal(t, expectedA, a)

	b := symbolTable.Define("b")
	assert.Equal(t, expectedB, b)
}

func Test_SymbolTable_ResolveGlobal(t *testing.T) {
	symbolTable := NewSymbolTable()
	symbolTable.Define("a")
	symbolTable.Define("b")

	expectedA := Symbol{
		Name:        "a",
		SymbolScope: GlobalScope,
		Index:       0,
	}
	expectedB := Symbol{
		Name:        "b",
		SymbolScope: GlobalScope,
		Index:       1,
	}

	a, ok := symbolTable.Resolve("a")
	assert.True(t, ok)
	assert.Equal(t, expectedA, a)

	b, ok := symbolTable.Resolve("b")
	assert.True(t, ok)
	assert.Equal(t, expectedB, b)

	_, ok = symbolTable.Resolve("c")
	assert.False(t, ok)
}

func Test_SymbolTable_ResolveLocal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")
	local.Define("d")

	symbol, ok := local.Resolve("a")
	assert.True(t, ok)
	assert.Equal(t, Symbol{
		Name:        "a",
		SymbolScope: GlobalScope,
		Index:       0,
	}, symbol)

	symbol, ok = local.Resolve("b")
	assert.True(t, ok)
	assert.Equal(t, Symbol{
		Name:        "b",
		SymbolScope: GlobalScope,
		Index:       1,
	}, symbol)

	symbol, ok = local.Resolve("c")
	assert.True(t, ok)
	assert.Equal(t, Symbol{
		Name:        "c",
		SymbolScope: LocalScope,
		Index:       0,
	}, symbol)

	symbol, ok = local.Resolve("d")
	assert.True(t, ok)
	assert.Equal(t, Symbol{
		Name:        "d",
		SymbolScope: LocalScope,
		Index:       1,
	}, symbol)
}

func Test_SymbolTable_resolveBuiltin(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.DefineBuiltin(0, "b")

	local := NewEnclosedSymbolTable(global)
	local.Define("c")

	symbol, ok := local.Resolve("b")
	assert.True(t, ok)
	assert.Equal(t, Symbol{
		Name:        "b",
		SymbolScope: BuiltinScope,
		Index:       0,
	}, symbol)
}

func Test_SymbolTable_resolveFreeVars(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")

	local1 := NewEnclosedSymbolTable(global)
	local1.Define("b")

	local2 := NewEnclosedSymbolTable(local1)
	local2.Define("c")

	symbol, ok := local2.Resolve("b")
	assert.True(t, ok)
	assert.Equal(t, Symbol{
		Name:        "b",
		SymbolScope: FreeScope,
		Index:       0,
	}, symbol)

	_, ok = local2.Resolve("d")
	assert.False(t, ok)

	assert.Equal(t, []Symbol{
		{
			Name:        "b",
			SymbolScope: LocalScope,
			Index:       0,
		},
	}, local2.FreeSymbols)
}
