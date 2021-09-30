package architecture

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/types"
)

//Emitter should be used to convert our internal assembly into actual assembly
type Emitter interface {
	Emit(bss ins.Instruction, instrs ins.Instruction) string
	EmitInstructions(ins.Instructions) string
	EmitNOOP(ins.NOOP) string
	EmitPool(ins.Pool) string
	EmitExit(ins.Exit) string
	EmitLabel(ins.Label) string
	EmitFunctionCall(ins.FunctionCall) string
	EmitBranch(ins.Branch) string
	EmitMove(ins.Move) string
	EmitCompare(ins.Compare) string
	EmitAdd(ins.Add) string
	EmitSub(ins.Sub) string
	EmitMult(ins.Mult) string
	EmitDiv(ins.Div) string
	EmitMod(ins.Mod) string
	EmitXor(ins.Xor) string
	EmitStringLiteral(ins.StringLiteral) string
	EmitBoolExpr(ins.BoolExpr) string
	EmitAnd(ins.And) string
	EmitOr(ins.Or) string
	EmitNeg(ins.Neg) string
	EmitSize(types.Size) string
	EmitLoad(ins.Load) string
	EmitStackInstr(ins.StackInstr) string
	EmitStore(ins.Store) string
	EmitStoreHeap(ins.StoreHeap) string
	EmitFreeHeap(ins.FreeHeap) string

	EmitOperand(ins.Operand) string
	EmitImmediate(ins.Immediate) string
	EmitPseudoImmediate(ins.PseudoImmediate) string
	EmitAddress(ins.Address) string
	EmitVariable(ins.Variable) string
	EmitRegister(ins.Register) string
}
