package eval

import (
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Eval_withErrors(t *testing.T) {
	testCases := []struct {
		input         string
		expectedError string
	}{
		{
			input:         "2 + true",
			expectedError: "type mismatch: integer + boolean",
		},
		{
			input:         "if (true) { 2 + true; let a = true; 2; }",
			expectedError: "type mismatch: integer + boolean",
		},
		{
			input:         "!10+true",
			expectedError: "type mismatch: !integer",
		},
		{
			input:         "!(10+true)",
			expectedError: "type mismatch: integer + boolean",
		},
		{
			input:         "-true",
			expectedError: "type mismatch: -boolean",
		},
		{
			input:         "2 - true",
			expectedError: "type mismatch: integer - boolean",
		},
		{
			input:         "2 * true",
			expectedError: "type mismatch: integer * boolean",
		},
		{
			input:         "2 / true",
			expectedError: "type mismatch: integer / boolean",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			l := lexer.New(strings.NewReader(testCase.input))
			program, err := parser.New(l).ParseProgram()

			assert.NoError(t, err)

			_, err = Eval(program)
			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}
