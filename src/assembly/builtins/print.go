package builtins

import (
	architecture "wacc_32/assembly/architectures"
	ins "wacc_32/assembly/instructions"
	"wacc_32/types"
)

//This package should be used to store code needed to generate functions that are built in

var (
	returnReg ins.Register
	sp        ins.Register
	lr        ins.Register
	pc        ins.Register
	arg0      ins.Register
	arg1      ins.Register
	arg2      ins.Register
)

//Use this to initialise the builtin registers
func Init(conf architecture.Config) {
	returnReg = conf.ReturnRegister
	arg0 = 0
	arg1 = 1
	arg2 = 2
	sp = conf.StackPointer
	lr = conf.LinkRegister
	pc = conf.ProgramCounter

	PrintLine = ins.Instructions{
		ins.NewLabel(PrintLineLabel),
		ins.NewPush(lr),
		ins.NewLoad(ins.Variable(LineBSSLabel), arg0, types.Word),
		ins.NewAdd(returnReg, returnReg, ins.Immediate(types.Word)),
		ins.NewFunctionCall("puts"),
		ins.NewLoad(ins.NewDirective("streams"), returnReg, types.Word),
		ins.NewLoad(ins.NewAddress(returnReg), returnReg, types.Word),
		ins.NewFunctionCall("fflush"),
		ins.NewPop(pc),
	}
	printString = ins.Instructions{
		ins.NewLabel("p_print_string"),
		ins.NewPush(lr),
		ins.NewLoad(ins.NewAddress(returnReg), arg1, types.Word),
		ins.NewAdd(arg2, returnReg, ins.Immediate(types.Word)),
		ins.NewLoad(ins.Variable("print_string"), arg0, types.Word),
		ins.NewAdd(arg0, arg0, ins.Immediate(types.Word)),
		ins.NewFunctionCall("printf"),
		ins.NewLoad(ins.NewDirective("streams"), returnReg, types.Word),
		ins.NewLoad(ins.NewAddress(returnReg), returnReg, types.Word),
		ins.NewFunctionCall("fflush"),
		ins.NewPop(pc),
	}
	printBool = ins.Instructions{
		ins.NewLabel("p_print_bool"),
		ins.NewPush(lr),
		ins.NewMove(arg0, arg1),
		ins.NewLoad(ins.Variable("print_bool"), arg0, types.Word),
		ins.NewAdd(arg0, arg0, ins.Immediate(types.Word)),
		ins.NewCompare(arg1, ins.Immediate(1)),
		ins.NewBranch("p_true", ins.EQ),
		ins.NewAdd(returnReg, returnReg, ins.Immediate(falseOffset)),
		ins.NewLabel("p_true"),
		ins.NewFunctionCall("printf"),
		ins.NewLoad(ins.NewDirective("streams"), returnReg, types.Word),
		ins.NewLoad(ins.NewAddress(returnReg), returnReg, types.Word),
		ins.NewFunctionCall("fflush"),
		ins.NewPop(pc),
	}
}

//The labels associated with println
const (
	PrintLineLabel = "p_println"
	LineBSSLabel   = "println"
)

//These are used to print a newline
var (
	PrintLine ins.Instructions
	LineBSS   = ins.NewStringLiteral(LineBSSLabel, `"\0"`)
)

const falseOffset = 5

var (
	printString = ins.Instructions{}
	printBool   = ins.Instructions{}
)

//PrintType returns the print function and bss variable for a type
func PrintType(wt types.WaccType) (printIns, BSSIns ins.Instruction, printLabel, BSSLabel string) {
	isArr := wt.Is(types.Array)
	if isArr && wt.GetChildren()[0].Is(types.Char) {
		return printString, newBSSFormatString("print_", types.Str), "p_print_string", "print_string"
	}
	if isArr || wt.Is(types.Pair) {
		return printNumericalFunction("PTR"), ins.NewStringLiteral("print_PTR", "\"%p\\0\""), "p_print_PTR", "print_PTR"
	}
	if wt.Is(types.Str) {
		return printString, newBSSFormatString("print_", wt), "p_print_string", "print_string"
	}
	if wt.Is(types.Char) {
		return "", ins.NOOP{}, "putchar", ""
	}
	if wt.Is(types.Boolean) {
		return printBool, newBSSFormatString("print_", wt), "p_print_bool", "print_bool"
	}
	return printNumericalFunction(wt.String()), newBSSFormatString("print_", wt), "p_print_" + wt.String(), "print_" + wt.String()
}

func printNumericalFunction(typeString string) ins.Instructions {
	return ins.Instructions{
		ins.NewLabel("p_print_" + typeString),
		ins.NewPush(lr),
		ins.NewMove(arg0, arg1),
		ins.NewLoad(ins.Variable("print_"+typeString), returnReg, types.Word),
		ins.NewAdd(returnReg, returnReg, ins.Immediate(types.Word)),
		ins.NewFunctionCall("printf"),
		ins.NewLoad(ins.NewDirective("streams"), returnReg, types.Word),
		ins.NewLoad(ins.NewAddress(returnReg), returnReg, types.Word),
		ins.NewFunctionCall("fflush"),
		ins.NewPop(pc),
	}
}

func newBSSFormatString(header string, wt types.WaccType) ins.Instruction {
	return ins.NewStringLiteral(header+wt.String(), "\""+wt.GetFormatString()+"\\0\"")
}
