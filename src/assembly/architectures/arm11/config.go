package arm11

import architecture "wacc_32/assembly/architectures"

func Config() architecture.Config {
	return architecture.Config{
		ProgramCounter:  15,
		LinkRegister:    14,
		StackPointer:    13,
		ScratchRegister: 12,
		Accumulator:     11,
		ReturnRegister:  0,
		CalleeSavedRegs: 0xF,
		CallerSavedRegs: 0x7F0,
		NRegs:           11,
	}
}
