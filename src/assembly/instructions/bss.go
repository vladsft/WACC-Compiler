package instructions

var (
	_ Instruction = StringLiteral{}
)

// Structure for the Block Starting Symbol
// This contains statically-allocated variables that are declared but they are not assigned a value yet.
type bss struct {
	ID   string
	Size int
}

type StringLiteral struct {
	bss
	String string
}

// Set of all escape characters
var escapeCharSet = map[rune]struct{}{
	'\\': {},
	'0':  {},
	'b':  {},
	't':  {},
	'n':  {},
	'f':  {},
	'r':  {},
	'"':  {},
	'\'': {},
}

func lenUnescaped(str string) int {
	length := len(str)
	acceptEscaped := false
	for _, char := range str {
		if _, ok := escapeCharSet[char]; acceptEscaped && ok {
			length--
		}
		acceptEscaped = char == '\\'
	}
	return length - 2
}

// NewStringLiteral creates a string literal, and automatically inserts the sentinel character
func NewStringLiteral(ID string, value string) Instruction {
	strLit := StringLiteral{
		String: value,
	}
	strLit.ID = ID
	strLit.Size = lenUnescaped(value)
	return strLit
}
