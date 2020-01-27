package object

type ObjectType string

const (
	IntegerType ObjectType = "integer"
	BooleanType ObjectType = "boolean"
	NullType    ObjectType = "null"
	ReturnType  ObjectType = "return"
)

type Ordering int8

const (
	EQ Ordering = 0
	LT Ordering = -1
	GT Ordering = 1
)

type Object interface {
	Type() ObjectType
	Inspect() string
	Equal(other Object) (bool, error)
}

type Comparable interface {
	Compare(other Comparable) (Ordering, error)
}
