package eval

import (
	"spike-interpreter-go/spike/eval/object"

	"github.com/pkg/errors"
)

type Environment struct {
	variables map[string]object.Object
}

func NewEnvironment() *Environment {
	variables := make(map[string]object.Object)
	return &Environment{variables: variables}
}

func (e Environment) Set(name string, value object.Object) {
	e.variables[name] = value
}

func (e Environment) Get(name string) (object.Object, error) {
	if value, ok := e.variables[name]; ok {
		return value, nil
	}

	return nil, errors.Errorf("undefined identifier: %s", name)
}
