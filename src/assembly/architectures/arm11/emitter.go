package arm11

import (
	"fmt"
	"strconv"
	"strings"
	architecture "wacc_32/assembly/architectures"
	ins "wacc_32/assembly/instructions"
	"wacc_32/types"
)

var _ architecture.Emitter = Emitter{}

type Emitter struct{}

const (
	arm11ImmSize = 1024
	addStr       = "\tadds %s, %s, %s"
	subStr       = "\tsubs %s, %s, %s"
)

func splitOperand(op ins.Operand) []ins.Operand {
	var val int
	var newOp func(int) ins.Operand
	switch operand := op.(type) {
	case ins.Immediate:
		val = int(operand)
		newOp = func(n int) ins.Operand {
			return ins.Immediate(n)
		}
	case ins.PseudoImmediate:
		val = int(operand)
		newOp = func(n int) ins.Operand {
			return ins.PseudoImmediate(n)
		}
	default:
		return []ins.Operand{op}
	}
	nOps := (val + 1) / arm11ImmSize
	ops := make([]ins.Operand, nOps+1)
	for i := 0; i < nOps; i++ {
		ops[i] = newOp(arm11ImmSize)
	}
	ops[nOps] = newOp(val % arm11ImmSize)
	return ops
}

func (arm Emitter) Emit(bss, instr ins.Instruction) string {
	codeString := ".data\n"
	codeString += arm.EmitInstruction(bss)
	codeString += "\n.text\n.global main\n"
	codeString += ".streams:\n\t.word stdout\n"
	codeString += arm.EmitInstruction(instr)
	return codeString + "\n"
}

func (arm Emitter) EmitInstruction(instruction ins.Instruction) string {
	switch instr := instruction.(type) {
	case ins.Instructions:
		return arm.EmitInstructions(instr)
	case ins.NOOP:
		return arm.EmitNOOP(instr)
	case ins.Pool:
		return arm.EmitPool(instr)
	case ins.Exit:
		return arm.EmitExit(instr)
	case ins.Label:
		return arm.EmitLabel(instr)
	case ins.FunctionCall:
		return arm.EmitFunctionCall(instr)
	case ins.Branch:
		return arm.EmitBranch(instr)
	case ins.Move:
		return arm.EmitMove(instr)
	case ins.Compare:
		return arm.EmitCompare(instr)
	case ins.Add:
		return arm.EmitAdd(instr)
	case ins.Sub:
		return arm.EmitSub(instr)
	case ins.Mult:
		return arm.EmitMult(instr)
	case ins.Div:
		return arm.EmitDiv(instr)
	case ins.Mod:
		return arm.EmitMod(instr)
	case ins.Xor:
		return arm.EmitXor(instr)
	case ins.StringLiteral:
		return arm.EmitStringLiteral(instr)
	case ins.BoolExpr:
		return arm.EmitBoolExpr(instr)
	case ins.And:
		return arm.EmitAnd(instr)
	case ins.Or:
		return arm.EmitOr(instr)
	case ins.Neg:
		return arm.EmitNeg(instr)
	case types.Size:
		return arm.EmitSize(instr)
	case ins.Load:
		return arm.EmitLoad(instr)
	case ins.StackInstr:
		return arm.EmitStackInstr(instr)
	case ins.Store:
		return arm.EmitStore(instr)
	case ins.StoreHeap:
		return arm.EmitStoreHeap(instr)
	case ins.FreeHeap:
		return arm.EmitFreeHeap(instr)
	case ins.Operand:
		return arm.EmitOperand(instr)
	case ins.Immediate:
		return arm.EmitImmediate(instr)
	case ins.PseudoImmediate:
		return arm.EmitPseudoImmediate(instr)
	case ins.Address:
		return arm.EmitAddress(instr)
	case ins.Variable:
		return arm.EmitVariable(instr)
	case ins.Register:
		return arm.EmitRegister(instr)
	}
	return ""
}

func (arm Emitter) EmitInstructions(is ins.Instructions) string {
	strs := make([]string, len(is))
	for i, instr := range is {
		strs[i] = arm.EmitInstruction(instr)
	}
	return strings.Join(strs, "\n")

}
func (arm Emitter) EmitNOOP(_ ins.NOOP) string {
	return ""
}
func (arm Emitter) EmitExit(e ins.Exit) string {
	var code string
	switch op := e.Code.(type) {
	case ins.Address:
		code = "\tldr r0, " + arm.EmitAddress(op)
	default:
		code = "\tmov r0, " + arm.EmitOperand(e.Code)
	}
	return code + "\n\tbl exit"
}

func (arm Emitter) EmitLabel(l ins.Label) string {
	return l.Name + ":"
}

func (arm Emitter) EmitFunctionCall(fc ins.FunctionCall) string {
	return "\tbl " + fc.Name
}

