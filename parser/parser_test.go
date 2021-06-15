package parser

import (
	. "mbs/common"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseName(t *testing.T) {
	testParseName(t, "abc ", " ", "abc")
	testParseName(t, "a123 ", " ", "a123")
	testParseName(t, "a123{", "{", "a123")
	testParseName(t, "a123=", "=", "a123")
	testParseName(t, " abc = 123;", " = 123;", "abc")
	testParseName(t, " abc = 123; b = 456;", " = 123; b = 456;", "abc")
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
	code := "\"Hello World\"; b:=123;"
	expectedExpr := String{Data: "Hello World"}
	expectedCode := "; b:=123;"

	code, expr, err := ParseString(code)

	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
}

func TestParseBoolean(t *testing.T) {
	code := "false; b:=123;"
	expectedExpr := Boolean{Data: false}
	expectedCode := "; b:=123;"

	code, expr, err := ParseBoolean(code)

	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
}

func TestParseInteger(t *testing.T) {
	code := "12345; b:=123;"
	expectedExpr := Integer{Data: 12345}
	expectedCode := "; b:=123;"

	code, expr, err := ParseInteger(code)

	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
}

func TestParseFloat(t *testing.T) {
	code := "123.51; b:=123;"
	expectedExpr := Float{Data: 123.51}
	expectedCode := "; b:=123;"

	code, expr, err := ParseFloat(code)

	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
}

func TestParseOperator(t *testing.T) {
	// Doesn't work if ParseExpression doesn't work
	code := "12+34; b:=123;"
	firstExpr := Integer{Data: 12}
	secondExpr := Integer{Data: 34}
	expectedExpr := Operator{Symbol: "+", FirstExp: firstExpr, SecondExp: secondExpr}
	expectedCode := "; b:=123;"

	code, expr, err := ParseOperator(code)

	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
}

func TestParseFunction(t *testing.T) {
	// TODO
}

func TestParseExpression(t *testing.T) {
	testParseExpression(t, "\"Hi\"; b:=123;", String{Data: "Hi"}, "; b:=123;")
	testParseExpression(t, "\"\"; b:=123;", String{Data: ""}, "; b:=123;")
	testParseExpression(t, "54.01; b:=123;", Float{Data: 54.01}, "; b:=123;")
	testParseExpression(t, "-54.01; b:=123;", Float{Data: -54.01}, "; b:=123;")
	testParseExpression(t, "987; b:=123;", Integer{Data: 987}, "; b:=123;")
	testParseExpression(t, "-987; b:=123;", Integer{Data: -987}, "; b:=123;")
	testParseExpression(t, "true; b:=123;", Boolean{Data: true}, "; b:=123;")
	testParseExpression(t, "5*2; b:=123;", Operator{Symbol: "*", FirstExp: Integer{Data: 5}, SecondExp: Integer{Data: 2}}, "; b:=123;")
	// TODO
	// testParseExpression(t, "print("\""Hello"\""), ...)
	testParseExpressionNegative(t, "*")
	testParseExpressionNegative(t, "/1")
	testParseExpressionNegative(t, "-.-.#+")
	testParseExpressionNegative(t, "\"abc")
	testParseExpressionNegative(t, "abc\"")
	testParseExpressionNegative(t, "abc")
}

func testParseExpression(t *testing.T, expression string, expectedExpression Expr, expectedCode string) {
	code, expr, err := ParseExpression(expression)
	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpression, code, expectedCode)
}

func testParseExpressionNegative(t *testing.T, expression string) {
	_, expr, err := ParseExpression(expression)
	if err == nil {
		t.Errorf(`got (%+v) wanted nil `, expr)
	}
}

func checkErrorAndCompareExpressionsAndCode(t *testing.T, err error, expr Expr, expectedExpr Expr, code string, expectedCode string) {
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(expr, expectedExpr) {
		t.Errorf(`got (Expr: "%+v") wanted (Expr: "%+v")`, expr, expectedExpr)
	}

	if code != expectedCode {
		t.Errorf(`got (Code: "%s") wanted (Code: "%s")`, code, expectedCode)
	}
}
func TestParseWriteVar(t *testing.T) {
	expectedName := "a"
	expectedCode := " ; b = 456  ;  \n\r c = 546;"
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
