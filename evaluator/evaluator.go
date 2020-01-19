package evaluator

import(
	"../ast"
	"../object"
)

var (
	TRUE = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL = &object.Null{}

)

func Eval (node ast.Node) object.Object {
	switch node := node.(type) {
		case *ast.Program:
			return evalStatements(node.Statements)
		case *ast.ExpressionStatement:
			return Eval(node.Expression)
		// Expression
		case *ast.IntegerLiteral:
			return &object.Integer{Value:node.Value}
		case *ast.Boolean:
			return nativeBoolToBooleanObject(node.Value)
		case *ast.PrefixExpression:
			right := Eval(node.Right)
			return evalPrefixExpression(node.Operator, right)
		case *ast.InfixExpression:
			left := Eval(node.Left)
			right := Eval(node.Right)
			return evalInfixExpression(node.Operator,left, right)
		case *ast.IfExpression:
			return evalIfExpression(node)
		case *ast.BlockStatement:
			return evalStatements(node.Statements)

	}
	return NULL
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
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
			return NULL
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
		return NULL 
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
		// We can write so since we only have two boolean object here
		case operator == "==":
			return nativeBoolToBooleanObject(left == right)
		case operator == "!=":
			return nativeBoolToBooleanObject(left != right)
		default:
			return NULL
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
			return NULL
	}
}

// Avoid Creating multiple TRUE or FALSE boolean objects
func nativeBoolToBooleanObject(Value bool) object.Object{
	if Value {
		return TRUE
	}
	return FALSE
}

func evalIfExpression(ie *ast.IfExpression) object.Object {
	condition := Eval(ie.Condition)
	if isTruthy(condition) {
		return Eval(ie.Consequence)
	} else if ie.Alternative != nil{ 
		return Eval(ie.Alternative)
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

