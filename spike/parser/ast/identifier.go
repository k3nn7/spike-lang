package ast

import "spike-interpreter-go/spike/lexer"

type Identifier struct {
	Token lexer.Token
	Value string
}

func (identifier *Identifier) TokenLiteral() string {
	return identifier.Value
}

func (identifier *Identifier) expression() {}
