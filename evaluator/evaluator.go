package evaluator

import(
	"fmt"
	"../ast"
	"../object"
)

var (
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL = &object.Null{}
)

func Eval (node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
		case *ast.Program:
			return evalProgram(node, env)
		case *ast.ExpressionStatement:
			return Eval(node.Expression, env)
		// Expression
		case *ast.IntegerLiteral:
			return &object.Integer{Value:node.Value}
		case *ast.Boolean:
			return nativeBoolToBooleanObject(node.Value)
		case *ast.PrefixExpression:
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalPrefixExpression(node.Operator, right)
		case *ast.InfixExpression:
			left := Eval(node.Left, env)
			if isError(left) {
				return left
			}
			right := Eval(node.Right, env)
			if isError(right) {
				return right
			}
			return evalInfixExpression(node.Operator,left, right)
		case *ast.IfExpression:
			return evalIfExpression(node, env)
		case *ast.BlockStatement:
			return evalBlockStatement(node, env)
		case *ast.ReturnStatement:
			val := Eval(node.ReturnValue, env)
			if isError(val) {
				return val
			}
			return &object.ReturnValue{Value: val}
		case *ast.LetStatement: 
			val := Eval(node.Value, env)
			if isError(val) {
				return val
			}
			env.Set(node.Name.Value, val)
		case *ast.Identifier:
			return evalIdentifier(node, env)
		case *ast.StringLiteral:
			return &object.String{Value:node.Value}
		case *ast.FunctionLiteral:
			params  := node.Parameters
			body 	:= node.Body
			return &object.Function{Parameters:params, Body:body, Env:env}
		case *ast.CallExpression:
			function := Eval(node.Function, env)
			if isError(function) {
				return function
			}

			args := evalExpressions(node.Arguments, env)
			if len(args) == 1 && isError(args[0]) {
				return args[0]
			}
			return applyFunction(function, args)
	}
	return NULL
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch result := result.(type) {
			case *object.ReturnValue:
				return result.Value
			case *object.Error:
				return result
		}

		// It finally get unwrap here
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}


	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object{
	switch operator {
		case "!":
			return evalBangOperatorExpression(right)
		case "-":
			return evalMinusPrefixOperatorExpression(right)
		default:
			return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(Obj object.Object) object.Object {
	switch Obj {
		case TRUE:
			return FALSE
		case FALSE:
			return TRUE
		case NULL:
			return TRUE
		default:
			return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	result := right.(*object.Integer)
	return &object.Integer{Value:-result.Value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object{
	switch {
		// Direct comparison is not suitable for integer comparison 
		// since we always allocate new instance for integers
		case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
			return evalIntegerInfixExpression(operator, left, right)
		case left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ:
			return evalStringInfixExpression(operator, left, right)
		// We can write so since we only have two boolean object here
		case operator == "==":
			return nativeBoolToBooleanObject(left == right)
		case operator == "!=":
			return nativeBoolToBooleanObject(left != right)
		case left.Type() != right.Type():
			return newError("type mismatch: %s %s %s",
				left.Type(), operator, right.Type())
		default:
			return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object{
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
		case "+":
			return &object.Integer{Value:leftValue + rightValue}
		case "-":
			return &object.Integer{Value:leftValue - rightValue}
		case "*":
			return &object.Integer{Value:leftValue * rightValue}
		case "/":
			return &object.Integer{Value:leftValue / rightValue}
		case ">":
			return nativeBoolToBooleanObject(leftValue > rightValue)
		case "<":
			return nativeBoolToBooleanObject(leftValue < rightValue)
		case "==":
			return nativeBoolToBooleanObject(leftValue == rightValue)
		case "!=":
			return nativeBoolToBooleanObject(leftValue != rightValue)
		default:
			return newError("unknown operator: %s %s %s",
               left.Type(), operator, right.Type())
	}
}

// Avoid Creating multiple TRUE or FALSE boolean objects
func nativeBoolToBooleanObject(Value bool) object.Object{
	if Value {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition 
	}
	if isTruthy(condition) {
		return Eval(ie.Consequence,env)
	} else if ie.Alternative != nil{ 
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func isTruthy(obj object.Object) bool{
	switch obj {
		case FALSE:
			return false
		case TRUE:
			return true
		case NULL:
			return false
		default:
			return true
	}
}

// ResultValue cannot be unwrapm due to nested statements
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURN_VALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func newError(format string, a ...interface{}) object.Object {
	return &object.Error{Message:fmt.Sprintf(format, a...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if val, ok := env.Get(node.Value); ok {
		return val
	}

	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}

	return newError("identifier not found: " + node.Value)
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, exp := range exps {
		evaluated := Eval(exp, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}
	
	return result
}

func applyFunction(fn object.Object, params []object.Object) object.Object {
	switch fn := fn.(type) {
		case *object.Function:
			extendedEnv := extendedFunctionEnv(fn, params)
			evaluated := Eval(fn.Body, extendedEnv)
			return unwrapReturnValue(evaluated)
		case *object.Builtin:
			return fn.Fn(params...)
		default:
			return newError("not a function: %s", fn.Type())
	}
}

func extendedFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value 
	}
	return obj
}

func evalStringInfixExpression(operator string, left,right object.Object) object.Object{
	if operator != "+" {
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value
	return &object.String{Value:leftValue + rightValue}
}