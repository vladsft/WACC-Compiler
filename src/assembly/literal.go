package assembly

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//VisitLiteral visits AST node ast.Literal
func (cg *CodeGenerator) VisitLiteral(node ast.Literal, ctx *ins.Context) ins.ExprInstr {
	var op ins.Operand
	if node.GetValue() == nil {
		return ins.NewExprInstr(ins.Instructions{}, ins.Immediate(0))
	}
	wt := node.EvalType(symboltable.SymbolTable{})
	if wt.Is(types.Array) {
		arr := node.GetValue().([]ast.Expression)
		size := int(types.TypeSize(wt.GetChildren()[0]))
		return cg.visitArrayLiteral(arr, size, ctx)
	}

	if wt.Is(types.UserDefinedType) {
		fieldTypes := wt.(types.UserType).GetFieldTypes()
		return cg.visitUserTypeLiteral(fieldTypes, node.GetValue().([]ast.Expression), ctx)
	}

	switch wt {
	case types.Integer:
		op = ins.Immediate(node.GetValue().(int))
	case types.Boolean:
		if node.GetValue().(bool) {
			op = ins.Immediate(1)
		} else {
			op = ins.Immediate(0)
		}
	case types.Char:
		op = ins.Immediate(int(parseEscapeChar(node.GetValue().(string)))) //Ignores the apostrophe
	case types.Str:
		id := cg.addStringToBSS(node.GetValue().(string))
		reg, err := cg.regs.AcquireCallerSaved()
		if err != nil {
			panic(err)
		}
		return ins.NewExprInstr(
			ins.Instructions{ins.NewLoad(ins.Variable(id), reg, types.Word)},
			reg,
		)
	case types.Sema:
		return cg.visitSemaLiteral(node.GetValue().(int))
	default:
		op = ins.Immediate(0)
	}
	return ins.NewExprInstr(ins.Instructions{}, op)
}

func (cg *CodeGenerator) visitArrayLiteral(arr []ast.Expression, size int, ctx *ins.Context) ins.ExprInstr {
	r0 := cg.ReturnRegister
	indexReg, _ := cg.regs.AcquireCallerSaved()
	defer cg.regs.ReleaseCallerSaved(indexReg)
	instrs := ins.Instructions{
		ins.NewMove(ins.Immediate(len(arr)), indexReg),
		ins.NewStoreHeap(indexReg, ins.Immediate(types.Word+len(arr)*size), 0, 0),
		ins.NewAdd(indexReg, r0, ins.Immediate(types.Word)),
	}
	for _, expr := range arr {
		exprInstr := cg.VisitExpression(expr, ctx)
		instrs = append(instrs, exprInstr.Instrs)
		eInstrs, reg, release := cg.loadIfNotRegister(exprInstr.Result, exprSize(expr))
		if release {
			cg.regs.ReleaseCallerSaved(reg)
		}
		instrs = append(instrs, eInstrs)
		instrs = append(instrs,
			ins.NewStore(types.Size(size), reg, ins.NewAddress(indexReg)),
			ins.NewAdd(indexReg, indexReg, ins.Immediate(size)),
		)
	}
	returnReg, _ := cg.regs.AcquireCallerSaved()
	return ins.NewExprInstr(append(instrs, ins.NewMove(r0, returnReg)), returnReg)
}

func (cg *CodeGenerator) visitUserTypeLiteral(fieldTypes []types.WaccType, exprs []ast.Expression, ctx *ins.Context) ins.ExprInstr {
	size := 0
	instrs := ins.Instructions{}

	// Calculating the size of the structure
	for _, t := range fieldTypes {
		tSize := types.TypeSize(t)
		size += int(tSize)
	}
	// MALLOC the struct object
	dummyReg, _ := cg.regs.AcquireCalleeSaved()
	instrs = append(instrs, ins.NewStoreHeap(dummyReg, ins.Immediate(size), 0, 0))
	cg.regs.ReleaseCalleeSaved(dummyReg)

	offset := 0

	for _, expr := range exprs {
		eSize := exprSize(expr)
		exprInstr := cg.VisitExpression(expr, ctx)
		loadInstrs, val, _ := cg.loadIfNotRegister(exprInstr.Result, eSize)
		cg.regs.ReleaseCallerSaved(val)
		instrs = append(instrs, exprInstr.Instrs)
		instrs = append(instrs, loadInstrs...)
		instrs = append(instrs, ins.NewStore(eSize, val, ins.NewAddress(cg.ReturnRegister, ins.Immediate(offset))))
		offset += int(eSize)
	}

	returnReg, _ := cg.regs.AcquireCallerSaved()
	return ins.NewExprInstr(append(instrs, ins.NewMove(cg.ReturnRegister, returnReg)), returnReg)
}

const semaSize = 16

func (cg *CodeGenerator) visitSemaLiteral(value int) ins.ExprInstr {
	r0, _ := cg.regs.AcquireCalleeSaved()
	r1, _ := cg.regs.AcquireCalleeSaved()
	r2, _ := cg.regs.AcquireCalleeSaved()
	defer cg.regs.ReleaseCalleeSaved(r2)
	defer cg.regs.ReleaseCalleeSaved(r1)
	defer cg.regs.ReleaseCalleeSaved(r0)

	reg, err := cg.regs.AcquireCallerSaved()
	if err != nil {
		panic(err)
	}
	return ins.NewExprInstr(
		ins.Instructions{
			ins.NewMove(ins.Immediate(semaSize), r0),
			ins.NewFunctionCall("malloc"),
			ins.NewStore(types.Word, r0, ins.NewAddress(cg.StackPointer)),
			ins.NewMove(ins.Immediate(value), r2),
			ins.NewMove(ins.Immediate(0), r1),
			ins.NewFunctionCall("sem_init"),
			ins.NewLoad(ins.NewAddress(cg.StackPointer), reg, types.Word),
		},
		reg,
	)
}
