package object

import (
	"fmt"
	"github.com/pkg/errors"
)

var Builtins = []*BuiltinFunction{
	{
		Name: "len",
		Function: func(args ...Object) (Object, error) {
			if len(args) != 1 {
				return nil, errors.New("1 function argument expected")
			}

			switch argument := args[0].(type) {
			case *String:
				return &Integer{Value: int64(len(argument.Value))}, nil

			case *Array:
				return &Integer{Value: int64(len(argument.Elements))}, nil
			}

			stringObject := args[0].(*String)
			return &Integer{Value: int64(len(stringObject.Value))}, nil
		},
	},
	{
		Name: "print",
		Function: func(args ...Object) (Object, error) {
			stringObject := args[0].(*String)
			fmt.Printf(stringObject.Value)

			return nil, nil
		},
	},
	{
		Name: "read",
		Function: func(args ...Object) (Object, error) {
			var result string
			_, err := fmt.Scan(&result)
			if err != nil {
				return nil, err
			}

			return &String{Value: result}, nil
		},
	},
}

func GetBuiltinByName(name string) *BuiltinFunction {
	for _, builtin := range Builtins {
		if builtin.Name == name {
			return builtin
		}
	}

	return nil
}
