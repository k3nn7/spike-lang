package object

var NullObject = Null{}

type Null struct{}

func (null *Null) Type() ObjectType {
	return NullType
}

func (null *Null) Inspect() string {
	return "null"
}

func (null *Null) Equal(other Object) (bool, error) {
	_, ok := other.(*Null)
	if !ok {
		return false, NotComparableError
	}

	return true, nil
}
