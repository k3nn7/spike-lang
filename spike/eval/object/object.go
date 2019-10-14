package object

type ObjectType string

const (
	IntegerType ObjectType = "integer"
	BooleanType ObjectType = "boolean"
	NullType    ObjectType = "null"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
