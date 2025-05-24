package main

// Eval evaluates an AST node and returns the resulting object
func Eval(node Node, env *Environment) Object {
	switch node := node.(type) {

	// Statements
	case *Program:
		return evalProgram(node.Statements, env)

	case *ExpressionStatement:
		return Eval(node.Expression, env)

	case *BlockStatement:
		return evalBlockStatement(node, env)

	case *LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Set(node.Name.Value, val)
		return NULL

	case *ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &ReturnValue{Value: val}

	case *FunctionStatement:
		fn := &Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
		env.Set(node.Name.Value, fn)
		return NULL

	case *IfStatement:
		return evalIfExpression(node, env)

	// Expressions
	case *IntegerLiteral:
		return &Integer{Value: node.Value}

	case *StringLiteral:
		return &String{Value: node.Value}

	case *BooleanLiteral:
		return nativeBoolToPyBoolean(node.Value)

	case *Identifier:
		return evalIdentifier(node, env)

	case *PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)

	default:
		return newError("unknown node type: %T", node)
	}
}

// evalProgram evaluates a sequence of statements
func evalProgram(stmts []Statement, env *Environment) Object {
	var result Object

	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *ReturnValue:
			return result.Value
		case *Error:
			return result
		}
	}

	return result
}

// evalBlockStatement evaluates a block of statements
func evalBlockStatement(block *BlockStatement, env *Environment) Object {
	var result Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == RETURN_OBJ || rt == ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

// evalIfExpression evaluates if-else expressions
func evalIfExpression(ie *IfStatement, env *Environment) Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

// evalIdentifier evaluates identifier expressions
func evalIdentifier(node *Identifier, env *Environment) Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: " + node.Value)
	}
	return val
}

// evalPrefixExpression evaluates prefix expressions like !x or -x
func evalPrefixExpression(operator string, right Object) Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression evaluates the ! operator
func evalBangOperatorExpression(right Object) Object {
	switch right {
	case RUNTIME_TRUE:
		return RUNTIME_FALSE
	case RUNTIME_FALSE:
		return RUNTIME_TRUE
	case NULL:
		return RUNTIME_TRUE
	default:
		return RUNTIME_FALSE
	}
}

// evalMinusPrefixOperatorExpression evaluates the - prefix operator
func evalMinusPrefixOperatorExpression(right Object) Object {
	if right.Type() != INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*Integer).Value
	return &Integer{Value: -value}
}

// evalInfixExpression evaluates infix expressions like x + y
func evalInfixExpression(operator string, left, right Object) Object {
	switch {
	case left.Type() == INTEGER_OBJ && right.Type() == INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == STRING_OBJ && right.Type() == STRING_OBJ:
		return evalStringInfixExpression(operator, left, right)
	case operator == "==":
		return nativeBoolToPyBoolean(left == right)
	case operator == "!=":
		return nativeBoolToPyBoolean(left != right)
	case operator == "&&":
		return evalLogicalAndExpression(left, right)
	case operator == "||":
		return evalLogicalOrExpression(left, right)
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalIntegerInfixExpression evaluates integer infix expressions
func evalIntegerInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*Integer).Value
	rightVal := right.(*Integer).Value

	switch operator {
	case "+":
		return &Integer{Value: leftVal + rightVal}
	case "-":
		return &Integer{Value: leftVal - rightVal}
	case "*":
		return &Integer{Value: leftVal * rightVal}
	case "/":
		if rightVal == 0 {
			return newError("division by zero")
		}
		return &Integer{Value: leftVal / rightVal}
	case "<":
		return nativeBoolToPyBoolean(leftVal < rightVal)
	case ">":
		return nativeBoolToPyBoolean(leftVal > rightVal)
	case "<=":
		return nativeBoolToPyBoolean(leftVal <= rightVal)
	case ">=":
		return nativeBoolToPyBoolean(leftVal >= rightVal)
	case "==":
		return nativeBoolToPyBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToPyBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s", operator)
	}
}

// evalStringInfixExpression evaluates string infix expressions
func evalStringInfixExpression(operator string, left, right Object) Object {
	leftVal := left.(*String).Value
	rightVal := right.(*String).Value

	switch operator {
	case "+":
		return &String{Value: leftVal + rightVal}
	case "==":
		return nativeBoolToPyBoolean(leftVal == rightVal)
	case "!=":
		return nativeBoolToPyBoolean(leftVal != rightVal)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalLogicalAndExpression evaluates && expressions with short-circuit evaluation
func evalLogicalAndExpression(left, right Object) Object {
	if !isTruthy(left) {
		return RUNTIME_FALSE
	}
	return nativeBoolToPyBoolean(isTruthy(right))
}

// evalLogicalOrExpression evaluates || expressions with short-circuit evaluation
func evalLogicalOrExpression(left, right Object) Object {
	if isTruthy(left) {
		return RUNTIME_TRUE
	}
	return nativeBoolToPyBoolean(isTruthy(right))
}

// evalExpressions evaluates a list of expressions
func evalExpressions(exps []Expression, env *Environment) []Object {
	var result []Object

	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

// applyFunction applies a function to its arguments
func applyFunction(fn Object, args []Object) Object {
	switch fn := fn.(type) {
	case *Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	default:
		return newError("not a function: %T", fn)
	}
}

// extendFunctionEnv creates a new environment for function execution
func extendFunctionEnv(fn *Function, args []Object) *Environment {
	env := NewEnclosedEnvironment(fn.Env)

	for paramIdx, param := range fn.Parameters {
		if paramIdx >= len(args) {
			env.Set(param.Value, NULL)
		} else {
			env.Set(param.Value, args[paramIdx])
		}
	}

	return env
}

// unwrapReturnValue unwraps return values
func unwrapReturnValue(obj Object) Object {
	if returnValue, ok := obj.(*ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

// Helper functions

// nativeBoolToPyBoolean converts a Go boolean to a TinyLang boolean object
func nativeBoolToPyBoolean(input bool) *Boolean {
	if input {
		return RUNTIME_TRUE
	}
	return RUNTIME_FALSE
}

// isError checks if an object is an error
func isError(obj Object) bool {
	if obj != nil {
		return obj.Type() == ERROR_OBJ
	}
	return false
}
