package ast

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}
