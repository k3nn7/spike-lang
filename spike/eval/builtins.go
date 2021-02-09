package eval

import (
	"errors"
	"fmt"
	"spike-interpreter-go/spike/object"
)

var builtins = map[string]*object.BuiltinFunction{
	"len": {
		Name: "len",
		Function: func(args ...object.Object) (object.Object, error) {
			if len(args) != 1 {
				return nil, errors.New("1 function argument expected")
			}

			switch argument := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(argument.Value))}, nil

			case *object.Array:
				return &object.Integer{Value: int64(len(argument.Elements))}, nil
			}

			stringObject := args[0].(*object.String)
			return &object.Integer{Value: int64(len(stringObject.Value))}, nil
		},
	},
	"print": {
		Function: func(args ...object.Object) (object.Object, error) {
			stringObject := args[0].(*object.String)
			fmt.Printf(stringObject.Value)

			return nil, nil
		},
	},
	"read": {
		Name: "read",
		Function: func(args ...object.Object) (object.Object, error) {
			var result string
			_, err := fmt.Scan(&result)
			if err != nil {
				return nil, err
			}

			return &object.String{Value: result}, nil
		},
	},
}
