package vm

import (
	"spike-interpreter-go/spike/compiler"
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/lexer"
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
	}

	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			stackTop := runInVM(t, testCase.code)
			assert.Equal(t, testCase.expectedStackTop, stackTop)
		})
	}
}

func runInVM(t *testing.T, input string) object.Object {
	l := lexer.New(strings.NewReader(input))
	p := parser.New(l)
	c := compiler.New()

	program, err := p.ParseProgram()
	assert.NoError(t, err)

	err = c.Compile(program)
	assert.NoError(t, err)

	vm := New(c.Bytecode())

	err = vm.Run()
	assert.NoError(t, err)

	return vm.LastPoppedStackElement()
}
