package parser

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Parser_code_sample(t *testing.T) {
	input := strings.NewReader(`let variable = 10;`)
	expectedProgram := Program{Statements: []Statement{
		&LetStatement{
			Token: lexer.Token{Type: lexer.Let, Literal: "let"},
			Name: &Identifier{
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
