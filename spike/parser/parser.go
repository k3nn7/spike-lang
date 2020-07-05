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
	alternative
	conjunction
	inequality
	equals
	sum
	product
	prefix
)

var precedences = map[lexer.TokenType]int{
	lexer.Plus:           sum,
	lexer.Minus:          sum,
	lexer.Asterisk:       product,
	lexer.Slash:          product,
	lexer.Equal:          equals,
	lexer.NotEqual:       equals,
	lexer.LessThan:       inequality,
	lexer.GreaterThan:    inequality,
	lexer.LessOrEqual:    inequality,
	lexer.GreaterOrEqual: inequality,
	lexer.And:            conjunction,
	lexer.Or:             alternative,
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
	parser.addPrefixParser(lexer.True, parser.parseBoolean)
	parser.addPrefixParser(lexer.False, parser.parseBoolean)
	parser.addPrefixParser(lexer.Bang, parser.parsePrefixExpression)
	parser.addPrefixParser(lexer.Minus, parser.parsePrefixExpression)
	parser.addPrefixParser(lexer.LeftParenthesis, parser.parseGroupedExpression)
	parser.addPrefixParser(lexer.If, parser.parseIfExpression)
	parser.addPrefixParser(lexer.Fn, parser.parseFunctionExpression)

	parser.addInfixParser(lexer.Plus, parser.parseInfixExpression)
	parser.addInfixParser(lexer.Asterisk, parser.parseInfixExpression)
	parser.addInfixParser(lexer.Minus, parser.parseInfixExpression)
	parser.addInfixParser(lexer.Slash, parser.parseInfixExpression)
	parser.addInfixParser(lexer.Equal, parser.parseInfixExpression)
	parser.addInfixParser(lexer.NotEqual, parser.parseInfixExpression)
	parser.addInfixParser(lexer.GreaterThan, parser.parseInfixExpression)
	parser.addInfixParser(lexer.LessThan, parser.parseInfixExpression)
	parser.addInfixParser(lexer.GreaterOrEqual, parser.parseInfixExpression)
	parser.addInfixParser(lexer.LessOrEqual, parser.parseInfixExpression)
	parser.addInfixParser(lexer.Or, parser.parseInfixExpression)
	parser.addInfixParser(lexer.And, parser.parseInfixExpression)

	return parser
}

func (parser *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{}

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
	case lexer.LeftBrace:
		return parser.parseBlockStatement()
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

func (parser *Parser) parseIfExpression() (ast.Expression, error) {
	ifExpression := &ast.IfExpression{Token: parser.currentToken}

	parser.advanceToken()
	if parser.currentToken.Type != lexer.LeftParenthesis {
		return ifExpression, errors.Errorf("expected left parenthesis, got %s", parser.currentToken.Type)
	}

	parser.advanceToken()
	condition, err := parser.parseExpression(lowest)
	if err != nil {
		return ifExpression, err
	}
	ifExpression.Condition = condition

	parser.advanceToken()
	if parser.currentToken.Type != lexer.RightParenthesis {
		return ifExpression, errors.Errorf("expected right parenthesis, got %s", parser.currentToken.Type)
	}

	parser.advanceToken()
	statement, err := parser.parseStatement()
	if err != nil {
		return ifExpression, err
	}
	ifExpression.Then = statement

	if parser.currentToken.Type != lexer.Else {
		return ifExpression, nil
	}

	parser.advanceToken()
	statement, err = parser.parseStatement()
	if err != nil {
		return ifExpression, err
	}
	ifExpression.Else = statement

	return ifExpression, nil
}

func (parser *Parser) parseFunctionExpression() (ast.Expression, error) {
	functionExpression := &ast.FunctionExpression{Token: parser.currentToken}

	parser.advanceToken()
	if parser.currentToken.Type != lexer.LeftParenthesis {
		return functionExpression, errors.Errorf("expected left parenthesis, got %s", parser.currentToken.Type)
	}

	for {
		parser.advanceToken()
		if parser.currentToken.Type == lexer.RightParenthesis {
			break
		}

		if parser.currentToken.Type != lexer.Identifier {
			return functionExpression, errors.Errorf("expected identifier, got %s", parser.currentToken.Type)
		}

		identifier, err := parser.parseIdentifier()
		if err != nil {
			return functionExpression, err
		}
		functionExpression.Parameters = append(functionExpression.Parameters, identifier.(*ast.Identifier))

		parser.advanceToken()
		if parser.currentToken.Type == lexer.RightParenthesis {
			break
		}

		if parser.currentToken.Type != lexer.Comma {
			return functionExpression, errors.Errorf("expected comma, got %s", parser.currentToken.Type)
		}
	}

	parser.advanceToken()
	statement, err := parser.parseStatement()
	if err != nil {
		return functionExpression, err
	}

	functionExpression.Body = statement

	return functionExpression, nil
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
	value, err := strconv.ParseInt(parser.currentToken.Literal, 10, 64)
	if err != nil {
		return nil, err
	}

	expression := &ast.Integer{
		Token: parser.currentToken,
		Value: value,
	}

	return expression, nil
}

func (parser *Parser) parseBoolean() (ast.Expression, error) {
	if parser.currentToken == lexer.TrueToken {
		return &ast.Boolean{Token: parser.currentToken, Value: true}, nil
	}

	return &ast.Boolean{Token: parser.currentToken, Value: false}, nil
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

func (parser *Parser) parseGroupedExpression() (ast.Expression, error) {
	parser.advanceToken()

	expression, _ := parser.parseExpression(lowest)

	parser.advanceToken()

	return expression, nil
}

func (parser *Parser) parseBlockStatement() (ast.Statement, error) {
	blockStatement := &ast.BlockStatement{
		Token:      parser.currentToken,
		Statements: make([]ast.Statement, 0),
	}

	for parser.advanceToken(); parser.currentToken.Type != lexer.RightBrace; parser.advanceToken() {
		statement, err := parser.parseStatement()
		if err != nil {
			return blockStatement, err
		}
		blockStatement.Statements = append(blockStatement.Statements, statement)

		if parser.peekToken.Type == lexer.Semicolon {
			parser.advanceToken()
		}
	}

	parser.advanceToken()
	return blockStatement, nil
}
