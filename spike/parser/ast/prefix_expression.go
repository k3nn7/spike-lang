package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type PrefixExpression struct {
	Token    lexer.Token
	Operator string
	Right    Expression
}

func (expression *PrefixExpression) expression() {}

func (expression *PrefixExpression) TokenLiteral() string {
	return expression.Token.Literal
}

func (expression *PrefixExpression) String() string {
	out := strings.Builder{}
	out.WriteString("(")
	out.WriteString(expression.Token.Literal)
	out.WriteString(expression.Right.String())
	out.WriteString(")")
	return out.String()
}
