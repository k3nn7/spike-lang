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
