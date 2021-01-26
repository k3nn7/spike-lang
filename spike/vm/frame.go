package vm

import (
	"spike-interpreter-go/spike/code"
	"spike-interpreter-go/spike/eval/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{
		fn: fn,
		ip: -1,
	}
}

func (frame *Frame) Instructions() code.Instructions {
	return frame.fn.Instructions
}
