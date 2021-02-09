package object

import (
	"github.com/pkg/errors"
)

type Environment struct {
	variables map[string]Object
	inner     *Environment
}

func NewEnvironment() *Environment {
	variables := make(map[string]Object)
	return &Environment{variables: variables}
}

func ExtendEnvironment(environment *Environment) *Environment {
	variables := make(map[string]Object)
	return &Environment{variables: variables, inner: environment}
}

func (e Environment) Set(name string, value Object) {
	e.variables[name] = value
}

func (e Environment) Get(name string) (Object, error) {
	if value, ok := e.variables[name]; ok {
		return value, nil
	}

	if e.inner != nil {
		return e.inner.Get(name)
	}

	return nil, errors.Errorf("undefined identifier: %s", name)
}
