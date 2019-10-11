package repl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	input := strings.NewReader("let x = 10 + 5 * 6\n")
	expectedOutput := ">> let x = (10 + (5 * 6))\n>> "
	output := &strings.Builder{}

	Start(input, output)

	assert.Equal(t, expectedOutput, output.String())
}