func (arm Emitter) EmitBranch(b ins.Branch) string {
	return fmt.Sprintf("\tb%s %s", b.Condition, b.Label)
}

func (arm Emitter) EmitMove(m ins.Move) string {
	instr := "mov"
	srcOp := m.Src
	switch src := m.Src.(type) {
	case ins.Immediate:
		srcOp = ins.PseudoImmediate(src)
		instr = "ldr"
	case ins.Address:
		instr = "ldr"
	}
	return fmt.Sprintf("\t%s %s, %s", instr, arm.EmitRegister(m.Dest), arm.EmitOperand(srcOp))
}

func (arm Emitter) EmitCompare(c ins.Compare) string {
	return fmt.Sprintf("\tcmp %s, %s", arm.EmitOperand(c.Left), arm.EmitOperand(c.Right))
}

func (arm Emitter) EmitAdd(a ins.Add) string {
	lOps := splitOperand(a.Left)
	rOps := splitOperand(a.Right)

	strs := make([]string, len(lOps)+len(rOps)-1)
	destStr := arm.EmitRegister(a.Dest)
	strs[0] = fmt.Sprintf(addStr, destStr, arm.EmitOperand(lOps[0]), arm.EmitOperand(rOps[0]))

	for i := 1; i < len(lOps); i++ {
		strs[i] = fmt.Sprintf(addStr, destStr, destStr, arm.EmitOperand(lOps[i]))
	}

	for i := len(lOps); i < len(strs); i++ {
		strs[i] = fmt.Sprintf(addStr, destStr, destStr, arm.EmitOperand(rOps[i+1-len(lOps)]))
	}

	return strings.Join(strs, "\n")
}

func (arm Emitter) EmitSub(s ins.Sub) string {
	lOps := splitOperand(s.Left)
	rOps := splitOperand(s.Right)

	strs := make([]string, len(lOps)+len(rOps)-1)
	destStr := arm.EmitRegister(s.Dest)
	strs[0] = fmt.Sprintf(subStr, destStr, arm.EmitOperand(lOps[0]), arm.EmitOperand(rOps[0]))

	for i := 1; i < len(lOps); i++ {
		strs[i] = fmt.Sprintf(subStr, destStr, destStr, arm.EmitOperand(lOps[i]))
	}

	for i := len(lOps); i < len(strs); i++ {
		strs[i] = fmt.Sprintf(subStr, destStr, destStr, arm.EmitOperand(rOps[i+1-len(lOps)]))
	}

	return strings.Join(strs, "\n")
}

func (arm Emitter) EmitMult(m ins.Mult) string {
	l, r := arm.EmitOperand(m.Left), arm.EmitOperand(m.Right)
	d := arm.EmitRegister(m.Dest)
	return fmt.Sprintf("\tsmull %s, r12, %s, %s\n\tcmp r12, %s, ASR #31\n\tcmpNE r0, #1<<31\n\tcmnVC r0, #1<<31", d, l, r, d)
}

func (arm Emitter) EmitDiv(d ins.Div) string {
	reg0Line := fmt.Sprintf("\tmov r0, %s\n", arm.EmitOperand(d.Left))
	reg1Line := fmt.Sprintf("\tmov r1, %s\n", arm.EmitOperand(d.Right))
	divByZero := fmt.Sprintf("\tbl p_check_divide_by_zero\n")
	branchLine := fmt.Sprintf("\tbl __aeabi_idiv\n")
	moveBack := fmt.Sprintf("\tmov %s, r0", arm.EmitRegister(d.Dest))
	return reg0Line + reg1Line + divByZero + branchLine + moveBack
}

func (arm Emitter) EmitMod(m ins.Mod) string {
	reg0Line := fmt.Sprintf("\tmov r0, %s\n", arm.EmitOperand(m.Left))
	reg1Line := fmt.Sprintf("\tmov r1, %s\n", arm.EmitOperand(m.Right))
	divByZero := fmt.Sprintf("\tbl p_check_divide_by_zero\n")
	branchLine := fmt.Sprintf("\tbl __aeabi_idivmod\n")
	moveBack := fmt.Sprintf("\tmov %s, r1", arm.EmitRegister(m.Dest))
	return reg0Line + reg1Line + divByZero + branchLine + moveBack
}

func (arm Emitter) EmitXor(x ins.Xor) string {
	return fmt.Sprintf("\teor %s, %s, %s", arm.EmitRegister(x.Dest), arm.EmitOperand(x.Left), arm.EmitOperand(x.Right))
}
func (arm Emitter) EmitStringLiteral(sl ins.StringLiteral) string {
	msgLine := fmt.Sprintf("%s:\n", sl.ID)
	wordLine := fmt.Sprintf("\t.word %d\n", sl.Size)
	asciiLine := fmt.Sprintf("\t.ascii %s", sl.String)

	return msgLine + wordLine + asciiLine
}
func (arm Emitter) EmitBoolExpr(be ins.BoolExpr) string {
	cmpLine := fmt.Sprintf("\tcmp %s, %s\n", arm.EmitOperand(be.Left), arm.EmitOperand((be.Right)))
	movTrue := fmt.Sprintf("\tmov%s %s, #1\n", be.True.String(), arm.EmitRegister(be.Dest))
	movFalse := fmt.Sprintf("\tmov%s %s, #0", be.False.String(), arm.EmitRegister(be.Dest))

	return cmpLine + movTrue + movFalse
}

