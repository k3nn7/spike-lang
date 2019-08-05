package repl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	input := strings.NewReader(" let x = 10\n")
	expectedOutput := ">> {Type:let Literal:let}\n{Type:identifier Literal:x}\n{Type:assign Literal:=}\n{Type:integer Literal:10}\n>> "
	output := &strings.Builder{}

	Start(input, output)

	assert.Equal(t, expectedOutput, output.String())
}
