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
	";": Semicolon,
}

// Keywords
const (
	Let    TokenType = "let"
	Return TokenType = "return"
)

var keywords = map[string]TokenType{
	"let":    Let,
	"return": Return,
}

// Other
const (
	Semicolon  TokenType = "semicolon"
	Eof        TokenType = "eof"
	Invalid    TokenType = "invalid"
	Identifier TokenType = "identifier"
	Integer    TokenType = "integer"
)
