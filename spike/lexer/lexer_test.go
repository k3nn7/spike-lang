package lexer

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Lexer_code_sample(t *testing.T) {
	// given
	input := strings.NewReader(" let variable = (10 + 20) * 5 ")
	expectedTokens := []Token{
		{Let, "let"},
		{Identifier, "variable"},
		{Assign, "="},
		{LeftParenthesis, "("},
		{Integer, "10"},
		{Plus, "+"},
		{Integer, "20"},
		{RightParenthesis, ")"},
		{Asterisk, "*"},
		{Integer, "5"},
	}

	lexer := NewLexer(input)

	// when
	tokens, err := iteratorToSlice(lexer)

	// then
	assert.NoError(t, err)
	assert.Exactly(t, expectedTokens, tokens)
}

func Test_Lexer_invalidToken(t *testing.T) {
	// given
	input := strings.NewReader("!")
	expectedTokens := []Token{
		{Invalid, "!"},
	}

	lexer := NewLexer(input)

	// when
	tokens, err := iteratorToSlice(lexer)

	// then
	assert.NoError(t, err)
	assert.Exactly(t, expectedTokens, tokens)
}

func iteratorToSlice(iterator TokenIterator) ([]Token, error) {
	result := make([]Token, 0)

	for token, err := iterator.NextToken(); token.Type != Eof; token, err = iterator.NextToken() {
		if err != nil {
			return nil, err
		}

		result = append(result, token)
	}

	return result, nil
}
