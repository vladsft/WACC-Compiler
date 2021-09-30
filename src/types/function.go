package types

var _ WaccType = function{}

type function struct {
	returnType WaccType
	paramTypes []WaccType
}

//NewFunction creates a new function type
func NewFunction(returnType WaccType, paramTypes []WaccType) WaccType {
	return function{
		returnType: returnType,
		paramTypes: paramTypes,
	}
}

//DefaultValue should not be called
func (f function) DefaultValue() interface{} {
	panic("Default Value on function should not have been called as we do not allow higher order function")
}

func (f function) GetFormatString() string {
	return "%p"
}

func (f function) Is(wt WaccType) bool {
	switch w := wt.(type) {
	case waccBaseType:
		return w == Function
	case function:
		for i, fType := range f.paramTypes {
			if w.paramTypes[i] != fType {
				return false
			}
		}
		return f.returnType == w.returnType
	default:
		return false
	}
}

func (f function) String() string {
	var str string
	for _, pType := range f.paramTypes {
		str += pType.String() + " -> "
	}
	return str + f.returnType.String()
}

func (f function) GetChildren() []WaccType {
	return append(f.paramTypes, f.returnType)
}
