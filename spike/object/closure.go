package object

import "fmt"

type Closue struct {
	Function      *CompiledFunction
	FreeVariables Object
}

func (closure *Closue) Type() ObjectType {
	return ClosureType
}

func (closure *Closue) Inspect() string {
	return fmt.Sprintf("Closure[%p]", closure)
}

func (closure *Closue) Equal(other Object) bool {
	return other == closure
}
