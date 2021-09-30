package ast

import (
	"wacc_32/errors"
	"wacc_32/types"
)

var (
	_ Statement = &StatLock{}
	_ Statement = &StatSema{}
)

type LockStatType int

const (
	Acquire LockStatType = iota + 1
	Release
)

var lockStatStrings = []string{"ACQUIRE", "RELEASE"}

func (lst LockStatType) String() string {
	return lockStatStrings[lst-1]
}

//StatLock represents a lock operation (either locking or unlocking)
type StatLock struct {
	ast
	sType LockStatType
	lock  *Ident
	pos   errors.Position
}

func NewAcquire(lock *Ident, pos errors.Position) *StatLock {
	return newStatLock(Acquire, lock, pos)
}
func NewRelease(lock *Ident, pos errors.Position) *StatLock {
	return newStatLock(Release, lock, pos)
}

func newStatLock(sType LockStatType, lock *Ident, pos errors.Position) *StatLock {
	return &StatLock{
		sType: sType,
		lock:  lock,
		pos:   pos,
	}
}

func (s StatLock) String() string {
	return format(s.sType.String(), s.lock.String())
}

func (s StatLock) GetName() string {
	return s.lock.name
}

func (s StatLock) GetType() LockStatType {
	return s.sType
}

func (s *StatLock) Check(ctx Context) {
	s.table = ctx.table
	s.lock.table = ctx.table
	if !s.lock.Check(ctx) {
		return
	}
	if t := s.lock.EvalType(*s.table); !t.Is(types.Lock) {
		ctx.SemanticErrChan <- errors.NewTypeError(s.pos, s.sType.String(), types.Lock, t)
	}
}

type StatSema struct {
	ast
	up   bool
	sema *Ident
	pos  errors.Position
}

func NewSemaUp(sema *Ident, pos errors.Position) Statement {
	return &StatSema{
		up:   true,
		sema: sema,
		pos:  pos,
	}
}
func NewSemaDown(sema *Ident, pos errors.Position) Statement {
	return &StatSema{
		up:   false,
		sema: sema,
		pos:  pos,
	}
}

func (s StatSema) GetIdent() *Ident {
	return s.sema
}

func (s StatSema) IsUp() bool {
	return s.up
}

func (s StatSema) getName() string {
	if s.up {
		return "up"
	}
	return "down"
}

func (s StatSema) String() string {
	return format(s.getName(), s.sema.String())
}

func (s StatSema) Check(ctx Context) {
	s.table = ctx.table
	s.sema.table = ctx.table
	if !s.sema.Check(ctx) {
		return
	}
	if t := s.sema.EvalType(*s.table); !t.Is(types.Sema) {
		ctx.SemanticErrChan <- errors.NewTypeError(s.pos, s.getName(), types.Sema, t)
	}
}
