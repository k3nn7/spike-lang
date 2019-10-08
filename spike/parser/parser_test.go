package parser

import (
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_parseValidCode(t *testing.T) {
	testCases := map[string]struct {
		code            string
		expectedProgram ast.Program
	}{
		"let statement": {
			code: `let variable = 10;`,
			expectedProgram: ast.Program{Statements: []ast.Statement{
				&ast.LetStatement{
					Token: lexer.Token{Type: lexer.Let, Literal: "let"},
					Name: &ast.Identifier{
						Token: lexer.Token{Type: lexer.Identifier, Literal: "variable"},
						Value: "variable",
					},
					Value: &ast.Integer{
						Token: lexer.Token{Type: lexer.Integer, Literal: "10"},
						Value: 10,
					},
				},
			}},
		},
		"return statement": {
			code: `return 2 + 2;`,
			expectedProgram: ast.Program{Statements: []ast.Statement{
				&ast.ReturnStatement{
					Token: lexer.Token{Type: lexer.Return, Literal: "return"},
					Result: &ast.InfixExpression{
						Token: lexer.Token{
							Type:    lexer.Plus,
							Literal: "+",
						},
						Left: &ast.Integer{
							Token: lexer.Token{
								Type:    lexer.Integer,
								Literal: "2",
							},
							Value: 2,
						},
						Operator: "+",
						Right: &ast.Integer{
							Token: lexer.Token{
								Type:    lexer.Integer,
								Literal: "2",
							},
							Value: 2,
						},
					},
				},
			}},
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			program, err := New(lexer.New(strings.NewReader(testCase.code))).ParseProgram()

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedProgram, program)
		})
	}
}

func Test_Parser_parsingError(t *testing.T) {
	testCases := map[string]struct {
		code          string
		expectedError string
	}{
		"missing assignment in let statement": {
			code:          "let variable 10;",
			expectedError: "expected assign operator, got integer",
		},
		"missing identifier in let statement": {
			code:          "let = 10;",
			expectedError: "expected identifier, got assign",
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			parser := New(lexer.New(strings.NewReader(testCase.code)))

			_, err := parser.ParseProgram()

			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}
