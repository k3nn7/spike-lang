package object

type Return struct {
	Value Object
}

func (r *Return) Type() ObjectType {
	return ReturnType
}

func (r *Return) Inspect() string {
	return r.Value.Inspect()
}

func (r *Return) Equal(other Object) bool {
	return false
}
