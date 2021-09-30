package assembly

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/types"
)

//VisitExpression visits AST node ast.Expression
func (cg *CodeGenerator) VisitExpression(node ast.Expression, ctx *ins.Context) ins.ExprInstr {
	return node.(ast.ExprInstrAcceptor).AcceptExprInstr(cg, ctx)
}

//VisitRHS visits AST node ast.RHS
func (cg *CodeGenerator) VisitRHS(node ast.RHS, ctx *ins.Context) ins.ExprInstr {
	return cg.VisitExpression(node, ctx)
}

//VisitUnOp visits AST node ast.UnOp
func (cg *CodeGenerator) VisitUnOp(node ast.UnOp, ctx *ins.Context) ins.ExprInstr {
	expr := node.GetExpr()
	exprInstr := cg.VisitExpression(expr, ctx)
	var instr ins.Instruction
	switch node.GetOpType() {
	case ast.Not:
		notInstrs, reg, _ := cg.loadIfNotRegister(exprInstr.Result, types.Byte)
		instrs := append(notInstrs, ins.NewXor(reg, reg, ins.Immediate(1)))
		instrs = append(exprInstr.Instrs, instrs)
		return ins.NewExprInstr(instrs, reg)
	case ast.Len:
		reg, _ := cg.regs.AcquireCallerSaved()
		instrs := ins.Instructions{
			ins.NewLoad(exprInstr.Result, reg, exprSize(expr)),
			ins.NewLoad(ins.NewAddress(reg), reg, exprSize(expr)),
		}
		return ins.NewExprInstr(instrs, reg)
	case ast.Chr:
		return exprInstr
	case ast.Ord:
		loadInstrs, reg, _ := cg.loadIfNotRegister(exprInstr.Result, types.Byte)
		return ins.NewExprInstr(loadInstrs, reg)
	case ast.Neg:
		loadInstrs, reg, _ := cg.loadIfNotRegister(exprInstr.Result, types.Word)
		arg0, _ := cg.regs.AcquireCalleeSaved()
		defer cg.regs.ReleaseCalleeSaved(arg0)
		instrs := append(loadInstrs,
			ins.NewNeg(reg),
			ins.NewMove(reg, arg0),
			ins.NewFunctionCall("p_check_int_overflow"),
		)
		cg.addOverflowCode()
		return ins.NewExprInstr(instrs, reg)
	case ast.TryLock:
		loadInstrs, reg, _ := cg.loadIfNotRegister(exprInstr.Result, types.Word)
		r0, _ := cg.regs.AcquireCalleeSaved()
		defer cg.regs.ReleaseCalleeSaved(r0)
		instrs := append(loadInstrs,
			ins.NewMove(reg, r0),
			ins.NewFunctionCall("pthread_mutex_trylock"),
			ins.NewBoolean(reg, cg.ReturnRegister, ins.Immediate(16), ins.NE), //16 is EBUSY - indicates whether
		)
		return ins.NewExprInstr(instrs, reg)

	}
	return ins.NewExprInstr(append(exprInstr.Instrs, instr), exprInstr.Result)
}

//VisitRHSNewPair visits AST node ast.RHSNewPair
func (cg *CodeGenerator) VisitRHSNewPair(node ast.RHSNewPair, ctx *ins.Context) ins.ExprInstr {
	//Allocate elements
	instrs := ins.Instructions{}
	size := 0
	regs := [2]ins.Register{}
	var offset types.Size
	//Visit each pair element in turn and store a pointer to it in a register
	for i := 0; i < 2; i++ {
		expr := node.GetExpr(i)
		exprInstr := cg.VisitExpression(expr, ctx)
		eSize := exprSize(expr)
		if i == 0 {
			offset = eSize
		}
		size += int(eSize)
		instructions, val, release := cg.loadIfNotRegister(exprInstr.Result, eSize)
		if release {
			defer cg.regs.ReleaseCallerSaved(val)
		}
		instrs = append(instrs, exprInstr.Instrs...)
		instrs = append(instrs, instructions...)
		regs[i] = val
	}

	ptrReg, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(ptrReg)
	returnReg, _ := cg.regs.AcquireCallerSaved()
	//Malloc SIZE bytes and put the pair element pointers inside
	instrs = append(instrs,
		ins.NewStoreHeap(regs[0], ins.Immediate(size), 0, 0),
		ins.NewStore(offset, regs[1], ins.NewAddress(ptrReg, ins.Immediate(offset))),
		ins.NewMove(ptrReg, returnReg),
	)

	return ins.NewExprInstr(instrs, returnReg)
}

