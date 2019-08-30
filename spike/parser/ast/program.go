package ast

import "strings"

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
