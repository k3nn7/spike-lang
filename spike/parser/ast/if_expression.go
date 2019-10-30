package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type IfExpression struct {
	Token     lexer.Token
	Condition Expression
	Then      Statement
	Else      Statement
}

func (expression *IfExpression) expression() {}

func (expression *IfExpression) TokenLiteral() string {
	return expression.Token.Literal
}

func (expression *IfExpression) String() string {
	out := strings.Builder{}
	out.WriteString("if ")
	out.WriteString(expression.Condition.String())
	out.WriteString(" ")
	out.WriteString(expression.Then.String())
	if expression.Else != nil {
		out.WriteString(" else ")
		out.WriteString(expression.Else.String())
	}

	return out.String()
}
