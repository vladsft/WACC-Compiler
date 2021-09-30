package assembly

import (
	"wacc_32/assembly/builtins"
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
)

//VisitFunction visits AST node ast.Function
func (cg *CodeGenerator) VisitFunction(node ast.Function, ctx *ins.Context) ins.Instruction {
	subCtx := &ins.Context{}

	name := node.GetName()

	statInstrs, early := cg.visitStatsTillEnd(node.GetStats(), subCtx)
	sym := node.GetSymbolTable()
	offset := sym.GetTotalOffset()
	withStack := wrapInstructions(
		ins.NewDecrementStack(offset, cg.StackPointer),
		statInstrs,
		ins.NewIncrementStack(offset, cg.StackPointer),
	)
	withPushPop := addHeader(
		ins.NewPush(cg.LinkRegister),
		withStack,
	)

	withLabel := addHeader(ins.NewLabel(name), withPushPop)
	if name == "main" && !early {
		withLabel = addTail(withLabel, ins.NewExit(ins.Immediate(0)))
		return addTail(withLabel, ins.NewPop(cg.ProgramCounter))
	}

	//Add concurrent header if the function is concurrent
	//The concurrent header sets up the stack of the function, from the pointer in r0
	concurrentHeader := cg.generateConcurrentHeader(name, offset, len(node.GetParams()))
	withLabel = append(concurrentHeader, withLabel...)

	cg.funcs[name] = append(withLabel, ins.NewPop(cg.ProgramCounter), ins.Pool{})
	return ins.NOOP{}
}

//VisitParamList visits AST node ast.ParamList
func (cg *CodeGenerator) VisitParamList(node ast.ParamList, ctx *ins.Context) ins.Instruction {
	return ins.NOOP{}
}

//VisitParam visits AST node ast.Param
func (cg *CodeGenerator) VisitParam(node ast.Param, ctx *ins.Context) ins.Instruction {
	return ins.NOOP{}
}

//VisitStatReturn visits AST node ast.StatReturn
func (cg *CodeGenerator) VisitStatReturn(node ast.StatReturn, ctx *ins.Context) ins.Instruction {
	exprInstrs := cg.VisitExpression(node.GetReturnExpr(), ctx)
	var instrs ins.Instructions
	switch op := exprInstrs.Result.(type) {
	case ins.Address:
		instrs = append(exprInstrs.Instrs,
			ins.NewLoad(exprInstrs.Result, cg.ReturnRegister, exprSize(node.GetReturnExpr())))
	case ins.Register:
		cg.regs.ReleaseCallerSaved(op)
		instrs = append(exprInstrs.Instrs, ins.NewMove(exprInstrs.Result, cg.ReturnRegister))
	default:
		instrs = append(exprInstrs.Instrs, ins.NewMove(exprInstrs.Result, cg.ReturnRegister))
	}
	offset := node.GetSymbolTable().GetOffset(node.GetFunctionName())
	instrs = append(instrs,
		ins.NewAdd(cg.StackPointer, cg.StackPointer, ins.Immediate(offset)),
		ins.NewPop(cg.ProgramCounter),
	)
	return instrs
}

func (cg *CodeGenerator) addOutOfBoundsCode() {
	cg.funcs[builtins.ArrayOutOfBoundsCheckLabel] = builtins.CheckArrayIndexOutOfBounds()
	cg.funcs[builtins.PrintErrorCheckLabel] = builtins.PrintError()
	cg.addErrMsgBSS(builtins.ArrayIndexNegativeError)
	cg.addErrMsgBSS(builtins.ArrayIndexTooLargeError)
}
func (cg *CodeGenerator) addDivideByZeroCode() {
	cg.funcs[builtins.DivideByZeroCheckLabel] = builtins.CheckDivideByZeroError()
	cg.funcs[builtins.PrintErrorCheckLabel] = builtins.PrintError()
	cg.addErrMsgBSS(builtins.DivideByZeroError)
}

func (cg *CodeGenerator) addOverflowCode() {
	cg.funcs[builtins.IntegerOverflowCheckLabel] = builtins.CheckIntegerOverflowError()
	cg.funcs[builtins.PrintErrorCheckLabel] = builtins.PrintError()
	cg.addErrMsgBSS(builtins.IntegerOverflowError)
}

func (cg *CodeGenerator) addNullPointerReferenceCode() {
	cg.funcs[builtins.NullPointerReferenceCheckLabel] = builtins.CheckNullPointerReferenceError()
	cg.funcs[builtins.PrintErrorCheckLabel] = builtins.PrintError()
	cg.addErrMsgBSS(builtins.NullPointerReferenceError)
}

func (cg *CodeGenerator) addInvalidThreadUnlockCode() {
	cg.funcs[builtins.InvalidThreadUnlockCheckLabel] = builtins.CheckInvalidThreadUnlock()
	cg.funcs[builtins.PrintErrorCheckLabel] = builtins.PrintError()
	cg.addErrMsgBSS(builtins.InvalidThreadUnlockError)
}

func (cg *CodeGenerator) addCheckSameThreadLockCode() {
	cg.funcs[builtins.SameThreadLockCheckLabel] = builtins.CheckSameThreadLock()
	cg.funcs[builtins.PrintErrorCheckLabel] = builtins.PrintError()
	cg.addErrMsgBSS(builtins.SameThreadLockError)
}
