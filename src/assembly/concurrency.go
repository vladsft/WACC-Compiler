package assembly

import (
	"wacc_32/assembly/builtins"
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/types"
)

//VisitWaccRoutine returns instructions to create and detach a pthread
//mov r0, <size>
//bl malloc
//mov r3, r0
//str arg1, [r3] ...
//sub sp, sp, #4
//mov r0, sp
//ldr r2, <functon>
//mov r1, #0
//bl pthread_create
//ldr r0, [sp]
//bl pthread_detach
//add sp, sp, #4
func (cg *CodeGenerator) VisitWaccRoutine(node ast.WaccRoutine, ctx *ins.Context) ins.Instruction {
	r0, _ := cg.regs.AcquireCalleeSaved()
	r1, _ := cg.regs.AcquireCalleeSaved()
	r2, _ := cg.regs.AcquireCalleeSaved()
	r3, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(r0)
	defer cg.regs.ReleaseCalleeSaved(r1)
	defer cg.regs.ReleaseCalleeSaved(r2)
	defer cg.regs.ReleaseCalleeSaved(r3)

	//1. Allocate memory for the arguments
	argInstrs := ins.Instructions{}
	var argSize types.Size
	for _, arg := range node.GetArgs() {
		exprInstrs := cg.VisitExpression(arg, ctx)
		size := exprSize(arg)
		loadInstrs, reg, _ := cg.loadIfNotRegister(exprInstrs.Result, size)
		cg.regs.ReleaseCallerSaved(reg)
		argInstrs = append(argInstrs, exprInstrs.Instrs...)
		argInstrs = append(argInstrs, loadInstrs)
		argInstrs = append(argInstrs, ins.NewStore(size, reg, ins.NewAddress(r3, ins.Immediate(argSize))))
		argSize += size
	}
	if len(argInstrs) > 0 {
		argInstrs = append(
			ins.Instructions{
				ins.NewStoreHeap(r0, ins.Immediate(argSize), 0, 0),
				ins.NewMove(r0, r3),
			}, argInstrs...)
	}

	//2. Create thread
	name, _ := node.FormatName()
	threadCreateInstrs := ins.Instructions{
		ins.NewSub(cg.StackPointer, cg.StackPointer, ins.Immediate(types.Word)),
		ins.NewMove(cg.StackPointer, r0),
		ins.NewLoad(ins.FunctionPointer(concHeader(name[1:])), r2, types.Word),
		ins.NewMove(ins.Immediate(0), r1),
		ins.NewFunctionCall("pthread_create"),
	}

	//3. Detach thread
	threadDetachInstrs := ins.Instructions{
		ins.NewLoad(ins.NewAddress(cg.StackPointer), r0, types.Word),
		ins.NewFunctionCall("pthread_detach"),
		ins.NewAdd(cg.StackPointer, cg.StackPointer, ins.Immediate(types.Word)),
	}
	return append(append(argInstrs, threadCreateInstrs...), threadDetachInstrs...)
}

func (cg *CodeGenerator) VisitStatLock(node ast.StatLock, ctx *ins.Context) ins.Instruction {
	r0, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(r0)

	offset := node.GetSymbolTable().GetOffset(node.GetName())
	loadLock := ins.Instructions{
		ins.NewLoad(ins.NewAddress(cg.StackPointer, ins.Immediate(offset)), r0, types.Word),
	}
	var function, runtimeCheck string
	switch node.GetType() {
	case ast.Acquire:
		function = "pthread_mutex_lock"
		cg.addCheckSameThreadLockCode()
		runtimeCheck = builtins.SameThreadLockCheckLabel
	case ast.Release:
		function = "pthread_mutex_unlock"
		cg.addInvalidThreadUnlockCode()
		runtimeCheck = builtins.InvalidThreadUnlockCheckLabel
	}
	return append(loadLock, ins.NewFunctionCall(function), ins.NewFunctionCall(runtimeCheck))
}

func concHeader(name string) string {
	return ".." + name + "_conc"
}

func (cg *CodeGenerator) generateConcurrentHeader(name string, offset, nArgs int) ins.Instructions {
	r0, _ := cg.regs.AcquireCalleeSaved()
	r1, _ := cg.regs.AcquireCalleeSaved()
	r2, _ := cg.regs.AcquireCalleeSaved()
	r3, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(r0)
	defer cg.regs.ReleaseCalleeSaved(r1)
	defer cg.regs.ReleaseCalleeSaved(r2)
	defer cg.regs.ReleaseCalleeSaved(r3)

	headerName := concHeader(name)
	label := ins.Instructions{
		ins.NewLabel(headerName),
	}
	if nArgs > 0 {
		label = append(label,
			ins.Instructions{
				ins.NewPush(cg.LinkRegister),
				ins.NewSub(cg.StackPointer, cg.StackPointer, ins.Immediate(offset)),
				ins.NewMove(r0, r1),
				ins.NewMove(r0, r3),
				ins.NewMove(cg.StackPointer, r0),
				ins.NewMove(ins.Immediate(offset), r2),
				ins.NewFunctionCall("memmove"),
				ins.NewAdd(cg.StackPointer, cg.StackPointer, ins.Immediate(offset)),
				ins.NewPop(cg.LinkRegister),
			}...,
		)
	}
	return label
}

func (cg *CodeGenerator) VisitStatSema(node ast.StatSema, ctx *ins.Context) ins.Instruction {
	expr := cg.VisitIdent(*node.GetIdent(), ctx)
	loadInstrs, reg, _ := cg.loadIfNotRegister(expr.Result, types.Word)
	cg.regs.ReleaseCallerSaved(reg)

	fName := "sem_wait"
	if node.IsUp() {
		fName = "sem_post"
	}
	return ins.Instructions{
		expr.Instrs,
		loadInstrs,
		ins.NewMove(reg, cg.ReturnRegister),
		ins.NewFunctionCall(fName),
	}
}
