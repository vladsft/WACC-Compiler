package instructions

import (
	"errors"
	"fmt"
)

const (
	allOnes = 0xFFFFFFFFFFFFFFFF
	maxRegs = 64
)

//RegisterManager manages 2 collections of registers: callee and caller saved registers
type RegisterManager struct {
	callee *registers
	caller *registers
}

//NewRegisterManager creates a RegisterManager with nRegs registers
func NewRegisterManager(calleeSaved, callerSaved uint64, nRegs int) *RegisterManager {
	regMng := &RegisterManager{
		callee: newRegisters(calleeSaved, nRegs),
		caller: newRegisters(callerSaved, nRegs),
	}
	return regMng
}

func (rm *RegisterManager) AcquireCalleeSaved() (Register, error) {
	return rm.callee.acquire()
}

func (rm *RegisterManager) AcquireCallerSaved() (Register, error) {
	return rm.caller.acquire()
}

func (rm *RegisterManager) ReleaseCalleeSaved(reg Register) {
	rm.callee.release(reg)
}

func (rm *RegisterManager) ReleaseCallerSaved(reg Register) {
	rm.caller.release(reg)
}

func (rm *RegisterManager) GetUsedCallerSaved() []Register {
	return rm.caller.getUsedRegs()
}

func (rm *RegisterManager) GetUsedCalleeSaved() []Register {
	return rm.callee.getUsedRegs()
}

//registers tracks which registers are in use
type registers struct {
	nRegs       int
	useableRegs uint64
	regs        uint64
}

//newRegisters creates a collection of registers with nRegs registers
func newRegisters(regs uint64, nRegs int) *registers {
	return &registers{
		nRegs:       nRegs,
		useableRegs: regs ^ allOnes,
		regs:        regs ^ allOnes,
	}
}

func (rs *registers) get(reg Register) error {
	shifted := uint64(1 << reg)
	if rs.useableRegs&shifted == 1 {
		return errors.New("Register already in use")
	}
	rs.useableRegs |= shifted
	return nil
}

//acquire returns a free register, errors if none are left
func (rs *registers) acquire() (Register, error) {
	var reg Register
	for i := 0; i < rs.nRegs; i++ {
		if (1<<i)&rs.useableRegs == 0 {
			rs.useableRegs |= 1 << i
			return Register(i), nil
		} else if (1<<i)&rs.regs == 0 {
			reg = Register(i)
		}
	}
	rs.release(reg)
	return -1, errors.New("Could not acquire register")
}

//release marks the register as "free", panics if the register is already free
func (rs *registers) release(reg Register) {
	shifted := uint64(1 << reg)
	if shifted&(rs.useableRegs^rs.regs) == 0 {
		panic(fmt.Errorf("cannot release %d, not in use", reg))
	}
	rs.useableRegs &= allOnes ^ shifted
}

//getUsedRegs returns a list of all used registers
func (rs *registers) getUsedRegs() []Register {
	regs := make([]Register, 0, rs.nRegs)
	for i := 0; i < rs.nRegs; i++ {
		if (1<<i)&(rs.useableRegs^rs.regs) != 0 {
			regs = append(regs, Register(i))
		}
	}
	return regs
}
