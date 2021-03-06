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
	OpIndex
	OpCall
	OpReturnValue
	OpReturn
	OpSetLocal
	OpGetLocal
	OpGetBuiltin
	OpClosure
	OpGetFreeVar
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
		OperandWidths: []int{2 * Byte},
	},
	OpGetGlobal: {
		Name:          "OpGetGlobal",
		OperandWidths: []int{2 * Byte},
	},
	OpArray: {
		Name:          "OpArray",
		OperandWidths: []int{2 * Byte},
	},
	OpHash: {
		Name:          "OpHash",
		OperandWidths: []int{2 * Byte},
	},
	OpIndex: {
		Name:          "OpIndex",
		OperandWidths: []int{},
	},
	OpCall: {
		Name:          "OpCall",
		OperandWidths: []int{1 * Byte},
	},
	OpReturnValue: {
		Name:          "OpReturnValue",
		OperandWidths: []int{},
	},
	OpReturn: {
		Name:          "OpReturn",
		OperandWidths: []int{},
	},
	OpSetLocal: {
		Name:          "OpSetLocal",
		OperandWidths: []int{1 * Byte},
	},
	OpGetLocal: {
		Name:          "OpGetLocal",
		OperandWidths: []int{1 * Byte},
	},
	OpGetBuiltin: {
		Name:          "OpGetBuiltin",
		OperandWidths: []int{1 * Byte},
	},
	OpClosure: {
		Name:          "OpClosure",
		OperandWidths: []int{2 * Byte, 1 * Byte},
	},
	OpGetFreeVar: {
		Name:          "OpGetFreeVar",
		OperandWidths: []int{1 * Byte},
	},
}

type Instructions []byte

func (instructions Instructions) String() string {
	var result bytes.Buffer

	i := 0
	for i < len(instructions) {
		definition, err := Lookup(Opcode(instructions[i]))
		if err != nil {
			_, err = fmt.Fprintf(&result, "ERROR: %s\n", err)
			if err != nil {
				panic(err)
			}
			continue
		}

		operands, operandBytes := ReadOperands(definition, instructions[i+1:])
		_, err = fmt.Fprintf(&result, "%04d %s\n", i, formatInstruction(definition, operands))
		if err != nil {
			panic(err)
		}

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
	case 2:
		return fmt.Sprintf("%s %d %d", definition.Name, operands[0], operands[1])
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
		case 1 * Byte:
			instruction[offset] = byte(operand)
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
		case 1 * Byte:
			operands[i] = int(instructions[offset])
		case 2 * Byte:
			operands[i] = int(binary.BigEndian.Uint16(instructions[offset:]))
		}

		offset += width
	}

	return operands, offset
}
