package instructions

//Instruction is our internal assembly representation
//It can be converted to an assembly language
type Instruction interface{}

//Context is carried by visitors
type Context struct{}

type StackInstrType int

//Instructions is a list of instructions
type Instructions []Instruction

//NOOP is a no op it generates nothing
type NOOP struct{}

type Pool struct{}
type Exit struct {
	Code Operand
}

//NewExit creates an exit instruction with an exit code
func NewExit(code Operand) Instruction {
	return Exit{
		Code: code,
	}
}

type Label struct {
	Name string
}

//NewLabel creates a new label
func NewLabel(name string) Instruction {
	return Label{
		Name: name,
	}
}

type FunctionCall struct {
	Name string
}

//NewFunctionCall branches to a function
func NewFunctionCall(name string) Instruction {
	return FunctionCall{
		Name: name,
	}
}

type Branch struct {
	Label     string
	Condition Cond
}

//NewBranch branches to a label if the condition is met
func NewBranch(label string, condition Cond) Instruction {
	return Branch{
		Label:     label,
		Condition: condition,
	}
}

type Move struct {
	Src  Operand
	Dest Register
}

//NewMove moves src to dest
func NewMove(src Operand, dest Register) Instruction {
	return Move{
		Src:  src,
		Dest: dest,
	}
}

type Compare struct {
	Left  Operand
	Right Operand
}

func NewCompare(left Operand, right Operand) Instruction {
	return Compare{
		Left:  left,
		Right: right,
	}
}

// func TypeSize(wt types.WaccType) Size {
// 	switch wt {
// 	case types.Boolean:
// 		fallthrough
// 	case types.Char:
// 		return Byte
// 	default:
// 		return Word
// 	}
// }