//VisitRHSFunctionCall visits AST node ast.RHSFunctionCall
func (cg *CodeGenerator) VisitRHSFunctionCall(node ast.RHSFunctionCall, ctx *ins.Context) ins.ExprInstr {
	args := node.GetArgs()
	instrs := make([]ins.Instruction, 0)

	name, _ := node.FormatName()
	offset := types.Word + node.GetSymbolTable().GetOffset(name)
	for _, arg := range args {
		argSize := exprSize(arg)

		negArgSize := ins.Immediate(-1 * offset)
		offset -= int(argSize)

		exprInstr := cg.VisitExpression(arg, ctx)
		instrs = append(instrs, exprInstr.Instrs)

		lInstrs, reg, _ := cg.loadIfNotRegister(exprInstr.Result, argSize)
		cg.regs.ReleaseCallerSaved(reg)

		instrs = append(instrs, lInstrs)
		instrs = append(instrs, ins.NewStore(argSize, reg, ins.NewAddress(cg.StackPointer, negArgSize)))

	}

	reg, _ := cg.regs.AcquireCallerSaved()
	instrs = append(instrs, ins.NewFunctionCall(name[1:]), ins.NewMove(cg.ReturnRegister, reg))

	return ins.ExprInstr{Instrs: instrs, Result: reg}
}

//VisitPairElem visits AST node ast.PairElem
func (cg *CodeGenerator) VisitPairElem(node ast.PairElem, ctx *ins.Context) ins.ExprInstr {
	offset := 0
	if node.GetPairElemPos() == ast.SND {
		pair := node.GetValue()
		pairType := pair.EvalType(*pair.GetSymbolTable())
		fstType := pairType.GetChildren()[0]
		offset = int(types.TypeSize(fstType))
	}
	expr := node.GetValue()
	pairExpr := cg.VisitExpression(expr, ctx)

	reg, _ := cg.regs.AcquireCallerSaved()
	cg.addNullPointerReferenceCode()

	r0, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(r0)
	return ins.NewExprInstr(ins.Instructions{
		ins.NewLoad(pairExpr.Result, reg, exprSize(expr)),
		ins.NewMove(reg, r0),
		ins.NewFunctionCall("p_check_null_pointer"),
	}, ins.NewAddress(reg, ins.Immediate(int(offset))))
}

func (cg *CodeGenerator) loadRegisterWithStack(op ins.Operand, size types.Size) (ins.Instructions, ins.Register, bool, bool) {
	reg, err := cg.regs.AcquireCallerSaved()
	stack := err != nil
	if err != nil {
		reg, _ = cg.regs.AcquireCallerSaved()
	}
	acquired := true
	switch operand := op.(type) {
	case ins.Address:
		acquired = false
		register := operand.Reg
		if register != cg.StackPointer {
			cg.regs.ReleaseCallerSaved(reg)
			reg = register
		} else {
			acquired = !stack
		}
		return ins.Instructions{ins.NewLoad(op, reg, size)}, reg, acquired, stack
	case ins.Immediate:
		return ins.Instructions{ins.NewMove(op, reg)}, reg, acquired, stack
	case ins.PseudoImmediate:
		return ins.Instructions{ins.NewMove(op, reg)}, reg, acquired, stack
	case ins.Variable:
		return ins.Instructions{ins.NewLoad(op, reg, types.Word)}, reg, acquired, stack
	case ins.Register:
		cg.regs.ReleaseCallerSaved(reg)
		return ins.Instructions{}, operand, false, false
	}
	return nil, -1, false, stack
}

func (cg *CodeGenerator) loadIfNotRegister(op ins.Operand, size types.Size) (ins.Instructions, ins.Register, bool) {
	reg, err := cg.regs.AcquireCallerSaved()
	if err != nil {
		panic(err)
	}
	switch operand := op.(type) {
	case ins.Address:
		acquired := false
		register := operand.Reg
		if register != cg.StackPointer && reg != cg.Accumulator {
			cg.regs.ReleaseCallerSaved(reg)
			reg = register
		} else {
			acquired = true
		}
		return ins.Instructions{ins.NewLoad(op, reg, size)}, reg, acquired
	case ins.Immediate:
		return ins.Instructions{ins.NewMove(op, reg)}, reg, true
	case ins.PseudoImmediate:
		return ins.Instructions{ins.NewMove(op, reg)}, reg, true
	case ins.Variable:
		return ins.Instructions{ins.NewLoad(op, reg, types.Word)}, reg, true
	case ins.Register:
		if reg != cg.Accumulator {
			cg.regs.ReleaseCallerSaved(reg)
		}
		return ins.Instructions{}, operand, false
	}
	return nil, -1, false
}

