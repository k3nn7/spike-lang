package object

import "fmt"

type Closure struct {
	Function      *CompiledFunction
	FreeVariables Object
}

func (closure *Closure) Type() ObjectType {
	return ClosureType
}

func (closure *Closure) Inspect() string {
	return fmt.Sprintf("Closure[%p]", closure)
}

func (closure *Closure) Equal(other Object) bool {
	return other == closure
}
