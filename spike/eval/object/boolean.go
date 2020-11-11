package object

import "fmt"

var (
	True  = Boolean{Value: true}
	False = Boolean{Value: false}
)

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() ObjectType {
	return BooleanType
}

func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}

func (boolean *Boolean) Equal(other Object) bool {
	otherBoolean, ok := other.(*Boolean)
	if !ok {
		return false
	}

	return boolean.Value == otherBoolean.Value
}

func (boolean *Boolean) GetHashKey() HashKey {
	if boolean.Value {
		return HashKey{Type: BooleanType, Value: 1}
	}

	return HashKey{Type: BooleanType, Value: 0}
}
