package parser

import (
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"

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

func (parser *Parser) ParseProgram() (ast.Program, error) {
	program := ast.Program{}

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

func (parser *Parser) parseStatement() (ast.Statement, error) {
	switch parser.currentToken.Type {
	case lexer.Let:
		return parser.parseLetStatement()
	default:
		return nil, errors.Errorf("Invalid statement: %s", parser.currentToken.Literal)
	}
}

func (parser *Parser) parseLetStatement() (ast.Statement, error) {
	letStatement := &ast.LetStatement{Token: parser.currentToken}

	parser.advanceToken()

	if parser.currentToken.Type != lexer.Identifier {
		return letStatement, errors.Errorf("expected identifier, got %s", parser.currentToken.Type)
	}

	letStatement.Name = &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	parser.advanceToken()

	if parser.currentToken.Type != lexer.Assign {
		return letStatement, errors.Errorf("expected assign operator, got %s", parser.currentToken.Type)
	}

	for parser.currentToken.Type != lexer.Semicolon {
		parser.advanceToken()
	}

	return letStatement, nil
}
