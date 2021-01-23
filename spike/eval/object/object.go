package object

type ObjectType string

const (
	IntegerType          ObjectType = "integer"
	StringType           ObjectType = "string"
	BooleanType          ObjectType = "boolean"
	NullType             ObjectType = "null"
	ReturnType           ObjectType = "return"
	FunctionType         ObjectType = "function"
	BuiltinFunctionType  ObjectType = "builtinFunction"
	ArrayType            ObjectType = "array"
	HashType             ObjectType = "hash"
	CompiledFunctionType ObjectType = "c ompiledFunction"
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
	Equal(other Object) bool
}

type Comparable interface {
	Compare(other Comparable) (Ordering, error)
}

type Hashable interface {
	GetHashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}
