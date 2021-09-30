package assembly

import (
	"strconv"
	"wacc_32/assembly/builtins"
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/symboltable"
	"wacc_32/types"
)

func (cg *CodeGenerator) visitStatsTillEnd(stats []ast.Statement, ctx *ins.Context) (instrs []ins.Instruction, early bool) {
	instrs = make([]ins.Instruction, 0, len(stats))
	length := 0
	for _, stat := range stats {
		length++
		instrs = append(instrs, cg.VisitAST(ast.AST(stat), ctx))
		switch stat.(type) {
		case *ast.StatReturn:
			early = true
		case *ast.StatExit:
			early = true
		default:
			continue
		}
	}
	return instrs[:length], early
}

//VisitAST visits AST node ast.AST
func (cg *CodeGenerator) VisitAST(node ast.AST, ctx *ins.Context) ins.Instruction {
	return node.(ast.InstructionAcceptor).AcceptInstruction(cg, ctx)
}

//VisitStatement visits AST node ast.Statement
func (cg *CodeGenerator) VisitStatement(node ast.Statement, ctx *ins.Context) ins.Instruction {
	return cg.VisitAST(node, ctx)
}

//VisitStatSkip visits AST node ast.StatSkip
func (cg *CodeGenerator) VisitStatSkip(node ast.StatSkip, ctx *ins.Context) ins.Instruction {
	return ins.NOOP{}
}

//VisitStatRead visits AST node ast.StatRead
func (cg *CodeGenerator) VisitStatRead(node ast.StatRead, ctx *ins.Context) ins.Instruction {
	toRead := node.GetToRead()
	toReadType := toRead.EvalType(*node.GetSymbolTable())
	exprInstr := cg.VisitExpression(toRead, ctx)
	instrs := exprInstr.Instrs

	rInstrs, bss, readLabel, bssLabel := builtins.ReadType(toReadType)
	cg.bssVars[bssLabel] = bss
	cg.funcs[readLabel] = rInstrs

	addr := exprInstr.Result.(ins.Address)

	instrs = append(instrs,
		ins.NewAdd(cg.ReturnRegister, ins.Register(cg.StackPointer), ins.Immediate(addr.Offset)),
		ins.NewFunctionCall(readLabel),
	)

	return ins.Instructions(instrs)
}

//VisitStatFree visits AST node ast.StatFree
func (cg *CodeGenerator) VisitStatFree(node ast.StatFree, ctx *ins.Context) ins.Instruction {
	exprInstr := cg.VisitExpression(node.GetExpression(), ctx)

	cg.addNullPointerReferenceCode()

	instrs := ins.Instructions{
		exprInstr.Instrs,
		ins.NewMove(exprInstr.Result, cg.ReturnRegister),
		ins.NewFunctionCall("p_check_null_pointer"),
	}
	//When freeing a lock, call pthread_mutex_destroy
	if node.GetExpression().EvalType(*node.GetSymbolTable()).Is(types.Lock) {
		instrs = append(instrs,
			ins.NewMove(exprInstr.Result, cg.ReturnRegister),
			ins.NewFunctionCall("pthread_mutex_destroy"),
		)
	}

	return append(instrs, ins.NewFunctionCall(builtins.FreeLabel))
}

type printer interface {
	GetExprToPrint() ast.Expression
	GetSymbolTable() *symboltable.SymbolTable
}

func (cg *CodeGenerator) visitPrinter(node printer, ctx *ins.Context) ins.Instructions {
	toPrint := node.GetExprToPrint()
	expr := cg.VisitExpression(toPrint, ctx)
	instrs := expr.Instrs
	printType := toPrint.EvalType(*node.GetSymbolTable())
	pInstrs, bss, printLabel, bssLabel := builtins.PrintType(printType)
	cg.bssVars[bssLabel] = bss
	cg.funcs[printLabel] = pInstrs

	rInstrs, result, _ := cg.loadIfNotRegister(expr.Result, exprSize(toPrint))
	cg.regs.ReleaseCallerSaved(result)
	arg0, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(arg0)
	instrs = append(instrs,
		rInstrs,
		ins.NewMove(result, arg0),
		ins.NewFunctionCall(printLabel),
	)

	return ins.Instructions(instrs)
}

//VisitStatPrint visits AST node ast.StatPrint
func (cg *CodeGenerator) VisitStatPrint(node ast.StatPrint, ctx *ins.Context) ins.Instruction {
	return cg.visitPrinter(node, ctx)
}

//VisitStatPrintln visits AST node ast.StatPrintln
func (cg *CodeGenerator) VisitStatPrintln(node ast.StatPrintln, ctx *ins.Context) ins.Instruction {
	instrs := cg.visitPrinter(node, ctx)
	cg.bssVars[builtins.LineBSSLabel] = builtins.LineBSS
	cg.funcs[builtins.PrintLineLabel] = builtins.PrintLine
	return addTail(instrs, ins.NewFunctionCall("p_println"))
}

