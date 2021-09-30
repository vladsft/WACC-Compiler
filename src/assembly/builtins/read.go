package builtins

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/types"
)

//ReadType returns the scanf function and bss variable for a type
func ReadType(wt types.WaccType) (readIns, BSSIns ins.Instruction, readLabel, BSSLabel string) {
	if wt.Is(types.Str) {
		return readStringFunction(), newBSSFormatString("read_", wt), "p_read_string", "read_string"
	}
	return readNumericalFunction(wt.String()), newBSSFormatString("read_", wt), "p_read_" + wt.String(), "read_" + wt.String()

}

func readStringFunction() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel("p_read_string"),
		ins.NewPush(lr),
		ins.NewLoad(ins.NewAddress(returnReg), arg1, types.Word),
		ins.NewAdd(arg2, returnReg, ins.Immediate(types.Word)),
		ins.NewFunctionCall("scanf"),
		ins.NewMove(ins.Immediate(0), returnReg),
		ins.NewPop(pc),
	}
}

func readNumericalFunction(typeString string) ins.Instructions {
	return ins.Instructions{
		ins.NewLabel("p_read_" + typeString),
		ins.NewPush(lr),
		ins.NewMove(returnReg, arg1),
		ins.NewLoad(ins.Variable("read_"+typeString), returnReg, types.Word),
		ins.NewAdd(returnReg, returnReg, ins.Immediate(types.Word)),
		ins.NewFunctionCall("scanf"),
		ins.NewPop(pc),
	}
}
