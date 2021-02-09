package object

import (
	"spike-interpreter-go/spike/parser/ast"
	"strings"
)

type Function struct {
	Parameters  []*ast.Identifier
	Body        ast.Statement
	Environment *Environment
}

func (function *Function) Type() ObjectType {
	return FunctionType
}

func (function *Function) Inspect() string {
	out := strings.Builder{}

	return out.String()
}

func (function *Function) Equal(other Object) bool {
	return other == function
}
