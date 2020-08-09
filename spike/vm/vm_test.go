package vm

import (
	"spike-interpreter-go/spike/compiler"
	"spike-interpreter-go/spike/eval/object"
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Run(t *testing.T) {
	input := "1 + 2"
	expectedStackTop := &object.Integer{Value: 2}

	l := lexer.New(strings.NewReader(input))
	p := parser.New(l)

	program, err := p.ParseProgram()
	assert.NoError(t, err)

	c := compiler.New()
	err = c.Compile(program)
	assert.NoError(t, err)

	vm := New(c.Bytecode())

	err = vm.Run()
	assert.NoError(t, err)

	stackTop := vm.StackTop()
	assert.Equal(t, expectedStackTop, stackTop)
}
