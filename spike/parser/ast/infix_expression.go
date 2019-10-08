package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type InfixExpression struct {
	Token    lexer.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (expression *InfixExpression) expression() {}

func (expression *InfixExpression) TokenLiteral() string {
	return expression.Token.Literal
}

func (expression *InfixExpression) String() string {
	out := strings.Builder{}
	out.WriteString("(")
	out.WriteString(expression.Left.String())
	out.WriteString(" ")
	out.WriteString(expression.Operator)
	out.WriteString(" ")
	out.WriteString(expression.Right.String())
	out.WriteString(")")

	return out.String()
}
