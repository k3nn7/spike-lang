package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Environment_Get_existing_variable_noOuterEnvironment(t *testing.T) {
	// given
	environment := NewEnvironment()
	value := &True

	environment.Set("key", value)

	// when
	result, err := environment.Get("key")

	// then
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func Test_Environment_Get_not_existing_variable_noOuterEnvironment(t *testing.T) {
	// given
	environment := NewEnvironment()

	// when
	_, err := environment.Get("key")

	// then
	assert.Error(t, err)
}

func Test_Environment_Get_existing_enclosedEnvironment(t *testing.T) {
	// given
	inner := NewEnvironment()
	value := &True

	inner.Set("key", value)
	outer := ExtendEnvironment(inner)

	// when
	result, err := outer.Get("key")

	// then
	assert.NoError(t, err)
	assert.Equal(t, value, result)
}

func Test_Environment_Get_not_existing_enclosedEnvironment(t *testing.T) {
	// given
	inner := NewEnvironment()

	outer := ExtendEnvironment(inner)

	// when
	_, err := outer.Get("key")

	// then
	assert.Error(t, err)
}
