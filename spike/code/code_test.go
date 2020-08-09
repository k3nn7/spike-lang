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
	instruction1, err := Make(OpConstant, 1)
	assert.NoError(t, err)
	instruction2, err := Make(OpConstant, 2)
	assert.NoError(t, err)
	instruction3, err := Make(OpConstant, 65535)
	assert.NoError(t, err)
	instructions := []Instructions{instruction1, instruction2, instruction3}

	expectedOutput := `0000 OpConstant 1
0003 OpConstant 2
0006 OpConstant 65535
`

	concattedInstructions := Instructions{}
	for _, instruction := range instructions {
		concattedInstructions = append(concattedInstructions, instruction...)
	}

	assert.Equal(t, expectedOutput, concattedInstructions.String())
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