//VisitStatExit visits AST node ast.StatExit
func (cg *CodeGenerator) VisitStatExit(node ast.StatExit, ctx *ins.Context) ins.Instruction {
	expr := cg.VisitExpression(node.GetCode(), ctx)
	return ins.Instructions(append(expr.Instrs, ins.NewExit(expr.Result)))
}

//VisitStatFor visits AST node ast.StatFor
func (cg *CodeGenerator) VisitStatFor(node ast.StatFor, ctx *ins.Context) ins.Instruction {
	label := "for_" + strconv.Itoa(cg.numFors)
	initialLabel := "initial_" + label
	condLabel := "cond_" + label
	//changeLabel := "change_" + label
	endLabel := "for_end_" + label
	cg.numFors++

	initialStat := cg.VisitStatNewassign(node.GetInitial(), ctx)
	condInstrs := cg.VisitExpression(node.GetCond(), ctx)
	condLoad, condReg, _ := cg.loadIfNotRegister(condInstrs.Result, types.Byte)
	cg.regs.ReleaseCallerSaved(condReg)
	changeStat := cg.VisitStatAssign(node.GetChange(), ctx)

	bodyInstrs := cg.VisitStatement(node.GetBody(), ctx)

	return ins.Instructions{
		ins.NewLabel(initialLabel),
		initialStat,
		ins.NewLabel(condLabel),
		ins.Instructions(condInstrs.Instrs),
		condLoad,
		ins.NewCompare(condReg, ins.Immediate(1)),
		ins.NewBranch(endLabel, ins.NE),
		bodyInstrs,
		changeStat,
		ins.NewBranch(condLabel, ins.AL),
		ins.NewLabel(endLabel),
	}
}

//VisitStatWhile visits AST node ast.StatWhile
func (cg *CodeGenerator) VisitStatWhile(node ast.StatWhile, ctx *ins.Context) ins.Instruction {
	label := "while_" + strconv.Itoa(cg.numWhiles)
	condLabel := "cond_" + label
	endLabel := "while_end_" + label
	cg.numWhiles++

	condInstrs := cg.VisitExpression(node.GetCond(), ctx)
	condLoad, condReg, _ := cg.loadIfNotRegister(condInstrs.Result, types.Byte)
	cg.regs.ReleaseCallerSaved(condReg)

	bodyInstrs := cg.VisitStatement(node.GetBody(), ctx)

	return ins.Instructions{
		ins.NewLabel(condLabel),
		ins.Instructions(condInstrs.Instrs),
		condLoad,
		ins.NewCompare(condReg, ins.Immediate(1)),
		ins.NewBranch(endLabel, ins.NE),
		bodyInstrs,
		ins.NewBranch(condLabel, ins.AL),
		ins.NewLabel(endLabel),
	}
}

//VisitStatWhile visits AST node ast.StatWhile
func (cg *CodeGenerator) VisitStatDoWhile(node ast.StatDoWhile, ctx *ins.Context) ins.Instruction {
	label := "do_while_" + strconv.Itoa(cg.numDoWhiles)
	condLabel := "cond_" + label
	endLabel := "do_while_end_" + label
	cg.numDoWhiles++

	bodyInstrs := cg.VisitStatement(node.GetBody(), ctx)

	condInstrs := cg.VisitExpression(node.GetCond(), ctx)
	condLoad, condReg, _ := cg.loadIfNotRegister(condInstrs.Result, types.Byte)
	cg.regs.ReleaseCallerSaved(condReg)

	return ins.Instructions{
		// ins.NewLabel(condLabel),
		// ins.Instructions(condInstrs.Instrs),
		// condLoad,
		// ins.NewCompare(condReg, ins.Immediate(1)),
		// ins.NewBranch(endLabel, ins.NE),
		// bodyInstrs,
		// ins.NewBranch(condLabel, ins.AL),
		// ins.NewLabel(endLabel),
		ins.NewLabel(condLabel),
		bodyInstrs,
		ins.Instructions(condInstrs.Instrs),
		condLoad,
		ins.NewCompare(condReg, ins.Immediate(1)),
		ins.NewBranch(endLabel, ins.NE),
		ins.NewBranch(condLabel, ins.AL),
		ins.NewLabel(endLabel),
	}
}

//VisitStatBegin visits AST node ast.StatBegin
func (cg *CodeGenerator) VisitStatBegin(node ast.StatBegin, ctx *ins.Context) ins.Instruction {
	return cg.VisitStatement(node.GetStat(), ctx)
}

