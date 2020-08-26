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
		Build()

	expectedOutput := `0000 OpConstant 2
0003 OpConstant 65535
0006 OpAdd
0007 OpSub
0008 OpMul
0009 OpDiv
0010 OpPop
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
