package assembly

import (
	"wacc_32/assembly/builtins"
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
)

//VisitBinOp visits AST node ast.BinOp
func (cg *CodeGenerator) VisitBinOp(node ast.BinOp, ctx *ins.Context) ins.ExprInstr {
	lexpr := node.GetLeftExpr()
	rexpr := node.GetRightExpr()

	left := cg.VisitExpression(lexpr, ctx)
	lInstrs, lReg, _, lStack := cg.loadRegisterWithStack(left.Result, exprSize(lexpr))

	if lStack {
		lInstrs = addHeader(ins.NewPush(10), lInstrs)
	}

	right := cg.VisitExpression(rexpr, ctx)
	rInstrs, rReg, _, rStack := cg.loadRegisterWithStack(right.Result, exprSize(rexpr))
	if rStack {
		rInstrs = addHeader(ins.NewPush(10), rInstrs)
	} else {
		if rReg != 10 {
			cg.regs.ReleaseCallerSaved(rReg)
		}
	}

	instrs := append(left.Instrs, lInstrs...)
	instrs = append(instrs, right.Instrs...)
	instrs = append(instrs, rInstrs...)

	if lStack || rStack || lReg == rReg {
		rReg = cg.Accumulator
		instrs = append(instrs, ins.NewPop(cg.Accumulator))
	}

	var instr ins.Instruction
	cg.funcs["p_print_error"] = builtins.PrintError()
	switch node.GetOpType() {
	case ast.Star:
		instr = ins.NewMult(lReg, lReg, rReg)
		instr = addTail([]ins.Instruction{instr}, ins.NewFunctionCall("p_check_int_overflow"))
		cg.addOverflowCode()
	case ast.Div:
		instr = ins.NewDiv(lReg, lReg, rReg)
		cg.addDivideByZeroCode()
	case ast.Mod:
		instr = ins.NewMod(lReg, lReg, rReg)
		cg.addDivideByZeroCode()
	case ast.Plus:
		instr = ins.NewAdd(lReg, lReg, rReg)
		instr = addTail([]ins.Instruction{instr}, ins.NewFunctionCall("p_check_int_overflow"))
		cg.addOverflowCode()
	case ast.Minus:
		instr = ins.NewSub(lReg, lReg, rReg)
		instr = addTail([]ins.Instruction{instr}, ins.NewFunctionCall("p_check_int_overflow"))
		cg.addOverflowCode()
	case ast.Greater:
		instr = ins.NewBoolean(lReg, lReg, rReg, ins.GT)
	case ast.GreaterEq:
		instr = ins.NewBoolean(lReg, lReg, rReg, ins.GE)
	case ast.Less:
		instr = ins.NewBoolean(lReg, lReg, rReg, ins.LT)
	case ast.LessEq:
		instr = ins.NewBoolean(lReg, lReg, rReg, ins.LE)
	case ast.Equal:
		instr = ins.NewBoolean(lReg, lReg, rReg, ins.EQ)
	case ast.NotEq:
		instr = ins.NewBoolean(lReg, lReg, rReg, ins.NE)
	case ast.And:
		instr = ins.NewAnd(lReg, lReg, rReg)
	case ast.Or:
		instr = ins.NewOr(lReg, lReg, rReg)
	}

	instrs = append(instrs, instr)
	return ins.NewExprInstr(instrs, lReg)
}
