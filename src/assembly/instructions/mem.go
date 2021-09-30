package instructions

import "wacc_32/types"

//go:generate stringer -type=StackInstrType
const (
	push StackInstrType = iota + 1
	pop
)

// //Size represents the size of data in memory
// type Size int

// //The different sizes available
// const (
// 	Byte = 1 << iota
// 	HalfWord
// 	Word
// )

// var sizeStrings = []string{"b", "sh", ""}

// func (s Size) String() string {
// 	return sizeStrings[s>>1]
// }

var (
	_ Instruction = Load{}
	_ Instruction = StackInstr{}
	_ Instruction = StoreHeap{}
)

type Load struct {
	Src  Operand
	Dest Register
	Size types.Size
}

//NewLoad creates a load instruction which loads dest with src
func NewLoad(src Operand, dest Register, size types.Size) Instruction {
	return Load{
		Src:  src,
		Dest: dest,
		Size: size,
	}
}

//NewIncrementStack pops size bytes from the stack to nothing
func NewIncrementStack(size int, stackPointer Register) Instruction {
	return NewAdd(Register(stackPointer), Register(stackPointer), Immediate(size))
}

//NewDecrementStack pushes size bytes to the stack
func NewDecrementStack(size int, stackPointer Register) Instruction {
	return NewSub(stackPointer, stackPointer, Immediate(size))
}

type StackInstr struct {
	T    StackInstrType
	Regs []Register
}

//NewPush creates a push
func NewPush(regs ...Register) Instruction {
	return StackInstr{
		T:    push,
		Regs: regs,
	}
}

//NewPop creates a pop
func NewPop(reg ...Register) Instruction {
	return StackInstr{
		T:    pop,
		Regs: reg,
	}
}

type Store struct {
	Size types.Size
	Src  Register
	Dest Operand
}

//NewStore stores src in dest
func NewStore(size types.Size, src Register, dest Operand) Instruction {
	return Store{
		Size: size,
		Src:  src,
		Dest: dest,
	}
}

type StoreHeap struct {
	Op, Size           Operand
	Offset, Multiplier int
}

//NewStoreHeap stores the value of a register at a heap address
func NewStoreHeap(op, size Operand, offset, multiplier int) Instruction {
	return StoreHeap{
		Op:         op,
		Size:       size,
		Offset:     offset,
		Multiplier: multiplier,
	}
}

type FreeHeap struct {
	Reg Register
}

//NewFreeHeap frees the memory location stored in the register
//Note this is unsafe as it does no checking
func NewFreeHeap(reg Register) Instruction {
	return FreeHeap{
		Reg: reg,
	}
}
