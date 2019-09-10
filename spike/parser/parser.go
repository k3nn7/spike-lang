package parser

import (
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"
	"strconv"

	"github.com/pkg/errors"
)

type prefixParseFunc func() (ast.Expression, error)
type infixParseFunc func(expression ast.Expression) (ast.Expression, error)

type Parser struct {
	lexerInstance *lexer.Lexer
	currentToken  lexer.Token
	peekToken     lexer.Token
	prefixParsers map[lexer.TokenType]prefixParseFunc
	infixParsers  map[lexer.TokenType]infixParseFunc
}

func New(lexerInstance *lexer.Lexer) *Parser {
	parser := &Parser{lexerInstance: lexerInstance}
	parser.prefixParsers = make(map[lexer.TokenType]prefixParseFunc)

	parser.addPrefixParser(lexer.Identifier, parser.parseIdentifier)
	parser.addPrefixParser(lexer.Integer, parser.parseInteger)

	return parser
}

func (parser *Parser) addPrefixParser(tokenType lexer.TokenType, prefixParser prefixParseFunc) {
	parser.prefixParsers[tokenType] = prefixParser
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
	case lexer.Return:
		return parser.parseReturnStatement()
	default:
		return parser.parseExpressionStatement()
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

	parser.advanceToken()

	expression, err := parser.parseExpression()
	letStatement.Value = expression

	if parser.peekToken.Type == lexer.Semicolon {
		parser.advanceToken()
	}

	return letStatement, err
}

func (parser *Parser) parseReturnStatement() (ast.Statement, error) {
	returnStatement := &ast.ReturnStatement{Token: parser.currentToken}

	parser.advanceToken()

	for parser.currentToken.Type != lexer.Semicolon {
		parser.advanceToken()
	}

	return returnStatement, nil
}

func (parser *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	returnStatement := &ast.ExpressionStatement{}
	returnStatement.Expression, _ = parser.parseExpression()

	if parser.peekToken.Type == lexer.Semicolon {
		parser.advanceToken()
	}

	return returnStatement, nil
}

func (parser *Parser) parseExpression() (ast.Expression, error) {
	parsePrefixExpression, ok := parser.prefixParsers[parser.currentToken.Type]
	if ok {
		return parsePrefixExpression()
	}

	return nil, errors.New("invalid expression")
}

func (parser *Parser) parseIdentifier() (ast.Expression, error) {
	return &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}, nil
}

func (parser *Parser) parseInteger() (ast.Expression, error) {
	value, err := strconv.Atoi(parser.currentToken.Literal)
	if err != nil {
		return nil, err
	}

	return &ast.Integer{
		Token: parser.currentToken,
		Value: value,
	}, nil
}
