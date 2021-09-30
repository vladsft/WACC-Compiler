package assembly

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//Use this to add sandwich the instructions with head and tail
func wrapInstructions(head ins.Instruction, instrs ins.Instructions, tail ins.Instruction) ins.Instructions {
	withHead := append(ins.Instructions{head}, instrs)
	return append(withHead, tail)
}

//Use this to add sandwich the instructions with head and tail
func addHeader(head ins.Instruction, instrs ins.Instructions) ins.Instructions {
	return append(ins.Instructions{head}, instrs)
}

//Use this to add sandwich the instructions with head and tail
func addTail(instrs ins.Instructions, tail ins.Instruction) ins.Instructions {
	return append(instrs, tail)
}

func exprSize(e ast.Expression) types.Size {
	if table := e.GetSymbolTable(); table != nil {
		return types.TypeSize(e.EvalType(*table))
	}
	return types.TypeSize(e.EvalType(symboltable.SymbolTable{}))
}

var escapeCodes = map[byte]byte{
	'\\': '\\',
	't':  '\t',
	'n':  '\n',
	'b':  '\b',
	'f':  '\f',
	'r':  '\r',
	'"':  '"',
	'\'': '\'',
	'0':  byte(0),
}

//Char is a string of the form '\?x'
func parseEscapeChar(char string) byte {
	if char[1] == '\\' {
		return escapeCodes[char[2]]
	}
	return char[1]
}
