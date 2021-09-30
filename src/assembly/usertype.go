package assembly

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/types"
)

func (cg *CodeGenerator) VisitUserType(node ast.UserType, ctx *ins.Context) ins.Instruction {
	return ins.NOOP{}
}

func (cg *CodeGenerator) VisitFieldAccess(node ast.Ident, ctx *ins.Context) ins.ExprInstr {
	components := node.GetNameComponents()
	sym := node.GetSymbolTable()
	reg, _ := cg.regs.AcquireCallerSaved()

	loadInstrs := make([]ins.Instruction, len(components))
	t, _ := sym.GetType(components[0])
	loadInstrs[0] = ins.NewLoad(
		ins.NewAddress(cg.StackPointer, ins.Immediate(sym.GetOffset(components[0]))),
		reg, types.Word,
	)

	for i := 1; i < len(components); i++ {
		fieldName := components[i]
		uType := t.(types.UserType)
		uType, _ = ast.LookupUserType(uType, *sym)
		offset := 0
		fieldNames := uType.GetFieldNames()
		fieldTypes := uType.GetFieldTypes()

		for j := 0; j < len(fieldNames); j++ {
			if fieldNames[j] == fieldName {
				t = fieldTypes[j]
				break
			}
			offset += int(types.TypeSize(fieldTypes[j]))
		}
		loadInstrs[i] = ins.NewLoad(ins.NewAddress(reg, ins.Immediate(offset)), reg, types.Word)
	}
	addr := loadInstrs[len(loadInstrs)-1].(ins.Load).Src
	return ins.NewExprInstr(loadInstrs[:len(loadInstrs)-1], addr)
}
