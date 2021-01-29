package code

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Make(t *testing.T) {
	testCases := map[string]struct {
		opcode   Opcode
		operands []int
		expected []byte
	}{
		"OpConstant": {
			opcode:   OpConstant,
			operands: []int{65534},
			expected: []byte{
				byte(OpConstant),
				255,
				254,
			},
		},
		"OpAdd": {
			opcode:   OpAdd,
			operands: []int{},
			expected: []byte{
				byte(OpAdd),
			},
		},
	}

	for testCaseName, testCase := range testCases {
		t.Run(testCaseName, func(t *testing.T) {
			instruction, err := Make(testCase.opcode, testCase.operands...)
			assert.NoError(t, err)
			assert.Equal(t, testCase.expected, instruction)
		})
	}
}

func Test_Instructions_String(t *testing.T) {
	instructions := NewBuilder().
		Make(OpConstant, 2).
		Make(OpConstant, 65535).
		Make(OpAdd).
		Make(OpSub).
		Make(OpMul).
		Make(OpDiv).
		Make(OpPop).
		Make(OpMinus).
		Make(OpBang).
		Make(OpJumpNotTrue, 256).
		Make(OpJump, 128).
		Make(OpNull).
		Make(OpSetGlobal, 256).
		Make(OpGetGlobal, 256).
		Make(OpArray, 256).
		Make(OpHash, 256).
		Make(OpIndex).
		Make(OpCall, 3).
		Make(OpReturnValue).
		Make(OpReturn).
		Make(OpSetLocal, 255).
		Make(OpGetLocal, 255).
		Build()

	expectedOutput := `0000 OpConstant 2
0003 OpConstant 65535
0006 OpAdd
0007 OpSub
0008 OpMul
0009 OpDiv
0010 OpPop
0011 OpMinus
0012 OpBang
0013 OpJumpNotTrue 256
0016 OpJump 128
0019 OpNull
0020 OpSetGlobal 256
0023 OpGetGlobal 256
0026 OpArray 256
0029 OpHash 256
0032 OpIndex
0033 OpCall 3
0035 OpReturnValue
0036 OpReturn
0037 OpSetLocal 255
0039 OpGetLocal 255
`

	assert.Equal(t, expectedOutput, instructions.String())
}

func Test_ReadOperands(t *testing.T) {
	opcode := OpConstant
	expectedOperands := []int{65535}
	expectedOperandBytes := 2

	instruction, err := Make(opcode, expectedOperands...)
	assert.NoError(t, err)

	definition, err := Lookup(opcode)

	operandsRead, operandBytes := ReadOperands(definition, instruction[1:])

	assert.Equal(t, expectedOperandBytes, operandBytes)
	assert.Equal(t, expectedOperands, operandsRead)
}
