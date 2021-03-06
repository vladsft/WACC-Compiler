// This file was automatically generated by genny.
// Any changes will be lost if this file is regenerated.
// see https://github.com/cheekybits/genny

// Code generated by visitor_generator. DO NOT EDIT.
package ast

import "wacc_32/assembly/instructions"

type InstructionAcceptor interface {
	AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction
}

type ExprInstrAcceptor interface {
	AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr
}

//Accept calls v.VisitRHSNewPair(r)
func (r RHSNewPair) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitRHSNewPair(r, ctx)
}

//Accept calls v.VisitRHSFunctionCall(r)
func (r RHSFunctionCall) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitRHSFunctionCall(r, ctx)
}

//Accept calls v.VisitPairElem(p)
func (p PairElem) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitPairElem(p, ctx)
}

//Accept calls v.VisitMake(m)
func (m Make) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitMake(m, ctx)
}

//Accept calls v.VisitBinOp(b)
func (b BinOp) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitBinOp(b, ctx)
}

//Accept calls v.VisitStatLock(s)
func (s StatLock) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatLock(s, ctx)
}

//Accept calls v.VisitStatSema(s)
func (s StatSema) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatSema(s, ctx)
}

//Accept calls v.VisitArrayElem(a)
func (a ArrayElem) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitArrayElem(a, ctx)
}

//Accept calls v.VisitIdent(i)
func (i Ident) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitIdent(i, ctx)
}

//Accept calls v.VisitFunction(f)
func (f Function) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitFunction(f, ctx)
}

//Accept calls v.VisitParamList(p)
func (p ParamList) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitParamList(p, ctx)
}

//Accept calls v.VisitParam(p)
func (p Param) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitParam(p, ctx)
}

//Accept calls v.VisitProgram(p)
func (p Program) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitProgram(p, ctx)
}

//Accept calls v.VisitStatSkip(s)
func (s StatSkip) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatSkip(s, ctx)
}

//Accept calls v.VisitStatRead(s)
func (s StatRead) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatRead(s, ctx)
}

//Accept calls v.VisitStatFree(s)
func (s StatFree) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatFree(s, ctx)
}

//Accept calls v.VisitStatNewassign(s)
func (s StatNewassign) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatNewassign(s, ctx)
}

//Accept calls v.VisitStatPrint(s)
func (s StatPrint) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatPrint(s, ctx)
}

//Accept calls v.VisitStatPrintln(s)
func (s StatPrintln) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatPrintln(s, ctx)
}

//Accept calls v.VisitStatExit(s)
func (s StatExit) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatExit(s, ctx)
}

//Accept calls v.VisitStatFor(s)
func (s StatFor) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatFor(s, ctx)
}

//Accept calls v.VisitStatWhile(s)
func (s StatWhile) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatWhile(s, ctx)
}

//Accept calls v.VisitStatDoWhile(s)
func (s StatDoWhile) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatDoWhile(s, ctx)
}

//Accept calls v.VisitStatBegin(s)
func (s StatBegin) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatBegin(s, ctx)
}

//Accept calls v.VisitStatAssign(s)
func (s StatAssign) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatAssign(s, ctx)
}

//Accept calls v.VisitStatReturn(s)
func (s StatReturn) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatReturn(s, ctx)
}

//Accept calls v.VisitStatIf(s)
func (s StatIf) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatIf(s, ctx)
}

//Accept calls v.VisitWaccRoutine(w)
func (w WaccRoutine) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitWaccRoutine(w, ctx)
}

//Accept calls v.VisitStatMultiple(s)
func (s StatMultiple) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitStatMultiple(s, ctx)
}

//Accept calls v.VisitTernaryOp(t)
func (t TernaryOp) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitTernaryOp(t, ctx)
}

//Accept calls v.VisitLiteral(l)
func (l Literal) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitLiteral(l, ctx)
}

//Accept calls v.VisitUnOp(u)
func (u UnOp) AcceptExprInstr(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.ExprInstr {
	return v.VisitUnOp(u, ctx)
}

//Accept calls v.VisitUserType(u)
func (u UserType) AcceptInstruction(v InstructionExprInstrVisitor, ctx *instructions.Context) instructions.Instruction {
	return v.VisitUserType(u, ctx)
}
