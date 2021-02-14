package compiler

import (
	"spike-interpreter-go/spike/code"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/object"
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
		{
			code: "1 > 2",
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpGreaterThan).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "1 < 2",
			expectedConstants: []object.Object{
				&object.Integer{Value: 2},
				&object.Integer{Value: 1},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpGreaterThan).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "1 == 2",
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpEqual).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "1 != 2",
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpNotEqual).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "-1",
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpMinus).
				Make(code.OpPop).
				Build(),
		},
		{
			code:              "!false",
			expectedConstants: []object.Object{},
			expectedInstructions: code.NewBuilder().
				Make(code.OpFalse).
				Make(code.OpBang).
				Make(code.OpPop).
				Build(),
		},
		{
			code: "if (true) { 10 }; 3333;",
			expectedConstants: []object.Object{
				&object.Integer{Value: 10},
				&object.Integer{Value: 3333},
			},
			expectedInstructions: code.NewBuilder().
				// 0000
				Make(code.OpTrue).
				// 0001
				Make(code.OpJumpNotTrue, 10).
				// 0004
				Make(code.OpConstant, 0).
				// 0007
				Make(code.OpJump, 11).
				// 0010
				Make(code.OpNull).
				// 0011
				Make(code.OpPop).
				// 0012
				Make(code.OpConstant, 1).
				// 0015
				Make(code.OpPop).
				Build(),
		},
		{
			code: "if (true) { 10 } else { 20 }; 3333",
			expectedConstants: []object.Object{
				&object.Integer{Value: 10},
				&object.Integer{Value: 20},
				&object.Integer{Value: 3333},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpTrue).
				Make(code.OpJumpNotTrue, 10).
				Make(code.OpConstant, 0).
				Make(code.OpJump, 13).
				Make(code.OpConstant, 1).
				Make(code.OpPop).
				Make(code.OpConstant, 2).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `let one = 1; let two = 2;`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpConstant, 1).
				Make(code.OpSetGlobal, 1).
				Build(),
		},
		{
			code: `let one = 1; one;`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpGetGlobal, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `let one = 1; let two = one; two;`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpGetGlobal, 0).
				Make(code.OpSetGlobal, 1).
				Make(code.OpGetGlobal, 1).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `"spike"`,
			expectedConstants: []object.Object{
				&object.String{Value: "spike"},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `"spike " + "language"`,
			expectedConstants: []object.Object{
				&object.String{Value: "spike "},
				&object.String{Value: "language"},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpAdd).
				Make(code.OpPop).
				Build(),
		},
		{
			code:              `[]`,
			expectedConstants: []object.Object{},
			expectedInstructions: code.NewBuilder().
				Make(code.OpArray, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `[1, 2, 3]`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpConstant, 2).
				Make(code.OpArray, 3).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `[1 + 2, 2 - 3]`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpAdd).
				Make(code.OpConstant, 2).
				Make(code.OpConstant, 3).
				Make(code.OpSub).
				Make(code.OpArray, 2).
				Make(code.OpPop).
				Build(),
		},
		{
			code:              `{}`,
			expectedConstants: []object.Object{},
			expectedInstructions: code.NewBuilder().
				Make(code.OpHash, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `{1: 2, 3: 4}`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
				&object.Integer{Value: 4},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpConstant, 2).
				Make(code.OpConstant, 3).
				Make(code.OpHash, 4).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `{1 + 2: 2 - 3}`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpAdd).
				Make(code.OpConstant, 2).
				Make(code.OpConstant, 3).
				Make(code.OpSub).
				Make(code.OpHash, 2).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `[1, 2][0 + 1]`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 0},
				&object.Integer{Value: 1},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpArray, 2).
				Make(code.OpConstant, 2).
				Make(code.OpConstant, 3).
				Make(code.OpAdd).
				Make(code.OpIndex).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `{1: 2}[0 + 1]`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 0},
				&object.Integer{Value: 1},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpConstant, 1).
				Make(code.OpHash, 2).
				Make(code.OpConstant, 2).
				Make(code.OpConstant, 3).
				Make(code.OpAdd).
				Make(code.OpIndex).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `fn () { return 5 + 10 }`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 5},
				&object.Integer{Value: 10},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpConstant, 0).
						Make(code.OpConstant, 1).
						Make(code.OpAdd).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 2, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `fn () { 5 + 10 }`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 5},
				&object.Integer{Value: 10},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpConstant, 0).
						Make(code.OpConstant, 1).
						Make(code.OpAdd).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 2, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `fn () { 1; 2 }`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpConstant, 0).
						Make(code.OpPop).
						Make(code.OpConstant, 1).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 2, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `fn () { }`,
			expectedConstants: []object.Object{
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpReturn).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 0, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `fn() { 24 } ()`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 24},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpConstant, 0).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 1, 0).
				Make(code.OpCall, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `let f = fn() { 24 }; f()`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 24},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpConstant, 0).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 1, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpGetGlobal, 0).
				Make(code.OpCall, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `let a = 55; fn() { a }`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 55},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpGetGlobal, 0).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     0,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpConstant, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpClosure, 1, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `fn() { let a = 55; return a; }`,
			expectedConstants: []object.Object{
				&object.Integer{Value: 55},
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpConstant, 0).
						Make(code.OpSetLocal, 0).
						Make(code.OpGetLocal, 0).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     1,
					ParametersCount: 0,
				},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 1, 0).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `let f = fn(a) { a }; f(24);`,
			expectedConstants: []object.Object{
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpGetLocal, 0).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     1,
					ParametersCount: 1,
				},
				&object.Integer{Value: 24},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 0, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpGetGlobal, 0).
				Make(code.OpConstant, 1).
				Make(code.OpCall, 1).
				Make(code.OpPop).
				Build(),
		},
		{
			code: `let f = fn(a, b, c) { a; b; c; }; f(2, 4, 6);`,
			expectedConstants: []object.Object{
				&object.CompiledFunction{
					Instructions: code.NewBuilder().
						Make(code.OpGetLocal, 0).
						Make(code.OpPop).
						Make(code.OpGetLocal, 1).
						Make(code.OpPop).
						Make(code.OpGetLocal, 2).
						Make(code.OpReturnValue).
						Build(),
					LocalsCount:     3,
					ParametersCount: 3,
				},
				&object.Integer{Value: 2},
				&object.Integer{Value: 4},
				&object.Integer{Value: 6},
			},
			expectedInstructions: code.NewBuilder().
				Make(code.OpClosure, 0, 0).
				Make(code.OpSetGlobal, 0).
				Make(code.OpGetGlobal, 0).
				Make(code.OpConstant, 1).
				Make(code.OpConstant, 2).
				Make(code.OpConstant, 3).
				Make(code.OpCall, 3).
				Make(code.OpPop).
				Build(),
		},
		{
			code:              `len([])`,
			expectedConstants: []object.Object{},
			expectedInstructions: code.NewBuilder().
				Make(code.OpGetBuiltin, 0).
				Make(code.OpArray, 0).
				Make(code.OpCall, 1).
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
