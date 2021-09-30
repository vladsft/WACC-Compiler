package types

var _ WaccType = array{}

type array struct {
	depth    int
	baseType WaccType
}

//NewArray creates a new array - assumes dimension >= 1
func NewArray(wt WaccType, depth int) WaccType {
	if depth < 1 {
		panic("Incorrect depth given")
	}
	if wt.Is(Array) {
		w := wt.(array)
		return array{
			baseType: w.baseType,
			depth:    w.depth + 1,
		}
	}
	return array{
		baseType: wt,
		depth:    depth,
	}
}

//DefaultValue returns nil
func (a array) DefaultValue() interface{} {
	return nil
}

func (a array) GetFormatString() string {
	return ""
}

func (arr array) Is(wt WaccType) bool {
	switch w := wt.(type) {
	case waccBaseType:
		return w == Array
	case array:
		return w.baseType.Is(arr.baseType) && w.depth == arr.depth
	default:
		return false
	}
}

func (arr array) String() string {
	str := arr.baseType.String()
	for i := 0; i < arr.depth; i++ {
		str += "[]"
	}
	return str
}

func (arr array) GetChildren() []WaccType {
	if arr.depth == 1 {
		return []WaccType{arr.baseType}
	}
	return []WaccType{NewArray(arr.baseType, arr.depth-1)}
}
