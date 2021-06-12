package parser

import (
	"bytes"
	"testing"
)

func TestParseName(t *testing.T) {
	testParseName(t, "abc ", " ", "abc")
	testParseName(t, "a123 ", " ", "a123")
	testParseName(t, "a123{", "{", "a123")
	testParseName(t, "a123=", "=", "a123")
}

func TestParseNameNegative(t *testing.T) {
	testParseNameNegative(t, "123 ")
	testParseNameNegative(t, "= ")
	testParseNameNegative(t, "{ ")
	testParseNameNegative(t, "äzcxv")
	testParseNameNegative(t, "")
}

func testParseName(t *testing.T, code, resultCode, name string) {
	xcode, xname, xerr := ParseName([]byte(code))

	if !bytes.Equal(xcode, []byte(resultCode)) || !bytes.Equal(xname, []byte(name)) || xerr != nil {
		t.Errorf(`got ("%s", "%s", %s) wanted ("%s", "%s", nil)`,
			xcode, xname, xerr,
			resultCode, name)
	}
}

func testParseNameNegative(t *testing.T, code string) {
	_, _, err := ParseName([]byte(code))

	if err == nil {
		t.Errorf(`expected error when parsing "%s"`, code)
	}
}