//VisitStatNewassign visits AST node ast.StatNewassign
func (cg *CodeGenerator) VisitStatNewassign(node ast.StatNewassign, ctx *ins.Context) ins.Instruction {
	size := types.TypeSize(node.GetType())

	rhs := node.GetRHS()
	exprInstr := cg.VisitRHS(rhs, ctx)

	eInstrs, eReg, _ := cg.loadIfNotRegister(exprInstr.Result, exprSize(node.GetRHS()))
	defer cg.regs.ReleaseCallerSaved(eReg) //Always release the right hand side register
	offset := ins.Immediate(node.GetSymbolTable().GetOffset(node.GetName()))

	var lockInstrs ins.Instructions
	const lockSize = 6 * types.Word

	//If we have a lock, create an error checking lock
	if node.GetType().Is(types.Lock) {
		r0, _ := cg.regs.AcquireCalleeSaved()
		r1, _ := cg.regs.AcquireCalleeSaved()
		scratch, _ := cg.regs.AcquireCallerSaved()
		defer cg.regs.ReleaseCalleeSaved(r0)
		defer cg.regs.ReleaseCalleeSaved(r1)
		lockInstrs = ins.Instructions{
			ins.NewStoreHeap(r0, ins.Immediate(lockSize), 0, 0), //create mutex
			ins.NewMove(r0, eReg),
			ins.NewStoreHeap(r0, ins.Immediate(types.Word), 0, 0), //create mutexattr
			ins.NewMove(ins.Immediate(2), r1),                     //2 is PTHREAD_MUTEX_ERRORCHECK_NP
			ins.NewMove(r0, scratch),
			ins.NewFunctionCall("pthread_mutexattr_settype"),
			ins.NewMove(scratch, r1),
			ins.NewMove(eReg, r0),
			ins.NewFunctionCall("pthread_mutex_init"),
		}
	}

	return ins.Instructions{
		exprInstr.Instrs,
		eInstrs,
		lockInstrs,
		ins.NewStore(size, eReg, ins.NewAddress(cg.StackPointer, offset)),
	}
}

//VisitStatAssign visits AST node ast.StatAssign
func (cg *CodeGenerator) VisitStatAssign(node ast.StatAssign, ctx *ins.Context) ins.Instruction {
	defer recover()

	lhs := node.GetLHS()
	rhs := node.GetRHS()
	rExpr := cg.VisitExpression(rhs, ctx)
	rInstrs, regOp, _ := cg.loadIfNotRegister(rExpr.Result, exprSize(rhs))

	lExpr := cg.VisitExpression(lhs, ctx) //This will always be an address

	lReg := lExpr.Result.(ins.Address).Reg
	switch lhs.(type) {
	case *ast.PairElem:
		defer cg.regs.ReleaseCallerSaved(lReg)
	case *ast.ArrayElem:
		defer cg.regs.ReleaseCallerSaved(lReg)
	}

	cg.regs.ReleaseCallerSaved(regOp) //Always release the rhs register
	return append(rExpr.Instrs, rInstrs, lExpr.Instrs,
		ins.NewStore(types.TypeSize(lhs.EvalType(*node.GetSymbolTable())), regOp, lExpr.Result),
	)
}

//VisitStatEnhancedAssign visits AST node ast.StatEnhancedAssign
func (cg *CodeGenerator) VisitStatEnhancedAssign(node ast.StatEnhancedAssign, ctx *ins.Context) ins.Instruction {

	ident := node.GetIdent()
	rhs := node.GetRHS()
	exprInstr := cg.VisitRHS(rhs, ctx)
	lExpr := cg.VisitExpression(ident, ctx)

	eInstrs, eReg, _ := cg.loadIfNotRegister(exprInstr.Result, exprSize(node.GetRHS()))
	defer cg.regs.ReleaseCallerSaved(eReg)

	return ins.Instructions{
		exprInstr.Instrs,
		eInstrs,
		ins.NewStore(types.TypeSize(ident.EvalType(*node.GetSymbolTable())), eReg, lExpr.Result),
	}
}

//VisitStatIf visits AST node ast.StatIf
func (cg *CodeGenerator) VisitStatIf(node ast.StatIf, ctx *ins.Context) ins.Instruction {
	cond := node.GetCondition()

	condExpr := cg.VisitExpression(cond, ctx)

	switch condType := cond.(type) {
	case *ast.Literal:
		if condType.GetValue().(bool) {
			return cg.VisitStatement(node.GetIfStat(), ctx)
		}
		return cg.VisitStatement(node.GetElseStat(), ctx)
	}

	loadInstrs, reg, _ := cg.loadIfNotRegister(condExpr.Result, types.Byte)
	cg.regs.ReleaseCallerSaved(reg)

	instrs := condExpr.Instrs
	ifStat := cg.VisitStatement(node.GetIfStat(), ctx)
	elseStat := cg.VisitStatement(node.GetElseStat(), ctx)

	nIfs := strconv.Itoa(cg.nIfs)
	elseLabel := "else_" + nIfs
	endLabel := "end_" + nIfs

	cg.nIfs++
	return ins.Instructions{
		instrs,
		loadInstrs,
		ins.NewCompare(reg, ins.Immediate(0)),
		ins.NewBranch(elseLabel, ins.EQ),
		ifStat,
		ins.NewBranch(endLabel, ins.AL),
		ins.NewLabel(elseLabel),
		elseStat,
		ins.NewLabel(endLabel),
	}
}

//VisitStatMultiple visits AST node ast.StatMultiple
func (cg *CodeGenerator) VisitStatMultiple(node ast.StatMultiple, ctx *ins.Context) ins.Instruction {
	instrs := ins.Instructions{}
	for _, stat := range node {
		instrs = append(instrs, cg.VisitStatement(stat, ctx))
	}
	return instrs
}
