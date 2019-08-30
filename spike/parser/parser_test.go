package parser

import (
	"io"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_code_sample(t *testing.T) {
	input := strings.NewReader(`let variable = 10;`)
	expectedProgram := ast.Program{Statements: []ast.Statement{
		&ast.LetStatement{
			Token: lexer.Token{Type: lexer.Let, Literal: "let"},
			Name: &ast.Identifier{
				Token: lexer.Token{Type: lexer.Identifier, Literal: "variable"},
				Value: "variable",
			},
		},
	}}

	parser := New(lexer.New(input))

	program, err := parser.ParseProgram()

	assert.NoError(t, err)
	assert.Equal(t, expectedProgram, program)
}

func Test_Parser_parsingError(t *testing.T) {
	testCases := map[string]struct {
		input         io.Reader
		expectedError string
	}{
		"missing assignment in let statement": {
			input:         strings.NewReader("let variable 10;"),
			expectedError: "expected assign operator, got integer",
		},
		"missing identifier in let statement": {
			input:         strings.NewReader("let = 10;"),
			expectedError: "expected identifier, got assign",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			parser := New(lexer.New(testCase.input))

			_, err := parser.ParseProgram()

			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}
