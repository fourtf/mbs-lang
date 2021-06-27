package typechecker

import (
	. "mbs/common"
)

/*This typechecker validates the type-safety of every expression in our AST*/

// stores all the declared variables and their type to look up their type when they are used later on in the code
var variables map[string]Type = make(map[string]Type)

func TypeCheckBlock(block *Block) bool {
	outerScopeVars := make(map[string]Type) // holds the variables declared outside of the current block
	for k, v := range variables {
		outerScopeVars[k] = v
	}
	// type-checking every expression inside of the current block
	for _, expr := range block.Statements {
		typesValid := TypeCheckExpr(expr)
		if !typesValid {
			return false
		}
	}
	variables = outerScopeVars // resets the variables to the ones declared outside of the  current block to "delete" the vars declared inside of it
	return true
}

// type-checking of expressions that can occur outside of another expression
func TypeCheckExpr(expr Expr) bool {
	switch exprType := expr.Type(); exprType {
	case WriteVarType:
		return TypeCheckWriteVar(expr.(WriteVar))
	case IfType:
		return TypeCheckIf(expr.(If))
	case ForType:
		return TypeCheckFor(expr.(For))
	case FunctionCallType:
		valid, _ := TypeCheckFunctionCall(expr.(FunctionCall))
		return valid
	}
	return false
}

// type-checking of expressions that can occur inside of another expression
func TypeCheckRightExpr(expr Expr) Type {
	switch exprType := expr.Type(); exprType {
	case OperatorType:
		return TypeCheckOperator(expr.(Operator))
	case FunctionCallType:
		valid, returnType := TypeCheckFunctionCall(expr.(FunctionCall))
		if !valid {
			return NopType
		}
		return returnType
	case ReadVarType:
		return TypeCheckReadVar(expr.(ReadVar))
	case IntegerType, FloatType, BooleanType, StringType:
		return exprType
	}
	return NopType
}

var (
	typeEqualCompOps = []string{"==", "!="}
	boolCompOps      = []string{"&&", "||"}
	arithmCompOps    = []string{">", "<", ">=", "<="}
	arithmOps        = []string{"+", "-", "*", "/"}
)

func TypeCheckOperator(operator Operator) Type {
	// checking the type of the expressions left and right of our operator
	firstExpType := TypeCheckRightExpr(operator.FirstExp)
	secondExpType := TypeCheckRightExpr(operator.SecondExp)

	// checking if the types can be used with the given operator
	for _, symbol := range typeEqualCompOps {
		if symbol == operator.Symbol && firstExpType == secondExpType {
			return BooleanType
		}
	}

	for _, symbol := range boolCompOps {
		if symbol == operator.Symbol && firstExpType == BooleanType && secondExpType == BooleanType {
			return BooleanType
		}
	}

	for _, symbol := range arithmCompOps {
		if symbol == operator.Symbol && (firstExpType == IntegerType || firstExpType == FloatType) && (secondExpType == IntegerType || secondExpType == FloatType) {
			return BooleanType
		}
	}

	for _, symbol := range arithmOps {
		if symbol == operator.Symbol {
			if firstExpType == FloatType && (secondExpType == FloatType || secondExpType == IntegerType) {
				return FloatType
			}
			if secondExpType == FloatType && (firstExpType == FloatType || firstExpType == IntegerType) {
				return FloatType
			}
			if firstExpType == IntegerType && secondExpType == IntegerType {
				return IntegerType
			}
		}
	}

	if operator.Symbol == "+" && firstExpType == StringType && secondExpType == StringType {
		return StringType
	}
	return NopType
}

// returns if the types in this function call are valid and what type is returned by the function
func TypeCheckFunctionCall(function FunctionCall) (bool, Type) {
	if function.Name == "println" && TypeCheckRightExpr(function.Argument) == StringType {
		return true, NopType
	} else if function.Name == "readln" && function.Argument.Type() == NopType {
		return true, StringType
	}
	return false, NopType
}

func TypeCheckWriteVar(writeVar WriteVar) bool {
	exprType := TypeCheckRightExpr(writeVar.Expr)
	if exprType == NopType {
		return false
	}
	variables[writeVar.Name] = exprType
	return true
}

func TypeCheckIf(ifExpr If) bool {
	if !TypeCheckCondition(ifExpr.Condition) {
		return false
	}

	return TypeCheckBlock(&ifExpr.Body)
}

func TypeCheckFor(forExpr For) bool {
	initType := forExpr.Init.Type()
	if initType == WriteVarType {
		if !TypeCheckWriteVar(forExpr.Init.(WriteVar)) {
			return false
		}
	} else if initType != NopType {
		return false
	}

	if !TypeCheckCondition(forExpr.Condition) {
		if forExpr.Condition.Type() != NopType {
			return false
		}
	}

	advType := forExpr.Advancement.Type()
	if advType == WriteVarType {
		if !TypeCheckWriteVar(forExpr.Advancement.(WriteVar)) {
			return false
		}
	} else if advType != NopType {
		return false
	}

	return TypeCheckBlock(&forExpr.Body)
}

func TypeCheckCondition(expr Expr) bool {
	switch exprType := expr.Type(); exprType {
	case BooleanType:
		return true
	case OperatorType:
		return TypeCheckOperator(expr.(Operator)) == BooleanType
	case ReadVarType:
		return TypeCheckReadVar(expr.(ReadVar)) == BooleanType
	case FunctionCallType:
		valid, returnType := TypeCheckFunctionCall(expr.(FunctionCall))
		return valid && returnType == BooleanType
	}

	return false
}

func TypeCheckReadVar(readVar ReadVar) Type {
	if tipe, ok := variables[readVar.Name]; ok {
		return tipe
	}
	return NopType
}
