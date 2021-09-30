package assembly

import (
	"strconv"
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/types"
)

//VisitTernaryOp visits AST node ast.TernaryOp
func (cg *CodeGenerator) VisitTernaryOp(node ast.TernaryOp, ctx *ins.Context) ins.ExprInstr {
	cond := node.GetCondition()

	condExpr := cg.VisitExpression(cond, ctx)

	switch condType := cond.(type) {
	case *ast.Literal:
		if condType.GetValue().(bool) {
			return cg.VisitExpression(node.GetIfExpr(), ctx)
		}
		return cg.VisitExpression(node.GetElseExpr(), ctx)
	}

	loadInstrs, reg, _ := cg.loadIfNotRegister(condExpr.Result, types.Byte)
	cg.regs.ReleaseCallerSaved(reg)

	instrs := condExpr.Instrs
	ifExpr := cg.VisitExpression(node.GetIfExpr(), ctx)
	ifInstrs, ifReg, _, ifStack := cg.loadRegisterWithStack(ifExpr.Result,
		exprSize(node.GetIfExpr()))
	cg.regs.ReleaseCallerSaved(ifReg)

	if ifStack {
		ifInstrs = addHeader(ins.NewPush(10), ifInstrs)
	}
	elseExpr := cg.VisitExpression(node.GetElseExpr(), ctx)
	elseInstrs, _, _, elseStack := cg.loadRegisterWithStack(elseExpr.Result,
		exprSize(node.GetElseExpr()))
	if elseStack {
		elseInstrs = addHeader(ins.NewPush(10), elseInstrs)
	}

	if ifStack || elseStack {
		instrs = append(instrs, ins.NewPop(cg.Accumulator))
	}

	nIfs := strconv.Itoa(cg.nIfs)
	elseLabel := "else_" + nIfs
	endLabel := "end_" + nIfs

	cg.nIfs++

	return ins.NewExprInstr(ins.Instructions{
		instrs,
		loadInstrs,
		ins.NewCompare(reg, ins.Immediate(0)),
		ins.NewBranch(elseLabel, ins.EQ),
		ifInstrs,
		ifExpr.Instrs,
		ins.NewBranch(endLabel, ins.AL),
		ins.NewLabel(elseLabel),
		elseInstrs,
		elseExpr.Instrs,
		ins.NewLabel(endLabel),
	}, ifReg)
}
