package eval

import (
	"spike-interpreter-go/spike/object"
)

var builtins = map[string]*object.BuiltinFunction{
	"len":   object.GetBuiltinByName("len"),
	"print": object.GetBuiltinByName("print"),
	"read":  object.GetBuiltinByName("read"),
}
