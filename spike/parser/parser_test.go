package parser

import (
	"io"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_parseValidCode(t *testing.T) {
	testCases := map[string]struct {
		code            io.Reader
		expectedProgram ast.Program
	}{
		"let statement": {
			code: strings.NewReader(`let variable = 10;`),
			expectedProgram: ast.Program{Statements: []ast.Statement{
				&ast.LetStatement{
					Token: lexer.Token{Type: lexer.Let, Literal: "let"},
					Name: &ast.Identifier{
						Token: lexer.Token{Type: lexer.Identifier, Literal: "variable"},
						Value: "variable",
					},
				},
			}},
		},
		"return statement": {
			code: strings.NewReader(`return 2 + 2;`),
			expectedProgram: ast.Program{Statements: []ast.Statement{
				&ast.ReturnStatement{
					Token: lexer.Token{Type: lexer.Return, Literal: "return"},
				},
			}},
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			parser := New(lexer.New(testCase.code))

			program, err := parser.ParseProgram()

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedProgram, program)
		})
	}
}

func Test_Parser_parsingError(t *testing.T) {
	testCases := map[string]struct {
		code          io.Reader
		expectedError string
	}{
		"missing assignment in let statement": {
			code:          strings.NewReader("let variable 10;"),
			expectedError: "expected assign operator, got integer",
		},
		"missing identifier in let statement": {
			code:          strings.NewReader("let = 10;"),
			expectedError: "expected identifier, got assign",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			parser := New(lexer.New(testCase.code))

			_, err := parser.ParseProgram()

			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}
