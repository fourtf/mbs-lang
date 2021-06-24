package typechecker

import (
	. "mbs/common"
	"testing"
)

func TestTypeCheckExpr(t *testing.T) {
	testTypeCheckExpr(t, WriteVar{Name: "abc", Expr: FunctionCall{Name: "readln", Argument: Nop{}}})
	testTypeCheckExpr(t, If{
		Condition: Operator{Symbol: "==", FirstExp: ReadVar{Name: "abc"}, SecondExp: String{Data: "abc"}},
		Body: &Block{Statements: []Expr{
			FunctionCall{Name: "println", Argument: String{Data: "Equal!"}},
			WriteVar{Name: "def", Expr: Integer{Data: 5}},
			For{
				Init:        WriteVar{Name: "i", Expr: Integer{Data: 0}},
				Condition:   Operator{Symbol: "<", FirstExp: ReadVar{Name: "i"}, SecondExp: ReadVar{Name: "def"}},
				Advancement: WriteVar{Name: "i", Expr: Operator{Symbol: "+", FirstExp: ReadVar{Name: "i"}, SecondExp: Integer{Data: 1}}},
				Body: &Block{Statements: []Expr{
					FunctionCall{Name: "println", Argument: String{Data: "iterating"}},
				}},
			}}},
	})
}

func testTypeCheckExpr(t *testing.T, expr Expr) {
	typesValid := TypeCheckExpr(expr)

	if !typesValid {
		t.Errorf(`Types are invalid at expression: "%+v" but should be valid`, expr)
	}
}

func testTypeCheckExprNegative(t *testing.T, expr Expr) {
	typesValid := TypeCheckExpr(expr)

	if typesValid {
		t.Errorf(`Types are valid at expression: "%+v" but should be invalid`, expr)
	}
}

func TestTypeCheckOperator(t *testing.T) {
	testTypeCheckOperator(t, Operator{Symbol: "==", FirstExp: Integer{Data: 1}, SecondExp: Integer{Data: 1}}, BooleanType)
	testTypeCheckOperator(t, Operator{Symbol: ">=", FirstExp: Integer{Data: 1}, SecondExp: Integer{Data: 1}}, BooleanType)
	testTypeCheckOperator(t, Operator{Symbol: "*", FirstExp: Integer{Data: 1}, SecondExp: Integer{Data: 1}}, IntegerType)
	testTypeCheckOperator(t, Operator{Symbol: "+", FirstExp: Float{Data: 1.0}, SecondExp: Integer{Data: 1}}, FloatType)
	testTypeCheckOperator(t, Operator{Symbol: "+", FirstExp: String{Data: "ab"}, SecondExp: String{Data: "cd"}}, StringType)
	testTypeCheckOperator(t, Operator{Symbol: "||", FirstExp: Boolean{Data: true}, SecondExp: Boolean{Data: true}}, BooleanType)
	testTypeCheckOperator(t, Operator{
		Symbol:    "==",
		FirstExp:  Integer{Data: 1},
		SecondExp: Operator{Symbol: "-", FirstExp: Integer{Data: 2}, SecondExp: Integer{Data: 1}}}, BooleanType)

	testTypeCheckOperatorNegative(t, Operator{Symbol: "<=", FirstExp: Boolean{Data: false}, SecondExp: Boolean{Data: true}})
	testTypeCheckOperatorNegative(t, Operator{Symbol: "!=", FirstExp: Float{Data: 1.0}, SecondExp: Integer{Data: 1}})
	testTypeCheckOperatorNegative(t, Operator{Symbol: "&&", FirstExp: Integer{Data: 1}, SecondExp: Integer{Data: 1}})
	testTypeCheckOperatorNegative(t, Operator{Symbol: "+", FirstExp: String{Data: "1"}, SecondExp: Integer{Data: 1}})
	testTypeCheckOperatorNegative(t, Operator{Symbol: "-", FirstExp: Boolean{Data: false}, SecondExp: Boolean{Data: true}})
}

func testTypeCheckOperator(t *testing.T, operator Operator, expectedType Type) {
	tipe := TypeCheckOperator(operator)

	if tipe != expectedType {
		t.Errorf(`expected type "%v" but got type "%v" after input of "%+v"`, expectedType, tipe, operator)
	}
}

func testTypeCheckOperatorNegative(t *testing.T, operator Operator) {
	tipe := TypeCheckOperator(operator)

	if tipe != NopType {
		t.Errorf(`expected type Nop but got type "%v" after input of "%+v"`, tipe, operator)
	}
}

func TestTypeCheckFunctionCall(t *testing.T) {
	testTypeCheckFunctionCall(t, FunctionCall{Name: "println", Argument: String{Data: "Hello World"}}, NopType)
	testTypeCheckFunctionCall(t, FunctionCall{Name: "readln", Argument: Nop{}}, StringType)

	testTypeCheckFunctionCallNegative(t, FunctionCall{Name: "readln", Argument: String{Data: "ABC"}})
	testTypeCheckFunctionCallNegative(t, FunctionCall{Name: "println", Argument: Nop{}})
	testTypeCheckFunctionCallNegative(t, FunctionCall{Name: "erfunden", Argument: String{Data: "ABC"}})
}

func testTypeCheckFunctionCall(t *testing.T, function FunctionCall, expectedType Type) {
	valid, tipe := TypeCheckFunctionCall(function)

	if !valid || tipe != expectedType {
		t.Errorf(`expected valid: "true", type: "%v" but got valid: "%v", type: "%v" after input of "%+v"`, expectedType, valid, tipe, function)
	}
}

func testTypeCheckFunctionCallNegative(t *testing.T, function FunctionCall) {
	valid, tipe := TypeCheckFunctionCall(function)

	if valid || tipe != NopType {
		t.Errorf(`expected valid: "false", type: Nop but got valid: "%v", type: "%v" after input of "%+v"`, valid, tipe, function)
	}
}
