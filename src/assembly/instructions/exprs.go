package instructions

var (
	_ Instruction = BoolExpr{}
)

type Cond int

//go:generate stringer -type=Cond
const (
	EQ Cond = iota + 1
	NE
	GT
	LE
	GE
	LT
	AL
	NV
	VS
)

type HasCondtion interface {
	GetCondition() Cond
}

type BoolExpr struct {
	Dest        Register
	Left, Right Operand
	True, False Cond
}

func (b BoolExpr) GetCondition() Cond {
	return b.False
}

//This is maps n to n+1 and n+1 to n where n is odd
func (c Cond) opposite() Cond {
	return c + (2 * (c % 2)) - 1
}

//NewBoolean returns a boolean expression, automatically inferring the false case
func NewBoolean(dest Register, left, right Operand, test Cond) BoolExpr {
	return BoolExpr{
		Left:  left,
		Right: right,
		Dest:  dest,
		True:  test,
		False: test.opposite(),
	}
}

type And struct {
	BoolExpr
}

//NewAnd stores the result of left && right in dest
func NewAnd(dest Register, left, right Operand) Instruction {
	return And{NewBoolean(dest, left, right, EQ)}
}

type Or struct {
	BoolExpr
}

//NewOr stores the result of left || right in dest
func NewOr(dest Register, left, right Operand) Instruction {
	return Or{NewBoolean(dest, left, right, EQ)}
}

type Neg struct {
	Reg Register
}

func NewNeg(reg Register) Instruction {
	return Neg{
		Reg: reg,
	}
}
