package ast

import (
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type ReturnStatement struct {
	Token  lexer.Token
	Result Expression
}

func (returnStatement *ReturnStatement) TokenLiteral() string {
	out := strings.Builder{}
	out.WriteString("return ")
	out.WriteString(returnStatement.Result.TokenLiteral())

	return out.String()
}

func (returnStatement *ReturnStatement) statement() {
}
