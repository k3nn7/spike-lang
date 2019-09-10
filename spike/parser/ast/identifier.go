package ast

import "spike-interpreter-go/spike/lexer"

type Identifier struct {
	Token lexer.Token
	Value string
}

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Token.Literal
}

func (identifier *Identifier) expression() {}

func (identifier *Identifier) String() string {
	return identifier.Value
}
