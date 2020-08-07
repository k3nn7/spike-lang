package object

import "fmt"

type BuiltinFunction struct {
	Name     string
	Function func(args ...Object) (Object, error)
}

func (builtin *BuiltinFunction) Type() ObjectType {
	return BuiltinFunctionType
}

func (builtin *BuiltinFunction) Inspect() string {
	return fmt.Sprintf("builtin(%s)", builtin.Name)
}

func (builtin *BuiltinFunction) Equal(other Object) (bool, error) {
	panic("implement me")
}
