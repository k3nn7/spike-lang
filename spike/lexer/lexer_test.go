package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Lexer_token(t *testing.T) {
	testCases := []struct {
		input         string
		expectedToken Token
	}{
		{
			input:         "true",
			expectedToken: TrueToken,
		},
		{
			input:         "false",
			expectedToken: FalseToken,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.input, func(t *testing.T) {
			l := New(strings.NewReader(testCase.input))

			token, err := l.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, testCase.expectedToken, token)

			token, err = l.NextToken()
			assert.NoError(t, err)
			assert.Equal(t, EOFToken, token)
		})
	}
}

func Test_Lexer_code_sample(t *testing.T) {
	// given
	input := strings.NewReader(`
let variable = (10 + 20) * 5; 
return variable2 ! VAR3 - true false / < > == !=
<= >= || &&
`)
	expectedTokens := []Token{
		LetToken,
		{Identifier, "variable"},
		AssignToken,
		LeftParenthesisToken,
		{Integer, "10"},
		PlusToken,
		{Integer, "20"},
		RightParenthesisToken,
		AsteriskToken,
		{Integer, "5"},
		SemicolonToken,
		ReturnToken,
		{Identifier, "variable2"},
		BangToken,
		{Identifier, "VAR3"},
		MinusToken,
		TrueToken,
		FalseToken,
		SlashToken,
		LessThanToken,
		GreaterThanToken,
		EqualToken,
		NotEqualToken,
		LessOrEqualToken,
		GreaterOrEqualToken,
		OrToken,
		AndToken,
	}

	lexer := New(input)

	// when
	tokens, err := iteratorToSlice(lexer)

	// then
	assert.NoError(t, err)
	assert.Exactly(t, expectedTokens, tokens)
}

func Test_Lexer_invalidToken(t *testing.T) {
	// given
	input := strings.NewReader("^")
	expectedTokens := []Token{
		{Invalid, "^"},
	}

	lexer := New(input)

	// when
	tokens, err := iteratorToSlice(lexer)

	// then
	assert.NoError(t, err)
	assert.Exactly(t, expectedTokens, tokens)
}

func iteratorToSlice(iterator TokenIterator) ([]Token, error) {
	result := make([]Token, 0)

	for token, err := iterator.NextToken(); token != EOFToken; token, err = iterator.NextToken() {
		if err != nil {
			return nil, err
		}

		result = append(result, token)
	}

	return result, nil
}
