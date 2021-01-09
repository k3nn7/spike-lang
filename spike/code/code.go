package code

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/pkg/errors"
)

type Opcode byte

const (
	Byte              = 1
	OpConstant Opcode = iota
	OpAdd
	OpSub
	OpMul
	OpDiv
	OpPop
	OpTrue
	OpFalse
	OpEqual
	OpNotEqual
	OpGreaterThan
	OpMinus
	OpBang
	OpJumpNotTrue
	OpJump
	OpNull
	OpSetGlobal
	OpGetGlobal
	OpArray
	OpHash
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {
		Name:          "OpConstant",
		OperandWidths: []int{2 * Byte},
	},
	OpAdd: {
		Name:          "OpAdd",
		OperandWidths: []int{},
	},
	OpSub: {
		Name:          "OpSub",
		OperandWidths: []int{},
	},
	OpMul: {
		Name:          "OpMul",
		OperandWidths: []int{},
	},
	OpDiv: {
		Name:          "OpDiv",
		OperandWidths: []int{},
	},
	OpPop: {
		Name:          "OpPop",
		OperandWidths: []int{},
	},
	OpTrue: {
		Name:          "OpTrue",
		OperandWidths: []int{},
	},
	OpFalse: {
		Name:          "OpFalse",
		OperandWidths: []int{},
	},
	OpEqual: {
		Name:          "OpEqual",
		OperandWidths: []int{},
	},
	OpNotEqual: {
		Name:          "OpNotEqual",
		OperandWidths: []int{},
	},
	OpGreaterThan: {
		Name:          "OpGreaterThan",
		OperandWidths: []int{},
	},
	OpMinus: {
		Name:          "OpMinus",
		OperandWidths: []int{},
	},
	OpBang: {
		Name:          "OpBang",
		OperandWidths: []int{},
	},
	OpJump: {
		Name:          "OpJump",
		OperandWidths: []int{2 * Byte},
	},
	OpJumpNotTrue: {
		Name:          "OpJumpNotTrue",
		OperandWidths: []int{2 * Byte},
	},
	OpNull: {
		Name:          "OpNull",
		OperandWidths: []int{},
	},
	OpSetGlobal: {
		Name:          "OpSetGlobal",
		OperandWidths: []int{2},
	},
	OpGetGlobal: {
		Name:          "OpGetGlobal",
		OperandWidths: []int{2},
	},
	OpArray: {
		Name:          "OpArray",
		OperandWidths: []int{2},
	},
	OpHash: {
		Name:          "OpHash",
		OperandWidths: []int{2},
	},
}

type Instructions []byte

func (instructions Instructions) String() string {
	var result bytes.Buffer

	i := 0
	for i < len(instructions) {
		definition, err := Lookup(Opcode(instructions[i]))
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
	case 0:
		return fmt.Sprintf("%s", definition.Name)
	case 1:
		return fmt.Sprintf("%s %d", definition.Name, operands[0])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", definition.Name)
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
