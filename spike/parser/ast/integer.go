package ast

import (
	"fmt"
	"spike-interpreter-go/spike/lexer"
)

type Integer struct {
	Token lexer.Token
	Value int
}

func (integer *Integer) TokenLiteral() string {
	return integer.Token.Literal
}

func (integer *Integer) expression() {}

func (integer *Integer) String() string {
	return fmt.Sprintf("%d", integer.Value)
}
