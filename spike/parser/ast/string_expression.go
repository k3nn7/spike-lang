package ast

import (
	"fmt"
	"spike-interpreter-go/spike/lexer"
)

type String struct {
	Token lexer.Token
	Value string
}

func (str *String) TokenLiteral() string {
	return str.Token.Literal
}

func (str *String) String() string {
	return fmt.Sprintf("\"%s\"", str.Value)
}

func (str *String) expression() {}
