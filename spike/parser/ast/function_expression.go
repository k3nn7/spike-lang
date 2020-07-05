package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type FunctionExpression struct {
	Token      lexer.Token
	Parameters []*Identifier
	Body       Statement
}

func (function *FunctionExpression) expression() {}

func (function *FunctionExpression) TokenLiteral() string {
	return function.Token.Literal
}

func (function *FunctionExpression) String() string {
	out := strings.Builder{}

	out.WriteString(function.Token.Literal)
	out.WriteString(" (")
	for i, parameter := range function.Parameters {
		out.WriteString(parameter.String())

		if i < len(function.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")

	out.WriteString(function.Body.String())

	return out.String()
}
