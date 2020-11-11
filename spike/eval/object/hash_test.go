package object

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHash_Equal(t *testing.T) {
	hash1 := &Hash{Pairs: map[HashKey]HashPair{
		HashKey{
			Type:  IntegerType,
			Value: 55,
		}: {
			Key:   &Integer{Value: 55},
			Value: &Integer{Value: 10},
		},
		HashKey{
			Type:  IntegerType,
			Value: 10,
		}: {
			Key:   &Integer{Value: 10},
			Value: &Integer{Value: 99},
		},
	}}

	hash2 := &Hash{Pairs: map[HashKey]HashPair{
		HashKey{
			Type:  IntegerType,
			Value: 10,
		}: {
			Key:   &Integer{Value: 10},
			Value: &Integer{Value: 99},
		},
		HashKey{
			Type:  IntegerType,
			Value: 55,
		}: {
			Key:   &Integer{Value: 55},
			Value: &Integer{Value: 10},
		},
	}}

	hash3 := &Hash{Pairs: map[HashKey]HashPair{
		HashKey{
			Type:  IntegerType,
			Value: 55,
		}: {
			Key:   &Integer{Value: 55},
			Value: &Integer{Value: 10},
		},
		HashKey{
			Type:  IntegerType,
			Value: 10,
		}: {
			Key:   &Integer{Value: 10},
			Value: &Integer{Value: 9},
		},
	}}

	other := &Integer{Value: 10}

	assert.False(t, hash1.Equal(other))
	assert.True(t, hash1.Equal(hash2))
	assert.False(t, hash1.Equal(hash3))
}
