package evaluator

import (
	"fmt"
	"interpreters/object"
)

var builtins = map[string]*object.Builtin{
	"len": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError(fmt.Sprintf("wrong number of arguments. got=%d want=%d", len(args), 1))
			}

			switch args[0].Type() {
			case object.INTEGER_OBJ:
				return newError(fmt.Sprintf("argument to `len` not supported, got %s", args[0].Type()))
			case object.STRING_OBJ:
				literal := args[0].(*object.String)
				return &object.Integer{Value: int64(len(literal.Value))}
			}

			return NULL
		},
	},
	"puts": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
