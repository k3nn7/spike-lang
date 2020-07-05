package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type CallExpression struct {
	Token     lexer.Token
	Function  Expression
	Arguments []Expression
}

func (call *CallExpression) TokenLiteral() string {
	return call.Token.Literal
}

func (call *CallExpression) String() string {
	out := strings.Builder{}

	out.WriteString(call.Function.String())
	out.WriteString("(")

	for i, argument := range call.Arguments {
		out.WriteString(argument.String())
		if i < len(call.Arguments)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(");")

	return out.String()
}

func (call *CallExpression) expression() {}
