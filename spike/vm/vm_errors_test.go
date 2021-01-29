package vm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Run_withError(t *testing.T) {
	testCases := []struct {
		code          string
		expectedError string
	}{
		{
			code:          `let f = fn(a) { a }; f(1, 2)`,
			expectedError: "mismatched number of function call arguments. Expected 1, got 2",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.code, func(t *testing.T) {
			_, err := runInVM(testCase.code)
			assert.Error(t, err)
			assert.EqualError(t, err, testCase.expectedError)
		})
	}
}
