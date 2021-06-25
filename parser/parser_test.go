package parser

import (
	"fmt"
	. "mbs/common"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestReadVar(t *testing.T) {
	testCase := func(code, expectedCode, expectedName string) {
		code, expr, err := ParseReadVar(code)

		checkErrorAndCompareExpressionsAndCode(t, err, expr, ReadVar{Name: expectedName}, code, expectedCode)
	}

	testCase("abc ", " ", "abc")
	testCase("a123 ", " ", "a123")
	testCase("a123{", "{", "a123")
	testCase("a123=", "=", "a123")
	testCase(" abc = 123;", " = 123;", "abc")
	testCase(" abc = 123; b = 456;", " = 123; b = 456;", "abc")
}

func TestReadVar_negative(t *testing.T) {
	testCase := func(t *testing.T, code string) {
		_, _, err := ParseName(code)

		if err == nil {
			t.Errorf(`expected error when parsing "%s"`, code)
		}
	}

	testCase(t, "123 ")
	testCase(t, "= ")
	testCase(t, "{ ")
	testCase(t, "äzcxv")
	testCase(t, "")
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
	code := "12+34; b:=123;"
	firstExpr := Integer{Data: 12}
	secondExpr := Integer{Data: 34}
	expectedExpr := Operator{Symbol: "+", FirstExp: firstExpr, SecondExp: secondExpr}
	expectedCode := "; b:=123;"

	code, expr, err := ParseOperator(code)

	checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
}

func TestParseFunctionCall(t *testing.T) {
	// TODO: switch order of arguments
	testCase := func(code string, expectedExpr Expr, expectedCode string) {
		code, expr, err := ParseFunctionCall(code)

		checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
	}

	testCase("asdf(123); b:=123;", FunctionCall{Name: "asdf", Argument: Integer{Data: 123}}, "; b:=123;")
}

func ExampleParseFunctionCall_nested() {
	code, expr, _ := ParseFunctionCall(" a ( b ( c ( 123 ) ) ); x")
	fmt.Println("code=" + code)
	if expr != nil {
		fmt.Println("expr=" + expr.Print())
	}

	// Output:
	// code=; x
	// expr=a(b(c(123)))
}

func ExampleParseFunctionCall_complicated() {
	code, expr, _ := ParseFunctionCall(" a ( b + 123 ) )")
	fmt.Println("code=" + code)
	if expr != nil {
		fmt.Println("expr=" + expr.Print())
	}

	// Output:
	// code= )
	// expr=a(b + 123)
}

func TestParseExpression(t *testing.T) {
	testParseExpression(t, "\"Hi\"; b:=123;", String{Data: "Hi"}, "; b:=123;")
	testParseExpression(t, `""; b:=123;`, String{Data: ""}, "; b:=123;")
	testParseExpression(t, `"\""; b:=123;`, String{Data: `"`}, "; b:=123;")
	testParseExpression(t, "54.01; b:=123;", Float{Data: 54.01}, "; b:=123;")
	testParseExpression(t, "-54.01; b:=123;", Float{Data: -54.01}, "; b:=123;")
	testParseExpression(t, "987; b:=123;", Integer{Data: 987}, "; b:=123;")
	testParseExpression(t, "-987; b:=123;", Integer{Data: -987}, "; b:=123;")
	testParseExpression(t, "true; b:=123;", Boolean{Data: true}, "; b:=123;")
	testParseExpression(t, "5*2; b:=123;", Operator{Symbol: "*", FirstExp: Integer{Data: 5}, SecondExp: Integer{Data: 2}}, "; b:=123;")
	testParseExpression(t, "abc", ReadVar{Name: "abc"}, "")
	testParseExpression(t, "abc\"", ReadVar{Name: "abc"}, `"`)
	// TODO
	// testParseExpression(t, "print("\""Hello"\""), ...)
	testParseExpressionNegative(t, "*")
	testParseExpressionNegative(t, "/1")
	testParseExpressionNegative(t, "-.-.#+")
	testParseExpressionNegative(t, "\"abc")
}

func testParseExpression(t *testing.T, expression string, expectedExpression Expr, expectedCode string) {
	t.Run(expression, func(t *testing.T) {
		code, expr, err := ParseExpression(expression)
		checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpression, code, expectedCode)
	})
}

