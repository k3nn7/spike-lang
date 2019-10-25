package eval

import (
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
	"spike-interpreter-go/spike/parser/ast"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Eval_AST(t *testing.T) {
	testCases := []struct {
		input    ast.Node
		expected object.Object
	}{
		{
			input: &ast.Integer{
				Token: lexer.Token{
					Type:    lexer.Integer,
					Literal: "99",
				},
				Value: 99,
			},
			expected: &object.Integer{
				Value: 99,
			},
		},
		{
			input: &ast.Boolean{
				Token: lexer.TrueToken,
				Value: true,
			},
			expected: &object.Boolean{Value: true},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input.String(), func(t *testing.T) {
			result, err := Eval(testCase.input)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}

func Test_Eval_program(t *testing.T) {
	testCases := []struct {
		input    string
		expected object.Object
	}{
		{
			input:    "10",
			expected: &object.Integer{Value: 10},
		},
		{
			input:    "true",
			expected: &object.True,
		},
		{
			input:    "!true",
			expected: &object.False,
		},
		{
			input:    "!false",
			expected: &object.True,
		},
		{
			input:    "!!true",
			expected: &object.True,
		},
		{
			input:    "!!false",
			expected: &object.False,
		},
		{
			input:    "-5",
			expected: &object.Integer{Value: -5},
		},
		{
			input:    "2 + 2",
			expected: &object.Integer{Value: 4},
		},
		{
			input:    "2 - 3",
			expected: &object.Integer{Value: -1},
		},
		{
			input:    "2 * 3",
			expected: &object.Integer{Value: 6},
		},
		{
			input:    "15 / 3",
			expected: &object.Integer{Value: 5},
		},
		{
			input:    "true == false",
			expected: &object.False,
		},
		{
			input:    "2 != 3",
			expected: &object.True,
		},
		{
			input:    "2 < 3",
			expected: &object.True,
		},
		{
			input:    "2 > 3",
			expected: &object.False,
		},
		{
			input:    "2 <= 3",
			expected: &object.True,
		},
		{
			input:    "3 <= 3",
			expected: &object.True,
		},
		{
			input:    "2 >= 3",
			expected: &object.False,
		},
		{
			input:    "3 >= 3",
			expected: &object.True,
		},
		{
			input:    "true || false",
			expected: &object.True,
		},
		{
			input:    "true && false",
			expected: &object.False,
		},
		{
			input:    "(2 > 3) || (true != false)",
			expected: &object.True,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			l := lexer.New(strings.NewReader(testCase.input))
			program, err := parser.New(l).ParseProgram()

			assert.NoError(t, err)

			result, err := Eval(program)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, result)
		})
	}
}
