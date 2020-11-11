package object

import "fmt"

type String struct {
	Value string
}

func (str *String) Type() ObjectType {
	return StringType
}

func (str *String) Inspect() string {
	return fmt.Sprintf("\"%s\"", str.Value)
}

func (str *String) Equal(other Object) bool {
	otherString, ok := other.(*String)
	if !ok {
		return false
	}

	return str.Value == otherString.Value
}

	return str.Value == otherString.Value, nil
}
