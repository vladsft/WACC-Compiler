package instructions

type operandType int

//Operand represents the different types of operand that can exist
type Operand interface {
	hidden()
}

//ExprInstr should be used for encoding instructions
type ExprInstr struct {
	Instrs Instructions
	Result Operand
}

//NewExprInstr returns a new expression instruction
func NewExprInstr(instrs []Instruction, result Operand) ExprInstr {
	return ExprInstr{
		Instrs: instrs,
		Result: result,
	}
}

//NewEmptyExprInstr returns an empty expression instruction
func NewEmptyExprInstr() ExprInstr {
	return ExprInstr{}
}

// Register represents a register
type Register int

type Immediate int

type PseudoImmediate int

type FunctionPointer string
type Address struct {
	Reg                Register
	Offset, Multiplier Immediate
}

//NewAddress returns an address, imms[0:2] are used
func NewAddress(reg Register, imms ...Immediate) Operand {
	addr := Address{
		Reg: reg,
	}
	if len(imms) == 2 {
		addr.Multiplier = imms[1]
		addr.Offset = imms[0]
	} else if len(imms) == 1 {
		addr.Offset = imms[0]
	}
	return addr
}

type Variable string

type Directive string

func NewDirective(name string) Operand {
	return Directive("." + name)
}

func (i Immediate) hidden()        {}
func (ps PseudoImmediate) hidden() {}
func (a Address) hidden()          {}
func (r Register) hidden()         {}
func (v Variable) hidden()         {}
func (f FunctionPointer) hidden()  {}
func (d Directive) hidden()        {}
