package parser

import (
	"spike-interpreter-go/spike/lexer"
	"spike-interpreter-go/spike/parser/ast"
	"strconv"

	"github.com/pkg/errors"
)

type prefixParseFunc func() (ast.Expression, error)
type infixParseFunc func(expression ast.Expression) (ast.Expression, error)

const (
	lowest = iota
	equals
	inequality
	sum
	product
	prefix
	call
)

var precedences = map[lexer.TokenType]int{
	lexer.Plus:     sum,
	lexer.Minus:    sum,
	lexer.Asterisk: product,
}

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
	parser.infixParsers = make(map[lexer.TokenType]infixParseFunc)

	parser.addPrefixParser(lexer.Identifier, parser.parseIdentifier)
	parser.addPrefixParser(lexer.Integer, parser.parseInteger)
	parser.addPrefixParser(lexer.Bang, parser.parsePrefixExpression)
	parser.addPrefixParser(lexer.Minus, parser.parsePrefixExpression)

	parser.addInfixParser(lexer.Plus, parser.parseInfixExpression)
	parser.addInfixParser(lexer.Asterisk, parser.parseInfixExpression)

	return parser
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

		if parser.peekToken.Type == lexer.Semicolon {
			parser.advanceToken()
		}
	}

	return program, nil
}

func (parser *Parser) addPrefixParser(tokenType lexer.TokenType, prefixParser prefixParseFunc) {
	parser.prefixParsers[tokenType] = prefixParser
}

func (parser *Parser) addInfixParser(tokenType lexer.TokenType, infixParser infixParseFunc) {
	parser.infixParsers[tokenType] = infixParser
}

func (parser *Parser) advanceToken() {
	parser.currentToken = parser.peekToken
	parser.peekToken, _ = parser.lexerInstance.NextToken()
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

	expression, err := parser.parseExpression(lowest)
	letStatement.Value = expression

	return letStatement, err
}

func (parser *Parser) parseReturnStatement() (ast.Statement, error) {
	returnStatement := &ast.ReturnStatement{Token: parser.currentToken}

	parser.advanceToken()

	expression, _ := parser.parseExpression(lowest)
	returnStatement.Result = expression

	return returnStatement, nil
}

func (parser *Parser) parseExpressionStatement() (*ast.ExpressionStatement, error) {
	var err error
	statement := &ast.ExpressionStatement{}
	statement.Expression, err = parser.parseExpression(lowest)

	if err != nil {
		return nil, err
	}

	if parser.peekToken.Type == lexer.Semicolon {
		parser.advanceToken()
	}

	return statement, nil
}

func (parser *Parser) parseExpression(precedence int) (ast.Expression, error) {
	var expression ast.Expression
	var err error
	parsePrefixExpression, ok := parser.prefixParsers[parser.currentToken.Type]
	if !ok {
		return expression, errors.Errorf("%q is not a valid prefix expression", parser.currentToken.Literal)
	}

	expression, err = parsePrefixExpression()
	if err != nil {
		return expression, err
	}

	for parser.peekToken.Type != lexer.Semicolon && precedence < precedences[parser.peekToken.Type] {
		parseInfixExpression, ok := parser.infixParsers[parser.peekToken.Type]
		if !ok {
			return expression, nil
		}

		parser.advanceToken()

		expression, err = parseInfixExpression(expression)
	}

	return expression, err
}

func (parser *Parser) parseIdentifier() (ast.Expression, error) {
	expression := &ast.Identifier{Token: parser.currentToken, Value: parser.currentToken.Literal}

	return expression, nil
}

func (parser *Parser) parseInteger() (ast.Expression, error) {
	value, err := strconv.Atoi(parser.currentToken.Literal)
	if err != nil {
		return nil, err
	}

	expression := &ast.Integer{
		Token: parser.currentToken,
		Value: value,
	}

	return expression, nil
}

func (parser *Parser) parsePrefixExpression() (ast.Expression, error) {
	prefixExpression := &ast.PrefixExpression{
		Token:    parser.currentToken,
		Operator: parser.currentToken.Literal,
	}

	parser.advanceToken()
	right, err := parser.parseExpression(prefix)
	prefixExpression.Right = right

	return prefixExpression, err
}

func (parser *Parser) parseInfixExpression(left ast.Expression) (ast.Expression, error) {
	expression := &ast.InfixExpression{
		Token:    parser.currentToken,
		Left:     left,
		Operator: parser.currentToken.Literal,
	}

	precedence, _ := precedences[parser.currentToken.Type]

	parser.advanceToken()
	expression.Right, _ = parser.parseExpression(precedence)

	return expression, nil
}
