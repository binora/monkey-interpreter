package evaluator

import (
	"fmt"
	"interpreters/ast"
	"interpreters/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {

	case *ast.StringLiteral:
		return &object.String{Value: node.Value}

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

	case *ast.FunctionLiteral:
		return evalFunctionLiteral(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}

		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}

		return evalIndexExpression(left, index, env)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		right := Eval(node.Right, env)
		return evalInfixExpression(node.Operator, left, right)

	case *ast.BlockStatement:
		return evalBlockStatement(node, env)

	case *ast.IfExpression:
		return evalIfExpression(node, env)

	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{Elements: elements}

	case *ast.HashLiteral:
		hash := &object.Hash{Pairs: make(map[object.HashKey]object.HashPair)}

		for key, value := range node.Pairs {
			keyObj := Eval(key, env)
			if isError(keyObj) {
				return keyObj
			}

			valueObj := Eval(value, env)
			if isError(valueObj) {
				return valueObj
			}

			hashKey, ok := keyObj.(object.Hashable)
			if !ok {
				return newError("unusable as hash key: %s", keyObj.Type())
			}
			hash.Pairs[hashKey.HashKey()] = object.HashPair{
				Key:   keyObj,
				Value: valueObj,
			}
		}
		return hash
	}

	return nil
}

func evalIndexExpression(left object.Object, index object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_OBJ:
		return evalHashIndexExpression(left, index)
	default:
		return newError("Index operator not supported: %s", left.Type())
	}
}

func evalHashIndexExpression(left object.Object, index object.Object) object.Object {
	hash := left.(*object.Hash)
	hashKeyObj, ok := index.(object.Hashable)
	if !ok {
		return newError("unusable as hash key: %s", index)
	}

	pair, ok := hash.Pairs[hashKeyObj.HashKey()]
	if !ok {
		return NULL
	}

	return pair.Value
}

func evalArrayIndexExpression(left object.Object, index object.Object) object.Object {
	array, _ := left.(*object.Array)
	i, _ := index.(*object.Integer)
	max := len(array.Elements) - 1
	if i.Value > int64(max) || i.Value < 0 {
		return NULL
	}

	return array.Elements[i.Value]
}

func applyFunction(function object.Object, args []object.Object) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return newError("not a function: %s", fn.Type())
	}
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnVal, ok := obj.(*object.ReturnValue); ok {
		return returnVal.Value
	}
	return obj
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func evalExpressions(arguments []ast.Expression, env *object.Environment) []object.Object {
	var results []object.Object

	for _, arg := range arguments {
		evaluated := Eval(arg, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		results = append(results, evaluated)
	}
	return results
}

func evalFunctionLiteral(node *ast.FunctionLiteral, env *object.Environment) object.Object {
	return &object.Function{
		Parameters: node.Parameters,
		Body:       node.Body,
		Env:        env,
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if ok {
		return val
	}
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	return newError("identifier not found: " + node.Value)
}

func isError(obj object.Object) bool {
	return obj.Type() == object.ERROR_OBJ
}

func evalBlockStatement(node *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range node.Statements {
		result = Eval(statement, env)

		if result != nil {
			if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
				return result
			}
		}
	}
	return result
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	}

	if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	}

	return NULL
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}

}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {

	if left.Type() != right.Type() {
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left, right)
	}

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}
	switch operator {
	case "==":
		return nativeBoolToBooleanObject(left == right)
	case "!=":
		return nativeBoolToBooleanObject(left != right)
	}
	return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
}

func evalStringInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftObj, rightObj := left.(*object.String), right.(*object.String)
	switch operator {
	case "+":
		return &object.String{Value: leftObj.Value + rightObj.Value}
	default:
		return newError(fmt.Sprintf("operator not supported: %s %s %s ", leftObj.Type(), operator, rightObj.Type()))
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftObject, rightObject := left.(*object.Integer), right.(*object.Integer)
	switch operator {
	case "-":
		return &object.Integer{Value: leftObject.Value - rightObject.Value}
	case "+":
		return &object.Integer{Value: leftObject.Value + rightObject.Value}
	case "/":
		return &object.Integer{Value: leftObject.Value / rightObject.Value}
	case "*":
		return &object.Integer{Value: leftObject.Value * rightObject.Value}
	case ">":
		return nativeBoolToBooleanObject(leftObject.Value > rightObject.Value)
	case "<":
		return nativeBoolToBooleanObject(leftObject.Value < rightObject.Value)
	case "==":
		return nativeBoolToBooleanObject(leftObject.Value == rightObject.Value)
	case "!=":
		return nativeBoolToBooleanObject(leftObject.Value != rightObject.Value)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func nativeBoolToBooleanObject(input bool) object.Object {
	if input {
		return TRUE
	}
	return FALSE
}

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperator(right)
	case "-":
		return evalMinusPrefixOperator(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperator(operand object.Object) object.Object {
	switch operand {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return FALSE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperator(operand object.Object) object.Object {
	if operand.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", operand.Type())
	}
	obj := operand.(*object.Integer)

	return &object.Integer{Value: -obj.Value}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
