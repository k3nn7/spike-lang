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
	testCases := []struct {
		code                 string
		expectedConstants    []object.Object
		expectedInstructions code.Instructions
	}{
		{
			code: "1 + 2",
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpAdd).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "3 - 4",
			expectedConstants: []object.Object{
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpSub).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "19 * 8",
			expectedConstants: []object.Object{
				&object.Integer{Value: 19},
				&object.Integer{Value: 8},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpMul).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "35 / 7",
			expectedConstants: []object.Object{
				&object.Integer{Value: 35},
				&object.Integer{Value: 7},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpDiv).
				Make(code.OpPop).
				Build(),
		},
		{
			code:              "true",
			expectedConstants: []object.Object{},
			expectedInstructions: code.NewBuilder().
				Make(code.OpTrue).
				Make(code.OpPop).
				Build(),
		},
		{
			code:              "false",
			expectedConstants: []object.Object{},
			expectedInstructions: code.NewBuilder().
				Make(code.OpFalse).
				Make(code.OpPop).
				Build(),
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			bytecode := compileCode(t, testCase.code)
			assert.Equal(t, testCase.expectedConstants, bytecode.Constants)
			assert.Equal(t, testCase.expectedInstructions.String(), bytecode.Instructions.String())
		})
	}
}

func compileCode(t *testing.T, input string) *Bytecode {
	l := lexer.New(strings.NewReader(input))
	p := parser.New(l)
	compiler := New()

	program, err := p.ParseProgram()
	assert.NoError(t, err)

	err = compiler.Compile(program)
	assert.NoError(t, err)

	return compiler.Bytecode()
}