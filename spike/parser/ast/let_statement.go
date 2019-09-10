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
	return let.Token.Literal
}

func (let *LetStatement) statement() {
}

func (let *LetStatement) String() string {
	out := strings.Builder{}
	out.WriteString(let.Token.Literal)
	out.WriteString(" ")
	out.WriteString(let.Name.String())
	out.WriteString(" = ")
	out.WriteString(let.Value.String())

	return out.String()
}
