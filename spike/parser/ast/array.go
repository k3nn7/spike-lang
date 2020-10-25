package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type Array struct {
	Token    lexer.Token
	Elements []Expression
}

func (array *Array) TokenLiteral() string {
	return array.Token.Literal
}

func (array *Array) String() string {
	out := strings.Builder{}

	out.WriteString("[")

	for i, element := range array.Elements {
		out.WriteString(element.String())
		if i < len(array.Elements)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString("]")

	return out.String()
}

func (array *Array) expression() {
}
