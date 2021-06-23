package parser

type ParseError struct {
	Message string
}

func (m *ParseError) Error() string {
	return m.Message
}

func NewParseErrorExpected(expected string) *ParseError {
	return &ParseError{Message: "Expected '" + expected + "'"}
}
