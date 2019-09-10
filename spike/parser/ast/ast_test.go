package ast

import (
	"spike-interpreter-go/spike/lexer"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_String(t *testing.T) {
	program := &Program{Statements: []Statement{
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
	}}
	expectedProgramString := "let var = var2\n"

	assert.Equal(t, expectedProgramString, program.String())
}
