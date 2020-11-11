package object

import (
	"fmt"
	"hash/fnv"
)

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

func (str *String) GetHashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(str.Value))

	return HashKey{
		Type:  StringType,
		Value: h.Sum64(),
	}
}
