package ast

import (
	"fmt"
	"sort"
	"spike-interpreter-go/spike/lexer"
	"strings"
)

type Hash struct {
	Token lexer.Token
	Pairs map[Expression]Expression
}

func (hash *Hash) TokenLiteral() string {
	return hash.Token.Literal
}

func (hash *Hash) String() string {
	out := strings.Builder{}

	pairs := make([]string, 0, len(hash.Pairs))
	for key, val := range hash.Pairs {
		pairs = append(pairs, fmt.Sprintf(
			"%s: %s",
			key.String(),
			val.String(),
		))
	}

	sort.Strings(pairs)

	out.WriteString(fmt.Sprintf(
		"{%s}",
		strings.Join(pairs, ", "),
	))

	return out.String()
}

func (hash *Hash) expression() {
}
