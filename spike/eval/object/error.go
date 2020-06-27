package object

import "fmt"

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ErrorType
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("Error: %s", e.Message)
}

func (e *Error) Equal(other Object) (bool, error) {
	return false, NotComparableError
}
