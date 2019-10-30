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
	LessThan         TokenType = "lessThan"
	GreaterThan      TokenType = "greaterThan"
	LessOrEqual      TokenType = "lessOrEqual"
	GreaterOrEqual   TokenType = "greaterOrEqual"
	Equal            TokenType = "equal"
	NotEqual         TokenType = "notEqual"
	And              TokenType = "and"
	Or               TokenType = "or"
	LeftBrace        TokenType = "leftBrace"
	RightBrace       TokenType = "rightBrace"
)

var oneCharOperators = map[string]Token{
	"=": AssignToken,
	"(": LeftParenthesisToken,
	")": RightParenthesisToken,
	"+": PlusToken,
	"-": MinusToken,
	"*": AsteriskToken,
	";": SemicolonToken,
	"!": BangToken,
	"/": SlashToken,
	"<": LessThanToken,
	">": GreaterThanToken,
	"{": LeftBraceToken,
	"}": RightBraceToken,
}

var twoCharOperators = map[string]Token{
	"==": EqualToken,
	"!=": NotEqualToken,
	"<=": LessOrEqualToken,
	">=": GreaterOrEqualToken,
	"&&": AndToken,
	"||": OrToken,
}

// Keywords
const (
	Let    TokenType = "let"
	Return TokenType = "return"
	True   TokenType = "true"
	False  TokenType = "false"
	If     TokenType = "if"
	Else   TokenType = "else"
)

var keywords = map[string]Token{
	"let":    LetToken,
	"return": ReturnToken,
	"true":   TrueToken,
	"false":  FalseToken,
	"if":     IfToken,
	"else":   ElseToken,
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
	LessThanToken         = Token{Type: LessThan, Literal: "<"}
	GreaterThanToken      = Token{Type: GreaterThan, Literal: ">"}
	EqualToken            = Token{Type: Equal, Literal: "=="}
	NotEqualToken         = Token{Type: NotEqual, Literal: "!="}
	LessOrEqualToken      = Token{Type: LessOrEqual, Literal: "<="}
	GreaterOrEqualToken   = Token{Type: GreaterOrEqual, Literal: ">="}
	AndToken              = Token{Type: And, Literal: "&&"}
	OrToken               = Token{Type: Or, Literal: "||"}
	IfToken               = Token{Type: If, Literal: "if"}
	ElseToken             = Token{Type: Else, Literal: "else"}
	LeftBraceToken        = Token{Type: LeftBrace, Literal: "{"}
	RightBraceToken       = Token{Type: RightBrace, Literal: "}"}
)
