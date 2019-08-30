package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

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
