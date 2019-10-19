package repl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	input := strings.NewReader("10\n")
	expectedOutput := ">> 10\n>> "
	output := &strings.Builder{}

	Start(input, output)

	assert.Equal(t, expectedOutput, output.String())
}
