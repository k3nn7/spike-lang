package object

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Boolean_Equal(t *testing.T) {
	testCases := []struct {
		left           Object
		right          Object
		expectedResult bool
	}{
		{
			left:           &True,
			right:          &False,
			expectedResult: false,
		},
		{
			left:           &False,
			right:          &True,
			expectedResult: false,
		},
		{
			left:           &True,
			right:          &True,
			expectedResult: true,
		},
		{
			left:           &False,
			right:          &False,
			expectedResult: true,
		},
	}

	for _, testCase := range testCases {
		testCaseName := fmt.Sprintf("%v == %v => %t", testCase.left, testCase.right, testCase.expectedResult)
		t.Run(testCaseName, func(t *testing.T) {
			result, err := testCase.left.Equal(testCase.right)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func Test_Boolean_Equal_errors(t *testing.T) {
	testCases := []struct {
		left          Object
		right         Object
		expectedError error
	}{
		{
			left:          &True,
			right:         &Integer{Value: 10},
			expectedError: NotComparableError,
		},
	}

	for _, testCase := range testCases {
		testCaseName := fmt.Sprintf("%#v == %#v => %s", testCase.left, testCase.right, testCase.expectedError)
		t.Run(testCaseName, func(t *testing.T) {
			_, err := testCase.left.Equal(testCase.right)

			assert.Equal(t, testCase.expectedError, err)
		})
	}
}
