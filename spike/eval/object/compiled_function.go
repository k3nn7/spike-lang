package object

import (
	"fmt"
	"spike-interpreter-go/spike/code"
)

type CompiledFunction struct {
	Instructions    code.Instructions
	LocalsCount     int
	ParametersCount int
}

func (function *CompiledFunction) Type() ObjectType {
	return CompiledFunctionType
}

func (function *CompiledFunction) Inspect() string {
	return fmt.Sprintf("CompiledFunction[%p]", function)
}

func (function *CompiledFunction) Equal(other Object) bool {
	return other == function
}
