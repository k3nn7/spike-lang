package parser

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Expressions(t *testing.T) {
	input := "foobar;"
	expectedProgram := "foobar\n"

	program, err := New(lexer.New(strings.NewReader(input))).ParseProgram()

	assert.NoError(t, err)
	assert.Equal(t, expectedProgram, program.String())
}
