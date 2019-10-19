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
