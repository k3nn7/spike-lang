package object

var NullObject = Null{}

type Null struct{}

func (null *Null) Type() ObjectType {
	return NullType
}

func (null *Null) Inspect() string {
	return "null"
}
