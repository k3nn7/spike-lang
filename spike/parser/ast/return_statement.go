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
	return returnStatement.Token.Literal
}

func (returnStatement *ReturnStatement) statement() {
}

func (returnStatement *ReturnStatement) String() string {
	out := strings.Builder{}
	out.WriteString("return ")
	out.WriteString(returnStatement.Result.String())

	return out.String()
}