func (arm Emitter) EmitAnd(a ins.And) string {
	return fmt.Sprintf("\tand %s, %s, %s", arm.EmitRegister(a.Dest), arm.EmitOperand(a.Left), arm.EmitOperand(a.Right))
}

func (arm Emitter) EmitOr(o ins.Or) string {
	return fmt.Sprintf("\torr %s, %s, %s", arm.EmitRegister(o.Dest), arm.EmitOperand(o.Left), arm.EmitOperand(o.Right))
}

func (arm Emitter) EmitNeg(n ins.Neg) string {
	return fmt.Sprintf("\tnegs %s, %s", arm.EmitRegister(n.Reg), arm.EmitRegister(n.Reg))
}

func (arm Emitter) EmitSize(types.Size) string {
	return ""
}

func (arm Emitter) EmitLoad(ld ins.Load) string {
	return fmt.Sprintf("\tldr%s %s, %s", ld.Size, arm.EmitRegister(ld.Dest), arm.EmitOperand(ld.Src))
}

func (arm Emitter) EmitStackInstr(st ins.StackInstr) string {
	regs := make([]string, len(st.Regs))
	for i, reg := range st.Regs {
		regs[i] = arm.EmitRegister(reg)
	}
	return fmt.Sprintf("\t%s {%s}", st.T, strings.Join(regs, ", "))
}

func (arm Emitter) EmitStore(st ins.Store) string {
	return fmt.Sprintf("\tstr%s %s, %s", st.Size, arm.EmitRegister(st.Src), arm.EmitOperand(st.Dest))
}

func (arm Emitter) EmitStoreHeap(st ins.StoreHeap) string {
	mov := "\tmov r0, " + arm.EmitOperand(st.Size) + "\n"
	bl := mov + fmt.Sprintf("\tbl malloc\n")
	addr := "[r0"
	if st.Offset > 0 {
		addr += ", " + strconv.Itoa(st.Offset)
		if st.Multiplier > 0 {
			addr += ", " + strconv.Itoa(st.Multiplier)
		}
	}
	return bl + fmt.Sprintf("\tstr %s, %s]", arm.EmitOperand(st.Op), addr)
}
func (arm Emitter) EmitFreeHeap(fh ins.FreeHeap) string {
	mov := fmt.Sprintf("\tmov r0, %s", arm.EmitRegister(fh.Reg))
	return mov + "\n\tbl free"

}

func (arm Emitter) EmitOperand(op ins.Operand) string {
	switch operand := op.(type) {
	case ins.Address:
		return arm.EmitAddress(operand)
	case ins.PseudoImmediate:
		return arm.EmitPseudoImmediate(operand)
	case ins.Register:
		return arm.EmitRegister(operand)
	case ins.Immediate:
		return arm.EmitImmediate(operand)
	case ins.Variable:
		return arm.EmitVariable(operand)
	case ins.FunctionPointer:
		return arm.EmitFunctionPointer(operand)
	case ins.Directive:
		return arm.EmitDirective(operand)
	}
	return ""
}

func (arm Emitter) EmitImmediate(i ins.Immediate) string {
	return "#" + strconv.Itoa(int(i))
}
func (arm Emitter) EmitPseudoImmediate(pi ins.PseudoImmediate) string {
	return "=" + strconv.Itoa(int(pi))
}
func (arm Emitter) EmitAddress(a ins.Address) string {
	str := "[" + arm.EmitRegister(a.Reg)
	if a.Offset != 0 {
		str += ", " + arm.EmitImmediate(a.Offset)
	}
	if a.Multiplier != 0 {
		str += ", " + arm.EmitImmediate(a.Multiplier)
	}
	return str + "]"
}

func (arm Emitter) EmitVariable(v ins.Variable) string {
	return "=" + string(v)
}

func (arm Emitter) EmitFunctionPointer(f ins.FunctionPointer) string {
	return "=" + string(f)
}

func (arm Emitter) EmitDirective(d ins.Directive) string {
	return string(d)
}

var arm11Regs = []string{"r0", "r1", "r2", "r3", "r4", "r5", "r6", "r7", "r8", "r9", "r10", "r11", "r12", "sp", "lr", "pc", "cpsr"}

func (arm Emitter) EmitRegister(r ins.Register) string {
	return arm11Regs[r]
}

func (arm Emitter) EmitPool(_ ins.Pool) string {
	return ".ltorg"
}
