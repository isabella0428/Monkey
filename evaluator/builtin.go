package evaluator

import ( 
	"../object"
)

var builtins = map[string] *object.Builtin {
	"len": &object.Builtin {
		Fn : func (args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value : int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return newError("argument to `len` not supported, got %s", args[0].Type())
			}
		},
	},
	"first": &object.Builtin {
		Fn : func (args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			return array.Elements[0]
		},
	},
	"last" : &object.Builtin {
		Fn : func (args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			num := len(array.Elements)
			if num > 0 {
				return array.Elements[num - 1]
			}
			return NULL
		},	
	},
	"rest": &object.Builtin {
		Fn : func (args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `first` must be ARRAY, got %s", args[0].Type())
			}
			
			array := args[0].(*object.Array)
			newElements := []object.Object{}

			if len(array.Elements) > 0 {
				for idx, e := range array.Elements {
					if idx != 0 {
						newElements = append(newElements, e)
					}
				}
				return &object.Array{Elements:newElements}
			}
			return NULL
		},
	},
	"push" : &object.Builtin {
		Fn : func(args...object.Object)object.Object {
			if len(args) != 2 {
				return newError("wrong number of arguments. got=%d, want=2", len(args))
			}

			if args[0].Type() != object.ARRAY_OBJ {
				return newError("argument to `push` must be ARRAY, got %s", args[0].Type())
			}

			array := args[0].(*object.Array)
			newElements := []object.Object{}

			for _, e := range array.Elements {
				newElements = append(newElements, e)
			}

			newElements = append(newElements, args[1])
			return &object.Array{Elements:newElements}
		},
	},
}

