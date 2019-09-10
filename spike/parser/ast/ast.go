package ast

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statement()
}

type Expression interface {
	Node
	expression()
}

type ExpressionStatement struct {
	expression Expression
}

func (statement *ExpressionStatement) TokenLiteral() string {
	return "expression"
}

func (statement ExpressionStatement) statement() {
}

func (statement *ExpressionStatement) String() string {
	return statement.expression.String()
}
