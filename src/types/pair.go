package types

import "fmt"

var _ WaccType = pair{}

type pair struct {
	fstType WaccType
	sndType WaccType
}

//NewPair creates a new pair
func NewPair(wt1, wt2 WaccType) WaccType {
	return pair{
		fstType: wt1,
		sndType: wt2,
	}
}

//DefaultValue should not be called
func (p pair) DefaultValue() interface{} {
	panic("Default Value on function should not have been called as we do not allow higher order function")
}

func (p pair) GetFormatString() string {
	return "%p"
}

func (p pair) Is(wt WaccType) bool {
	switch w := wt.(type) {
	case waccBaseType:
		return w == Pair
	case pair:
		return p.fstType.Is(w.fstType) && p.sndType.Is(w.sndType)
	default:
		return false
	}
}

func (p pair) String() string {
	return fmt.Sprintf("pair(%s,%s)", p.fstType, p.sndType)
}

func (p pair) GetChildren() []WaccType {
	return []WaccType{p.fstType, p.sndType}
}
