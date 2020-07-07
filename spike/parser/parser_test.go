package parser

import (
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_parseValidCode(t *testing.T) {
	testCases := []struct {
		code            string
		expectedProgram *ast.Program
	}{
		{
			code: `let variable = 10;`,
			expectedProgram: &ast.Program{Statements: []ast.Statement{
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
		{
			code: `return 2 + 2;`,
			expectedProgram: &ast.Program{Statements: []ast.Statement{
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

	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			program, err := New(lexer.New(strings.NewReader(testCase.code))).ParseProgram()

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedProgram, program)
		})
	}
}

func Test_Parser_ParseProgram(t *testing.T) {
	testCases := []struct {
		code        string
		expectedAst string
	}{
		{
			code:        "let variable = 2 + 2 * 2;",
			expectedAst: "let variable = (2 + (2 * 2))\n",
		},
		{
			code:        "return 2 + variable * 2;",
			expectedAst: "return (2 + (variable * 2))\n",
		},
		{
			code:        "if (true == false) { let a = 10; };",
			expectedAst: "if (true == false) {\n  let a = 10;\n}\n",
		},
		{
			code:        "if (true == false) { let a = 10; } else { let a = 20; };",
			expectedAst: "if (true == false) {\n  let a = 10;\n} else {\n  let a = 20;\n}\n",
		},
		{
			code:        "fn (x, y) { return x + y; }",
			expectedAst: "fn (x, y) {\n  return (x + y);\n}\n",
		},
		{
			code:        "fn (x, y) { let x = 2; return x; }",
			expectedAst: "fn (x, y) {\n  let x = 2;\n  return x;\n}\n",
		},
		{
			code:        "fn (x) { x; }",
			expectedAst: "fn (x) {\n  x;\n}\n",
		},
		{
			code:        "add(5);",
			expectedAst: "add(5);\n",
		},
		{
			code:        "add(x, 2 + 5);",
			expectedAst: "add(x, (2 + 5));\n",
		},
		{
			code:        "fn (x) { x; }(5)",
			expectedAst: "fn (x) {\n  x;\n}(5);\n",
		},
		{
			code:        "\"hello world\"",
			expectedAst: "\"hello world\"\n",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			program, err := New(lexer.New(strings.NewReader(testCase.code))).ParseProgram()

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedAst, program.String())
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
