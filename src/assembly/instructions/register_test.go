package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	nCallee     = 2
	nCaller     = 3
	nRegs       = nCallee + nCaller
	calleeSaved = 0x3  //first 2 regs are calle saved
	callerSaved = 0x1c //last 3 regs are caller saved
)

func TestAcquireIsOrdered(t *testing.T) {
	rm := newRegisters(0x1f, nRegs)
	a := Register(-1)
	for i := 0; i < nRegs; i++ {
		b, err := rm.acquire()
		if err != nil || b <= a {
			t.Fail()
		}
		a = b
	}
}

func TestAcquireBounds(t *testing.T) {
	rm := newRegisters(0x1f, nRegs)
	for i := 0; i < nRegs; i++ {
		_, err := rm.acquire()
		assert.Nil(t, err)
	}
	reg, err := rm.acquire()
	assert.Equal(t, Register(-1), reg)
	assert.Error(t, err)
}

func TestCanReleaseAcquiredRegs(t *testing.T) {
	rm := newRegisters(0x1f, nRegs)
	x, _ := rm.acquire()
	y, _ := rm.acquire()
	z, _ := rm.acquire()
	rm.release(z)
	rm.release(x)
	rm.release(y)
}

func TestPanicOnInvalidRelease(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Fail()
		}
	}()
	rm := newRegisters(0x1f, nRegs)
	rm.acquire()
	rm.acquire()
	rm.release(32)
}

func TestGetUsedRegs(t *testing.T) {
	rm := newRegisters(0x1f, nRegs)
	for i := 0; i < nRegs; i++ {
		rm.acquire()
	}
	assert.Equal(t, nRegs, len(rm.getUsedRegs()))
}

func TestAcquiresDontOverlap(t *testing.T) {
	rm := NewRegisterManager(calleeSaved, callerSaved, nRegs)
	for i := 0; i < 2; i++ {
		r, err := rm.AcquireCalleeSaved()
		assert.Nil(t, err)
		assert.Equal(t, r, Register(i))
	}
	for i := nCallee; i < nRegs; i++ {
		r, err := rm.AcquireCallerSaved()
		assert.Nil(t, err)
		assert.Equal(t, r, Register(i))
	}
}
