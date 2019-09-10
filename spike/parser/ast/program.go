package ast

import "strings"

type Program struct {
	Statements []Statement
}

func (program *Program) TokenLiteral() string {
	return "program"
}

func (program *Program) AddStatement(statement Statement) {
	program.Statements = append(program.Statements, statement)
}

func (program *Program) String() string {
	out := strings.Builder{}

	for _, statement := range program.Statements {
		out.WriteString(statement.String())
		out.WriteByte('\n')
	}

	return out.String()
}
