package object

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Integer_Equal(t *testing.T) {
	testCases := []struct {
		left           Object
		right          Object
		expectedResult bool
	}{
		{
			left:           &Integer{Value: 10},
			right:          &Integer{Value: -10},
			expectedResult: false,
		},
		{
			left:           &Integer{Value: 255},
			right:          &Integer{Value: 255},
			expectedResult: true,
		},
	}

	for _, testCase := range testCases {
		testCaseName := fmt.Sprintf("%#v == %#v => %t", testCase.left, testCase.right, testCase.expectedResult)
		t.Run(testCaseName, func(t *testing.T) {
			result, err := testCase.left.Equal(testCase.right)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}

func Test_Integer_Equal_errors(t *testing.T) {
	testCases := []struct {
		left          Object
		right         Object
		expectedError error
	}{
		{
			left:          &Integer{Value: 127},
			right:         &NullObject,
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

func Test_Integer_Compare(t *testing.T) {
	testCases := []struct {
		left           Comparable
		right          Comparable
		expectedResult Ordering
	}{
		{
			left:           &Integer{Value: 10},
			right:          &Integer{Value: -10},
			expectedResult: GT,
		},
		{
			left:           &Integer{Value: 255},
			right:          &Integer{Value: 255},
			expectedResult: EQ,
		},
		{
			left:           &Integer{Value: 127},
			right:          &Integer{Value: 255},
			expectedResult: LT,
		},
	}

	for _, testCase := range testCases {
		testCaseName := fmt.Sprintf("%#v <> %#v => %v", testCase.left, testCase.right, testCase.expectedResult)
		t.Run(testCaseName, func(t *testing.T) {
			result, err := testCase.left.Compare(testCase.right)

			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}
