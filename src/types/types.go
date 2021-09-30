package types

import (
	"fmt"
)

//WaccType interface
type WaccType interface {
	fmt.Stringer
	Is(wt WaccType) bool
	GetChildren() []WaccType
	GetFormatString() string
	DefaultValue() interface{}
}

var _ WaccType = Integer

var typeStrings = []string{"NULL", "int", "bool", "char", "string", "pair", "array", "function", "structure", "lock", "sema"}
var typeFormatStrings = []string{"%p", "%d", "true\\0false", " %c", "%.*s", "%p", "%p", "", "", "%p", "%p"}
var defaultValues = []interface{}{
	nil,
	0,
	true,
	"",
	"",
	nil,
	nil,
	[]interface{}{},
	nil,
	nil,
	nil,
	nil,
}

type waccBaseType int

//BaseTypes representing the different types in the Wacc language
const (
	None waccBaseType = iota + 1
	Integer
	Boolean
	Char
	Str
	Pair
	Array
	Function
	UserDefinedType
	Lock
	Sema
)

func (wbt waccBaseType) String() string {
	return typeStrings[wbt-1]
}

func (wbt waccBaseType) DefaultValue() interface{} {
	return defaultValues[wbt-1]
}

//GetFormatString returns the C format string specifier for the wacc type
func (wbt waccBaseType) GetFormatString() string {
	return typeFormatStrings[wbt-1]
}

//Is checks if the BaseType of WBT is the same as WT
func (wbt waccBaseType) Is(wt WaccType) bool {
	switch w := wt.(type) {
	case array:
		if wbt == Str && w.baseType == Char && w.depth == 1 {
			return true
		}
		return wbt == Array
	case pair:
		return wbt == Pair
	case function:
		return wbt == Function
	default:
		return wbt == w
	}
}

func (wbt waccBaseType) GetChildren() []WaccType {
	return []WaccType{}
}

//Size represents the size of data in memory
type Size int

//The different sizes available
const (
	Byte = 1 << iota
	HalfWord
	Word
)

var sizeStrings = []string{"b", "sh", ""}

func (s Size) String() string {
	return sizeStrings[s>>1]
}

func TypeSize(wt WaccType) Size {
	switch wt {
	case Boolean:
		fallthrough
	case Char:
		return Byte
	default:
		return Word
	}
}
