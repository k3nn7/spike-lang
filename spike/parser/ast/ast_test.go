package ast

import (
	"spike-interpreter-go/spike/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Node2String(t *testing.T) {
	testCases := []struct {
		ast      Node
		expected string
	}{
		{
			ast: &PrefixExpression{
				Token: lexer.Token{
					Type:    lexer.Bang,
					Literal: "!",
				},
				Operator: "!",
				Right: &Identifier{
					Token: lexer.Token{lexer.Identifier, "bool"},
					Value: "bool",
				},
			},
			expected: "(!bool)",
		},
		{
			ast: &Program{Statements: []Statement{
				&LetStatement{
					Token: lexer.Token{lexer.Let, "let"},
					Name: &Identifier{
						Token: lexer.Token{lexer.Identifier, "var"},
						Value: "var",
					},
					Value: &Identifier{
						Token: lexer.Token{lexer.Identifier, "var2"},
						Value: "var2",
					},
				},
			}},
			expected: "let var = var2\n",
		},
		{
			ast: &InfixExpression{
				Token: lexer.Token{Type: lexer.Plus, Literal: "+"},
				Left: &Integer{
					Token: lexer.Token{
						Type:    lexer.Integer,
						Literal: "55",
					},
					Value: 55,
				},
				Operator: "+",
				Right: &Integer{
					Token: lexer.Token{
						Type:    lexer.Integer,
						Literal: "99",
					},
					Value: 99,
				},
			},
			expected: "(55 + 99)",
		},
	}

	for _, testCase := range testCases {
		assert.Equal(t, testCase.expected, testCase.ast.String())
	}
}