func testParseExpressionNegative(t *testing.T, expression string) {
	t.Run(expression, func(t *testing.T) {
		_, expr, err := ParseExpression(expression)
		if err == nil {
			t.Errorf(`got (%+v) wanted nil `, expr)
		}
	})
}

func checkErrorAndCompareExpressionsAndCode(t *testing.T, err error, expr Expr, expectedExpr Expr, code string, expectedCode string) {
	if err != nil {
		t.Error(err)
	}

	if !cmp.Equal(expr, expectedExpr) {
		t.Errorf(`got (Expr: "%#v") wanted (Expr: "%#v")`, expr, expectedExpr)
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

func TestParseParentheses(t *testing.T) {
	testCase := func(code string, expectedExpr Expr, expectedCode string) {
		code, expr, err := ParseParentheses(code)

		checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
	}

	testCase("(123)", Integer{Data: 123}, "")
	testCase("(123);123", Integer{Data: 123}, ";123")
	testCase("(asdf(123));123", FunctionCall{Name: "asdf", Argument: Integer{Data: 123}}, ";123")
}

func TestParseFor(t *testing.T) {
	testCase := func(code string, expectedExpr Expr, expectedCode string) {
		t.Run(code, func(t *testing.T) {
			code, expr, err := ParseFor(code)

			checkErrorAndCompareExpressionsAndCode(t, err, expr, expectedExpr, code, expectedCode)
		})
	}

	testCase("for (;false;) {}", For{
		Init:        &Nop{},
		Condition:   Boolean{Data: false},
		Advancement: &Nop{},
		Body:        Block{Statements: []Expr{}},
	}, "")
}

func ExampleParseFor() {
	_, expr, err := ParseFor(`for (e = 1; e < 4; e = e + 1) {
		print("e");
	}`)

	if err != nil {
		fmt.Println("ERROR", err)
	} else {
		fmt.Println(expr.Print())
	}

	// Output:
	// for (e = 1; e < 4; e = e + 1) {
	// print("e")
	// }
}

func ExampleParseCode_simple() {
	input := `a = 123;
b = "abc";
c = true;
d = 4.2;`

	block, err := ParseCode(input)

	fmt.Println("error:", err != nil)
	if block != nil {
		fmt.Println(block.Print())
	}

	// Output:
	// error: false
	// a = 123
	// b = "abc"
	// c = true
	// d = 4.20000
}

func ExampleParseCode_full() {
	input := `a = 123;
b = "abc";
c = true;
d = 4.2;

if (c) {
    print("c is true");
}

if (a == 123) {
	print("a is 123");
}

if (c && true) {
	print("c && true");
}

if (b == "abc") {
	print("b is abc");
}

print(b + "123");

for (;false;) {
}

for (e = 1; e < 4; e = e + 1) {
	print("e");
}

input = readline();
print(input);`

	block, err := ParseCode(input)

	fmt.Println("error:", err != nil)
	if block != nil {
		fmt.Println(block.Print())
	}

	// Output:
	// error: false
	// a = 123
	// b = "abc"
	// c = true
	// d = 4.20000
	// if (c) {
	// print("c is true")
	// }
	//
	// if (a == 123) {
	// print("a is 123")
	// }
	//
	// if (c && true) {
	// print("c && true")
	// }
	//
	// if (b == "abc") {
	// print("b is abc")
	// }
	//
	// print(b + "123")
	// for (nop; false; nop) {
	// }
	// for (e = 1; e < 4; e = e + 1) {
	// print("e")
	// }
	// input = readline(nop)
	// print(input)
}
