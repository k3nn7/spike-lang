package vm

import (
	"spike-interpreter-go/spike/compiler"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/object"
	"spike-interpreter-go/spike/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	testCases := []struct {
		code             string
		expectedStackTop object.Object
	}{
		{
			code:             "1 + 2",
			expectedStackTop: &object.Integer{Value: 3},
		},
		{
			code:             "3 * 5",
			expectedStackTop: &object.Integer{Value: 15},
		},
		{
			code:             "30 / (1 + 2)",
			expectedStackTop: &object.Integer{Value: 10},
		},
		{
			code:             "100 / (5 - 6) * 2",
			expectedStackTop: &object.Integer{Value: -200},
		},
		{
			code:             "true",
			expectedStackTop: True,
		},
		{
			code:             "false",
			expectedStackTop: False,
		},
		{
			code:             "1 < 2",
			expectedStackTop: True,
		},
		{
			code:             "1 > 2",
			expectedStackTop: False,
		},
		{
			code:             "1 == 2",
			expectedStackTop: False,
		},
		{
			code:             "2 == 2",
			expectedStackTop: True,
		},
		{
			code:             "1 != 2",
			expectedStackTop: True,
		},
		{
			code:             "-5",
			expectedStackTop: &object.Integer{Value: -5},
		},
		{
			code:             "!false",
			expectedStackTop: True,
		},
		{
			code:             "if (true) { 10 }",
			expectedStackTop: &object.Integer{Value: 10},
		},
		{
			code:             "if (true) { 10 } else { 20 }",
			expectedStackTop: &object.Integer{Value: 10},
		},
		{
			code:             "if (false) { 10 } else { 20 }",
			expectedStackTop: &object.Integer{Value: 20},
		},
		{
			code:             "if (2 > 1) { 10 } else { 20 }",
			expectedStackTop: &object.Integer{Value: 10},
		},
		{
			code:             "if (2 < 1) { 10 } else { 20 }",
			expectedStackTop: &object.Integer{Value: 20},
		},
		{
			code:             "if (false) { 10 };",
			expectedStackTop: Null,
		},
		{
			code:             "let one = 1; one;",
			expectedStackTop: &object.Integer{Value: 1},
		},
		{
			code:             "let one = 1; let two = 2; one + two;",
			expectedStackTop: &object.Integer{Value: 3},
		},
		{
			code:             "let one = 1; let two = one + one; one + two;",
			expectedStackTop: &object.Integer{Value: 3},
		},
		{
			code:             `"spike"`,
			expectedStackTop: &object.String{Value: "spike"},
		},
		{
			code:             `"spike " + "language"`,
			expectedStackTop: &object.String{Value: "spike language"},
		},
		{
			code:             `[]`,
			expectedStackTop: &object.Array{Elements: []object.Object{}},
		},
		{
			code: `[1, 2, 3]`,
			expectedStackTop: &object.Array{Elements: []object.Object{
				&object.Integer{Value: 1},
				&object.Integer{Value: 2},
				&object.Integer{Value: 3},
			}},
		},
		{
			code: `[1 + 2, 2 + 3]`,
			expectedStackTop: &object.Array{Elements: []object.Object{
				&object.Integer{Value: 3},
				&object.Integer{Value: 5},
			}},
		},
		{
			code:             `{}`,
			expectedStackTop: &object.Hash{Pairs: map[object.HashKey]object.HashPair{}},
		},
		{
			code: `{1:2, 2:3}`,
			expectedStackTop: &object.Hash{Pairs: map[object.HashKey]object.HashPair{
				(&object.Integer{Value: 1}).GetHashKey(): {
					Key:   &object.Integer{Value: 1},
					Value: &object.Integer{Value: 2},
				},
				(&object.Integer{Value: 2}).GetHashKey(): {
					Key:   &object.Integer{Value: 2},
					Value: &object.Integer{Value: 3},
				},
			}},
		},
		{
			code: `{1+2:2-3}`,
			expectedStackTop: &object.Hash{Pairs: map[object.HashKey]object.HashPair{
				(&object.Integer{Value: 3}).GetHashKey(): {
					Key:   &object.Integer{Value: 3},
					Value: &object.Integer{Value: -1},
				},
			}},
		},
		{
			code:             `[1, 2, 3][1]`,
			expectedStackTop: &object.Integer{Value: 2},
		},
		{
			code:             `[1, 2, 3][1 + 1]`,
			expectedStackTop: &object.Integer{Value: 3},
		},
		{
			code:             `[][1]`,
			expectedStackTop: Null,
		},
		{
			code:             `{"name": "kenny", "age": 31}["age"]`,
			expectedStackTop: &object.Integer{Value: 31},
		},
		{
			code:             `{"name": "kenny", "age": 31}["surname"]`,
			expectedStackTop: Null,
		},
		{
			code:             `let f = fn () { 5 + 10 }; f();`,
			expectedStackTop: &object.Integer{Value: 15},
		},
		{
			code:             `let f = fn () { return 5 + 10 }; f();`,
			expectedStackTop: &object.Integer{Value: 15},
		},
		{
			code:             `let f = fn () { }; f();`,
			expectedStackTop: Null,
		},
		{
			code:             `let f = fn() { 1 }; let g = fn() { f }; g()()`,
			expectedStackTop: &object.Integer{Value: 1},
		},
		{
			code:             `let one = fn() { let one = 1; one }; one()`,
			expectedStackTop: &object.Integer{Value: 1},
		},
		{
			code:             `let a = 50; let dec = fn() { let b = 2; return a - b }; dec()`,
			expectedStackTop: &object.Integer{Value: 48},
		},
		{
			code:             `let a = fn() { let aa = 5; aa }; let b = fn () { let bb = 10; bb }; a() + b()`,
			expectedStackTop: &object.Integer{Value: 15},
		},
		{
			code:             `let add = fn() { let a = 10; a }; let a = 5; a + add()`,
			expectedStackTop: &object.Integer{Value: 15},
		},
		{
			code:             `fn() { let a = 7; a } () + 5`,
			expectedStackTop: &object.Integer{Value: 12},
		},
		{
			code:             `let a = fn() { let i = 1; let j = 2; i + j }; a()`,
			expectedStackTop: &object.Integer{Value: 3},
		},
		{
			code: `
			let g = 10;
			let a = fn() { let i = 1; i + g };
			let b = fn() { let i = 2; i + g };
			a() + b();`,
			expectedStackTop: &object.Integer{Value: 23},
		},
		{
			code: `
			let g = 10;
			let a = fn() { let i = 1; i + g };
			let b = fn() { let i = 2; a() + i };
			g + b();`,
			expectedStackTop: &object.Integer{Value: 23},
		},
		{
			code: `
			let f = fn(a) { a };
			f(555);`,
			expectedStackTop: &object.Integer{Value: 555},
		},
		{
			code: `
			let f = fn(a, b) { a + b };
			f(555, 222);`,
			expectedStackTop: &object.Integer{Value: 777},
		},
		{
			code:             `len("abc")`,
			expectedStackTop: &object.Integer{Value: 3},
		},
		{
			code:             `len([1, 2, 3, 4])`,
			expectedStackTop: &object.Integer{Value: 4},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			stackTop, err := runInVM(testCase.code)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedStackTop, stackTop)
		})
	}
}

func runInVM(input string) (object.Object, error) {
	l := lexer.New(strings.NewReader(input))
	p := parser.New(l)
	c := compiler.New()

	program, err := p.ParseProgram()
	if err != nil {
		return nil, err
	}

	err = c.Compile(program)
	if err != nil {
		return nil, err
	}

	vm := New(c.Bytecode())

	err = vm.Run()
	if err != nil {
		return nil, err
	}

	return vm.LastPoppedStackElement(), nil
}
