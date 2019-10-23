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
	Minus            TokenType = "minus"
	Asterisk         TokenType = "asterisk"
	Bang             TokenType = "bang"
	Slash            TokenType = "slash"
)

var operators = map[string]Token{
	"=": AssignToken,
	"(": LeftParenthesisToken,
	")": RightParenthesisToken,
	"+": PlusToken,
	"-": MinusToken,
	"*": AsteriskToken,
	";": SemicolonToken,
	"!": BangToken,
	"/": SlashToken,
}

// Keywords
const (
	Let    TokenType = "let"
	Return TokenType = "return"
	True   TokenType = "true"
	False  TokenType = "false"
)

var keywords = map[string]Token{
	"let":    LetToken,
	"return": ReturnToken,
	"true":   TrueToken,
	"false":  FalseToken,
}

// Other
const (
	Semicolon  TokenType = "semicolon"
	Eof        TokenType = "eof"
	Invalid    TokenType = "invalid"
	Identifier TokenType = "identifier"
	Integer    TokenType = "integer"
)

// Predefined tokens
var (
	EOFToken              = Token{Type: Eof, Literal: ""}
	TrueToken             = Token{Type: True, Literal: "true"}
	FalseToken            = Token{Type: False, Literal: "false"}
	LetToken              = Token{Type: Let, Literal: "let"}
	ReturnToken           = Token{Type: Return, Literal: "return"}
	AssignToken           = Token{Type: Assign, Literal: "="}
	LeftParenthesisToken  = Token{Type: LeftParenthesis, Literal: "("}
	RightParenthesisToken = Token{Type: RightParenthesis, Literal: ")"}
	PlusToken             = Token{Type: Plus, Literal: "+"}
	MinusToken            = Token{Type: Minus, Literal: "-"}
	AsteriskToken         = Token{Type: Asterisk, Literal: "*"}
	SemicolonToken        = Token{Type: Semicolon, Literal: ";"}
	BangToken             = Token{Type: Bang, Literal: "!"}
	SlashToken            = Token{Type: Slash, Literal: "/"}
)
