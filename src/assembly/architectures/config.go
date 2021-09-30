package architecture

import (
	ins "wacc_32/assembly/instructions"
)

//Config contains all information about an architecture that the wacc compiler needs
type Config struct {
	StackPointer    ins.Register
	ScratchRegister ins.Register
	Accumulator     ins.Register
	ReturnRegister  ins.Register
	LinkRegister    ins.Register
	ProgramCounter  ins.Register
	CalleeSavedRegs uint64
	CallerSavedRegs uint64
	NRegs           int
}
