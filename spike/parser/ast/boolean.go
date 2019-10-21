package ast

import "spike-interpreter-go/spike/lexer"

type Boolean struct {
	Token lexer.Token
	Value bool
}

func (boolean *Boolean) expression() {}

func (boolean *Boolean) TokenLiteral() string {
	return boolean.Token.Literal
}

func (boolean *Boolean) String() string {
	if boolean.Value {
		return "true"
	}

	return "false"
}
