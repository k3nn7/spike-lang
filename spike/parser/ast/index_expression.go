package ast

import (
	"fmt"
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type IndexExpression struct {
	Token lexer.Token
	Array Expression
	Index Expression
}

func (index *IndexExpression) TokenLiteral() string {
	return index.Token.Literal
}

func (index *IndexExpression) String() string {
	out := strings.Builder{}

	out.WriteString(fmt.Sprintf(
		"(%s[%s])",
		index.Array.String(),
		index.Index.String(),
	))

	return out.String()
}

func (index *IndexExpression) expression() {}
