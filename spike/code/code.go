package code

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
)

type Instructions []byte

func (instructions Instructions) String() string {
	var result bytes.Buffer

	i := 0
	for i < len(instructions) {
		definition, err := Lookup(Opcode(instructions[0]))
		if err != nil {
			fmt.Fprintf(&result, "ERROR: %s\n", err)
			continue
		}

		operands, operandBytes := ReadOperands(definition, instructions[i+1:])
		fmt.Fprintf(&result, "%04d %s\n", i, formatInstruction(definition, operands))

		i += 1 + operandBytes
	}

	return result.String()
}

func formatInstruction(definition *Definition, operands []int) string {
	operandCount := len(definition.OperandWidths)

	if len(operands) != operandCount {
		return fmt.Sprintf(
			"ERROR: operand len %d does not match defined %d\n",
			len(operands),
			operandCount,
		)
	}

	switch operandCount {
	case 1:
		return fmt.Sprintf("%s %d", definition.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", definition.Name)
}

type Opcode byte

const (
	Byte              = 1
	OpConstant Opcode = iota
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {Name: "OpConstant", OperandWidths: []int{2 * Byte}},
}

func Lookup(opcode Opcode) (*Definition, error) {
	if definition, ok := definitions[opcode]; ok {
		return definition, nil
	}

	return nil, errors.Errorf("opcode %d undefined", opcode)
}

func Make(opcode Opcode, operands ...int) ([]byte, error) {
	definition, err := Lookup(opcode)
	if err != nil {
		return nil, err
	}

	instructionLength := 1
	for _, operandWidth := range definition.OperandWidths {
		instructionLength += operandWidth
	}

	instruction := make([]byte, instructionLength)
	instruction[0] = byte(opcode)

	offset := 1
	for i, operand := range operands {
		operandWidth := definition.OperandWidths[i]
		switch operandWidth {
		case 2 * Byte:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(operand))
		}
		offset += operandWidth
	}

	return instruction, nil
}

func ReadOperands(definition *Definition, instructions Instructions) ([]int, int) {
	operands := make([]int, len(definition.OperandWidths))
	offset := 0

	for i, width := range definition.OperandWidths {
		switch width {
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(instructions[offset:]))
		}

		offset += width
	}

	return operands, offset
}