//VisitArrayElem visits AST node ast.ArrayElem
func (cg *CodeGenerator) VisitArrayElem(node ast.ArrayElem, ctx *ins.Context) ins.ExprInstr {
	indices := node.GetIndices()
	ident := cg.VisitIdent(*node.GetIdent(), ctx)
	loadInstrs, reg, _ := cg.loadIfNotRegister(ident.Result, exprSize(node.GetIdent()))
	size := types.TypeSize(node.EvalType(*node.GetSymbolTable()))

	cg.addOutOfBoundsCode()

	r0, _ := cg.regs.AcquireCalleeSaved()
	r1, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(r0)
	defer cg.regs.ReleaseCalleeSaved(r1)
	sizeReg, err := cg.regs.AcquireCallerSaved()
	if err != nil {
		panic(err)
	}
	defer cg.regs.ReleaseCallerSaved(sizeReg)
	instrs := append(ident.Instrs, loadInstrs...)
	instrs = append(instrs, ins.NewLoad(ins.NewAddress(reg), r1, types.Word))
	//Skip the first word (length)
	for _, expr := range indices {
		index := cg.VisitExpression(expr, ctx) //Result will always be a register or an immediate
		loadIndex, iReg, _ := cg.loadIfNotRegister(index.Result, exprSize(expr))
		cg.regs.ReleaseCallerSaved(iReg)

		instrs = append(instrs,
			ins.NewAdd(reg, reg, ins.Immediate(types.Word)),
		)
		instrs = append(instrs, index.Instrs...)
		instrs = append(instrs, loadIndex...)
		instrs = append(instrs,
			ins.NewMove(ins.Immediate(int(size)), sizeReg),
			ins.NewMult(iReg, iReg, sizeReg),
			ins.NewMove(iReg, r0),
			ins.NewMult(r1, r1, sizeReg),
			ins.NewFunctionCall("p_check_array_index_out_of_bounds"),
			ins.NewAdd(reg, reg, iReg),
			ins.NewLoad(ins.NewAddress(reg), reg, size),
		)
	}
	return ins.NewExprInstr(instrs[:len(instrs)-1], ins.NewAddress(reg))
}

//VisitIdent visits AST node ast.Ident
func (cg *CodeGenerator) VisitIdent(node ast.Ident, ctx *ins.Context) ins.ExprInstr {
	if node.IsNamespaced() {
		return cg.VisitFieldAccess(node, ctx)
	}
	offset := ins.Immediate(node.GetSymbolTable().GetOffset(node.String()))
	return ins.NewExprInstr(nil, ins.NewAddress(cg.StackPointer, offset))
}

func (cg *CodeGenerator) VisitMake(node ast.Make, ctx *ins.Context) ins.ExprInstr {
	exprInstr := cg.VisitExpression(node.GetLengthExpression(), ctx)
	loadInstrs, reg, _ := cg.loadIfNotRegister(exprInstr.Result, types.Word)

	typeSize := types.TypeSize(node.EvalType(*node.GetSymbolTable()))
	sizeReg, _ := cg.regs.AcquireCallerSaved()
	defer cg.regs.ReleaseCallerSaved(sizeReg)

	instrs := exprInstr.Instrs
	instrs = append(instrs, loadInstrs...)
	instrs = append(instrs, []ins.Instruction{
		ins.NewMove(ins.Immediate(typeSize), sizeReg),
		ins.NewMult(cg.ReturnRegister, reg, sizeReg),
		ins.NewAdd(cg.ReturnRegister, cg.ReturnRegister, ins.Immediate(types.Word)),
		ins.NewStoreHeap(reg, cg.ReturnRegister, 0, 0), //Store the length on the heap
		ins.NewMove(cg.ReturnRegister, reg),
	}...)
	return ins.NewExprInstr(instrs, reg)
}
