package ast

import (
	"fmt"
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type BlockStatement struct {
	Token      lexer.Token
	Statements []Statement
}

func (block *BlockStatement) statement() {}

func (block *BlockStatement) TokenLiteral() string {
	return block.Token.Literal
}

func (block *BlockStatement) String() string {
	out := strings.Builder{}
	out.WriteString("{\n")
	for _, statement := range block.Statements {
		out.WriteString(fmt.Sprintf("  %s;\n", statement.String()))
	}
	out.WriteString("}")

	return out.String()
}
