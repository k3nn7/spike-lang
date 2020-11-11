package object

import (
	"fmt"
	"strings"
)

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (hash *Hash) Type() ObjectType {
	return HashType
}

func (hash *Hash) Inspect() string {
	out := strings.Builder{}

	out.WriteString("{")
	inspectedPairs := make([]string, 0, len(hash.Pairs))
	for _, pair := range hash.Pairs {
		inspectedPairs = append(
			inspectedPairs,
			fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()),
		)
	}

	out.WriteString(strings.Join(inspectedPairs, ", "))
	out.WriteString("}")

	return out.String()
}

func (hash *Hash) Equal(other Object) bool {
	otherHash, ok := other.(*Hash)
	if !ok {
		return false
	}

	for key, val := range hash.Pairs {
		val2, ok := otherHash.Pairs[key]
		if !ok {
			return false
		}
		if !val.Value.Equal(val2.Value) {
			return false
		}
	}

	return true
}
