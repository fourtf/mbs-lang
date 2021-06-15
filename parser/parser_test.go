package parser

import (
	. "mbs/common"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseName(t *testing.T) {
	testParseName(t, "abc ", "", "abc")
	testParseName(t, "a123 ", "", "a123")
	testParseName(t, "a123{", "{", "a123")
	testParseName(t, "a123=", "=", "a123")
	testParseName(t, " abc = 123;", "=123;", "abc")
	testParseName(t, " abc = 123; b = 456;", "=123;b=456;", "abc")
}

func TestParseNameNegative(t *testing.T) {
	testParseNameNegative(t, "123 ")
	testParseNameNegative(t, "= ")
	testParseNameNegative(t, "{ ")
	testParseNameNegative(t, "äzcxv")
	testParseNameNegative(t, "")
}

func testParseName(t *testing.T, code, resultCode, name string) {
	xcode, xname, xerr := ParseName(code)

	if xcode != resultCode || xname != name || xerr != nil {
		t.Errorf(`got ("%s", "%s", %s) wanted ("%s", "%s", nil)`,
			xcode, xname, xerr,
			resultCode, name)
	}
}

func testParseNameNegative(t *testing.T, code string) {
	_, _, err := ParseName(code)

	if err == nil {
		t.Errorf(`expected error when parsing "%s"`, code)
	}
}

func TestParseString(t *testing.T) {
	expression := "\"Hello World\""
	expectedExpr := String{Data: "Hello World"}

	expr, err := ParseString(expression)

	checkErrorAndCompareExpressions(t, err, expr, expectedExpr)
}

func TestParseBoolean(t *testing.T) {
	expression := "false"
	expectedExpr := Boolean{Data: false}

	expr, err := ParseBoolean(expression)

	checkErrorAndCompareExpressions(t, err, expr, expectedExpr)
}

func TestParseInteger(t *testing.T) {
	expression := "12345"
	expectedExpr := Integer{Data: 12345}

	expr, err := ParseInteger(expression)

	checkErrorAndCompareExpressions(t, err, expr, expectedExpr)
}

func TestParseFloat(t *testing.T) {
	expression := "123.51"
	expectedExpr := Float{Data: 123.51}

	expr, err := ParseFloat(expression)

	checkErrorAndCompareExpressions(t, err, expr, expectedExpr)
}

func TestParseOperator(t *testing.T) {
	// Doesn't work if ParseExpression doesn't work
	expression := "12+34"
	firstExpr := Integer{Data: 12}
	secondExpr := Integer{Data: 34}
	expectedExpr := Operator{Symbol: "+", FirstExp: firstExpr, SecondExp: secondExpr}

	expr, err := ParseOperator(expression)

	checkErrorAndCompareExpressions(t, err, expr, expectedExpr)
}

func TestParseFunction(t *testing.T) {
	// TODO
}

func TestParseExpression(t *testing.T) {
	testParseExpression(t, "\"Hi\"", String{Data: "Hi"})
	testParseExpression(t, "\"\"", String{Data: ""})
	testParseExpression(t, "54.01", Float{Data: 54.01})
	testParseExpression(t, "54.", Float{Data: 54.0})
	testParseExpression(t, "987", Integer{Data: 987})
	testParseExpression(t, "true", Boolean{Data: true})
	testParseExpression(t, "5*2", Operator{Symbol: "*", FirstExp: Integer{Data: 5}, SecondExp: Integer{Data: 2}})
	// TODO
	// testParseExpression(t, "print("\""Hello"\""), ...)
	testParseExpressionNegative(t, "*")
	testParseExpressionNegative(t, "5+")
	testParseExpressionNegative(t, "/1")
	testParseExpressionNegative(t, "-.-.#+")
	testParseExpressionNegative(t, "\"abc")
	testParseExpressionNegative(t, "abc\"")
	testParseExpressionNegative(t, "abc")
}

func testParseExpression(t *testing.T, expression string, expectedExpression Expr) {
	expr, err := ParseExpression(expression)
	checkErrorAndCompareExpressions(t, err, expr, expectedExpression)
}

func testParseExpressionNegative(t *testing.T, expression string) {
	expr, err := ParseExpression(expression)
	if err == nil {
		t.Errorf(`got (%+v) wanted nil `, expr)
	}
}

func checkErrorAndCompareExpressions(t *testing.T, err error, expr Expr, expectedExpr Expr) {
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(expr, expectedExpr) {
		t.Errorf(`got (Expr: "%+v") wanted (Expr: "%+v")`, expr, expectedExpr)
	}
}
func TestParseWriteVar(t *testing.T) {
	expectedName := "a"
	expectedCode := "b=456;c=546;"
	expectedExpr := Integer{Data: 123}

	code, expr, err := ParseWriteVar(" a = 123 ; b = 456  ;  \n\r c = 546;")

	if code != expectedCode || expr == nil || err != nil {
		t.Errorf(`got (Code: "%s", Expr: "%s", Err: %s) wanted ("%s", "%+v", nil)`, code, expr, err, expectedCode, expectedExpr)
	}

	if writeVar, ok := expr.(WriteVar); ok {
		if writeVar.Name != expectedName || writeVar.Expr == nil {
			t.Errorf(`got (Name: "%s", Expr: nil) wanted (Name: "%s", Expr: "%+v")`, writeVar.Name, expectedName, expectedExpr)
		}
		if !cmp.Equal(writeVar.Expr, expectedExpr) {
			t.Errorf(`got (Expr: "%s") wanted (Expr: "%+v")`, writeVar.Expr, expectedExpr)
		}
	} else {
		t.Errorf("The expression is not of type WriteVar!")
	}
}
