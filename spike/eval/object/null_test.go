package object

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Null_Equal(t *testing.T) {
	testCases := []struct {
		left           Object
		right          Object
		expectedResult bool
	}{
		{
			left:           &NullObject,
			right:          &NullObject,
			expectedResult: true,
		},
	}

	for _, testCase := range testCases {
		testCaseName := fmt.Sprintf("%v == %v => %t", testCase.left, testCase.right, testCase.expectedResult)
		t.Run(testCaseName, func(t *testing.T) {
			result := testCase.left.Equal(testCase.right)

			assert.Equal(t, testCase.expectedResult, result)
		})
	}
}
