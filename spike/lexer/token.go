package lexer

type Token struct {
	Type    TokenType
	Literal string
}

type TokenType string

// Operators
const (
	Assign           TokenType = "assign"
	LeftParenthesis  TokenType = "leftParenthesis"
	RightParenthesis TokenType = "rightParenthesis"
	Plus             TokenType = "plus"
	Asterisk         TokenType = "*"
)

var operators = map[string]TokenType{
	"=": Assign,
	"(": LeftParenthesis,
	")": RightParenthesis,
	"+": Plus,
	"*": Asterisk,
}

// Keywords
const (
	Let TokenType = "let"
)

var keywords = map[string]TokenType{
	"let": Let,
}

// Other
const (
	Eof        TokenType = "eof"
	Invalid    TokenType = "invalid"
	Identifier TokenType = "identifier"
	Integer    TokenType = "integer"
)
