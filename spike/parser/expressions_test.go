package parser

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Expressions(t *testing.T) {
	testCases := map[string]struct {
		input           string
		expectedProgram string
	}{
		"single identifier": {
			input:           "foobar;",
			expectedProgram: "foobar\n",
		},
		"let statement with two identifiers": {
			input:           "let var1 = var2;",
			expectedProgram: "let var1 = var2\n",
		},
		"let statement with integer literal": {
			input:           "let var = 125;",
			expectedProgram: "let var = 125\n",
		},
		"negate identifier": {
			input:           "! boolVariable;",
			expectedProgram: "(!boolVariable)\n",
		},
		"negate integer": {
			input:           "! 0;",
			expectedProgram: "(!0)\n",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			program, err := New(lexer.New(strings.NewReader(testCase.input))).ParseProgram()

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedProgram, program.String())
		})
	}
}
