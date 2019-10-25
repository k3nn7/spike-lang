package object

import "fmt"

type Integer struct {
	Value int64
}

func (integer *Integer) Type() ObjectType {
	return IntegerType
}

func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

func (integer *Integer) Equal(other Object) (bool, error) {
	otherInteger, ok := other.(*Integer)
	if !ok {
		return false, NotComparableError
	}

	return integer.Value == otherInteger.Value, nil
}

func (integer *Integer) Compare(other Comparable) (Ordering, error) {
	otherInteger := other.(*Integer)

	if integer.Value > otherInteger.Value {
		return GT, nil
	} else if integer.Value < otherInteger.Value {
		return LT, nil
	}

	return EQ, nil
}
