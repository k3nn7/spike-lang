package parser

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	out := strings.Builder{}

	for _, statement := range program.Statements {
		out.WriteString(statement.TokenLiteral())
		out.WriteByte('\n')
	}

	return out.String()
}

func (program *Program) AddStatement(statement Statement) {
	program.Statements = append(program.Statements, statement)
}

type LetStatement struct {
	Token lexer.Token
	Name  *Identifier
	Value Expression
}

func (let *LetStatement) TokenLiteral() string {
	out := strings.Builder{}
	out.WriteString(let.Name.TokenLiteral())
	out.WriteString(" = ")
	out.WriteString(let.Value.TokenLiteral())

	return out.String()
}

func (let *LetStatement) statement() {
}

type Identifier struct {
	Token lexer.Token
	Value string
}

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Value
}

func (identifier *Identifier) expression() {}
