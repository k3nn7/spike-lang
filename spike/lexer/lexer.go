package lexer

import (
	"bufio"
	"io"
	"strings"
)

type TokenIterator interface {
	NextToken() (Token, error)
}

type Lexer struct {
	reader *bufio.Reader
}

func New(reader io.Reader) *Lexer {
	return &Lexer{reader: bufio.NewReader(reader)}
}

func (lexer *Lexer) NextToken() (Token, error) {
	err := lexer.skipWhitespace()
	if err != nil {
		return lexer.handleIOError(err)
	}

	return lexer.readNextToken()
}

func (lexer *Lexer) readNextToken() (Token, error) {
	operator, err := lexer.tryReadTwoCharOperator()
	if err != nil {
		return lexer.handleIOError(err)
	}
	if operator != nil {
		return *operator, nil
	}

	operator, err = lexer.tryReadOneCharOperator()
	if err != nil {
		return lexer.handleIOError(err)
	}
	if operator != nil {
		return *operator, nil
	}

	identifier, err := lexer.tryReadIdentifier()
	if err != nil {
		return lexer.handleIOError(err)
	}
	if identifier != nil {
		return *identifier, nil
	}

	integer, err := lexer.tryReadNumber()
	if err != nil {
		return lexer.handleIOError(err)
	}
	if integer != nil {
		return *integer, nil
	}

	str, err := lexer.tryReadString()
	if err != nil {
		return lexer.handleIOError(err)
	}
	if str != nil {
		return *str, nil
	}

	invalidToken, err := lexer.reader.ReadByte()
	return Token{Invalid, string(invalidToken)}, err
}

func (lexer *Lexer) skipWhitespace() error {
	var err error
	c := make([]byte, 0, 1)

	for c, err = lexer.reader.Peek(1); err == nil && isWhitespace(c[0]); c, err = lexer.reader.Peek(1) {
		_, err2 := lexer.reader.ReadByte()
		if err2 != nil {
			return err2
		}
	}
	return err
}

func (lexer *Lexer) tryReadTwoCharOperator() (*Token, error) {
	twoChars, err := lexer.reader.Peek(2)
	if err == io.EOF {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	t := lookupTwoCharOperator(string(twoChars))
	if t == nil {
		return nil, nil
	}

	_, err = lexer.reader.Read(twoChars)
	return t, err
}

func (lexer *Lexer) tryReadOneCharOperator() (*Token, error) {
	char, err := lexer.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	t := lookupOneCharOperator(string(char))
	if t == nil {
		return nil, nil

	}

	_, err = lexer.reader.Read(char)
	return t, err
}

func (lexer *Lexer) tryReadIdentifier() (*Token, error) {
	char, err := lexer.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if !isIdentifierFirstCharacter(char[0]) {
		return nil, nil
	}

	identifier, err := lexer.readIdentifier()
	if err != nil {
		return nil, err
	}

	keyword := lookupKeyword(identifier)
	if keyword != nil {
		return keyword, nil
	}

	return &Token{Identifier, identifier}, nil
}

func (lexer *Lexer) tryReadNumber() (*Token, error) {
	char, err := lexer.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if !isNumber(char[0]) {
		return nil, nil
	}

	number, err := lexer.readNumber()
	if err != nil {
		return nil, err
	}

	return &Token{Integer, number}, nil
}

func (lexer *Lexer) tryReadString() (*Token, error) {
	char, err := lexer.reader.Peek(1)
	if err != nil {
		return nil, err
	}

	if char[0] != '"' {
		return nil, nil
	}

	_, err = lexer.reader.ReadByte()
	if err != nil {
		return nil, err
	}

	str, err := lexer.readString()
	if err != nil {
		return nil, err
	}

	return &Token{String, str}, nil
}

func (lexer *Lexer) readIdentifier() (string, error) {
	var err error
	c := make([]byte, 0, 1)

	identifier := strings.Builder{}

	for c, err = lexer.reader.Peek(1); err == nil && isIdentifierCharacter(c[0]); c, err = lexer.reader.Peek(1) {
		b, err2 := lexer.reader.ReadByte()
		if err2 != nil {
			return "", err2
		}

		identifier.WriteByte(b)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return identifier.String(), nil
}

func (lexer *Lexer) readNumber() (string, error) {
	var err error
	c := make([]byte, 0, 1)

	number := strings.Builder{}

	for c, err = lexer.reader.Peek(1); err == nil && isNumber(c[0]); c, err = lexer.reader.Peek(1) {
		b, err2 := lexer.reader.ReadByte()
		if err2 != nil {
			return "", err2
		}

		number.WriteByte(b)
	}

	if err != nil && err != io.EOF {
		return "", err
	}

	return number.String(), nil
}

func (lexer *Lexer) readString() (string, error) {
	str := strings.Builder{}
	for {
		b, err := lexer.reader.ReadByte()
		if err != nil {
			return str.String(), err
		}

		if b == '"' {
			return str.String(), nil
		}

		str.WriteByte(b)
	}
}

func (lexer *Lexer) handleIOError(err error) (Token, error) {
	if err == io.EOF {
		return EOFToken, nil
	}

	return Token{}, err
}

func isIdentifierFirstCharacter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentifierCharacter(c byte) bool {
	return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || (c >= 'A' && c <= 'Z')
}

func isNumber(c byte) bool {
	return c >= '0' && c <= '9'
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

func lookupKeyword(literal string) *Token {
	token, ok := keywords[literal]
	if !ok {
		return nil
	}

	return &token
}

func lookupOneCharOperator(literal string) *Token {
	token, ok := oneCharOperators[literal]
	if !ok {
		return nil
	}

	return &token
}

func lookupTwoCharOperator(literal string) *Token {
	token, ok := twoCharOperators[literal]
	if !ok {
		return nil
	}

	return &token
}
