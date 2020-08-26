package code

type InstructionBuilder struct {
	instructions Instructions
	err          error
}

func NewBuilder() *InstructionBuilder {
	return &InstructionBuilder{}
}

func (builder *InstructionBuilder) Make(opcode Opcode, operands ...int) *InstructionBuilder {
	if builder.err != nil {
		return builder
	}

	instruction, err := Make(opcode, operands...)
	builder.instructions = append(builder.instructions, instruction...)
	builder.err = err

	return builder
}

func (builder *InstructionBuilder) Build() Instructions {
	if builder.err != nil {
		panic(builder.err)
	}
	return builder.instructions
}
