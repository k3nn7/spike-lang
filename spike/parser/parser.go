package parser

import (
	"spike-interpreter-go/spike/lexer"

	"github.com/pkg/errors"
)

type Parser struct {
	lexerInstance *lexer.Lexer
	currentToken  lexer.Token
	peekToken     lexer.Token
}

func New(lexerInstance *lexer.Lexer) *Parser {
	return &Parser{lexerInstance: lexerInstance}
}

func (parser *Parser) advanceToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken, _ = parser.lexerInstance.NextToken()
}

func (parser *Parser) ParseProgram() (Program, error) {
	program := Program{}

	parser.advanceToken()

	for parser.advanceToken(); parser.currentToken.Type != lexer.Eof; parser.advanceToken() {
		statement, err := parser.parseStatement()
		if err != nil {
			return program, err
		}

		program.AddStatement(statement)
	}

	return program, nil
}

func (parser *Parser) parseStatement() (Statement, error) {
	switch parser.currentToken.Type {
	case lexer.Let:
		return parser.parseLetStatement()
	default:
		return nil, errors.Errorf("Invalid statement: %s", parser.currentToken.Literal)
	}
}

func (parser *Parser) parseLetStatement() (Statement, error) {
	letStatement := &LetStatement{Token: parser.currentToken}

	parser.advanceToken()

	if parser.currentToken.Type != lexer.Identifier {
		return letStatement, errors.Errorf("Expecting identifier, got: %s", parser.currentToken.Literal)
	}

	letStatement.Name = &Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	parser.advanceToken()

	if parser.currentToken.Type != lexer.Assign {
		return letStatement, errors.Errorf("Expecting assign operator, got: %s", parser.currentToken.Literal)
	}

	for parser.currentToken.Type != lexer.Semicolon {
		parser.advanceToken()
	}

	return letStatement, nil
}
