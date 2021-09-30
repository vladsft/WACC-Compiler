package instructions

type binop struct {
	Dest        Register
	Left, Right Operand
}

type Add binop

//NewAdd stores the result of left + right in dest
func NewAdd(dest Register, left, right Operand) Instruction {
	return Add{
		Dest:  dest,
		Left:  left,
		Right: right,
	}
}

type Sub binop

//NewSub stores the result of left-right in dest
func NewSub(dest Register, left, right Operand) Instruction {
	return Sub{
		Dest:  dest,
		Left:  left,
		Right: right,
	}
}

type Mult binop

//NewMult stores the result of left * right in dest
func NewMult(dest Register, left, right Operand) Instruction {
	return Mult{
		Dest:  dest,
		Left:  left,
		Right: right,
	}
}

type Div binop

//NewDiv stores the result of left / right in dest
func NewDiv(dest Register, left, right Operand) Instruction {
	return Div{
		Dest:  dest,
		Left:  left,
		Right: right,
	}
}

type Mod binop

// NewMod stores the result of left % right in dest
func NewMod(dest Register, left, right Operand) Instruction {
	return Mod{
		Dest:  dest,
		Left:  left,
		Right: right,
	}
}

type Xor binop

//NewXor stores the result of left ^ right in dest
func NewXor(dest Register, left, right Operand) Instruction {
	return Xor{
		Dest:  dest,
		Left:  left,
		Right: right,
	}
}
