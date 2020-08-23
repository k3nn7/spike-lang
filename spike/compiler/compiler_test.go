package compiler

import (
	"spike-interpreter-go/spike/code"
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Compiler(t *testing.T) {
	input := "1 + 2"
	expectedConstants := []object.Object{
		&object.Integer{Value: 1},
		&object.Integer{Value: 2},
	}

	instruction1, err := code.Make(code.OpConstant, 0)
	assert.NoError(t, err)
	instruction2, err := code.Make(code.OpConstant, 1)
	assert.NoError(t, err)
	instruction3, err := code.Make(code.OpAdd)
	assert.NoError(t, err)
	instruction4, err := code.Make(code.OpPop)
	assert.NoError(t, err)
	expectedInstructions := concatInstructions(
		instruction1,
		instruction2,
		instruction3,
		instruction4,
	)

	l := lexer.New(strings.NewReader(input))
	p := parser.New(l)
	compiler := New()
	program, err := p.ParseProgram()

	assert.NoError(t, err)

	err = compiler.Compile(program)
	result := compiler.Bytecode()
	assert.NoError(t, err)
	assert.Equal(t, expectedConstants, result.Constants)
	assert.Equal(t, expectedInstructions.String(), result.Instructions.String())
}

func concatInstructions(instructions ...code.Instructions) code.Instructions {
	result := code.Instructions{}
	for _, instruction := range instructions {
		result = append(result, instruction...)
	}

	return result
}
