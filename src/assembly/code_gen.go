package assembly

import (
	architecture "wacc_32/assembly/architectures"
	"wacc_32/assembly/architectures/arm11"
	"wacc_32/assembly/builtins"
	ins "wacc_32/assembly/instructions"
	"wacc_32/ast"
)

//go:generate ./../visitor_generator/visitor_generator.sh instructions.Instruction instructions.ExprInstr *instructions.Context

var _ ast.InstructionExprInstrVisitor = &CodeGenerator{}

//CodeGenerator converts an AST into our internal assembly representation
type CodeGenerator struct {
	architecture.Config
	architecture.Emitter
	bssVars   map[string]ins.Instruction
	funcs     map[string]ins.Instruction
	regs      *ins.RegisterManager
	nIfs      int
	numWhiles int
	numDoWhiles int
	numFors   int
}

//newCodeGenerator creates a code generator with a specific number of registers
func newCodeGenerator(conf architecture.Config, emitter architecture.Emitter) *CodeGenerator {
	return &CodeGenerator{
		Config:  conf,
		Emitter: emitter,
		bssVars: make(map[string]ins.Instruction),
		funcs:   make(map[string]ins.Instruction),
		regs:    ins.NewRegisterManager(conf.CalleeSavedRegs, conf.CallerSavedRegs, conf.NRegs),
	}
}

//NewArm11CodeGenerator creates a CodeGenerator which emits arm11 code
func NewArm11CodeGenerator() *CodeGenerator {
	return newCodeGenerator(arm11.Config(), arm11.Emitter{})
}

func (cg *CodeGenerator) GenerateCode(tree ast.AST) string {
	builtins.Init(cg.Config)
	bss, instrs := cg.generateInternalCode(tree)
	return cg.Emit(bss, instrs)
}

//generateInternalCode visits the tree and returns an internal representation of the assembly code
func (cg *CodeGenerator) generateInternalCode(tree ast.AST) (bssVars, instrs ins.Instructions) {
	ctx := &ins.Context{}
	mainInstrs := ins.Instructions{cg.VisitAST(tree, ctx)}
	instrs = make([]ins.Instruction, len(cg.funcs)+len(mainInstrs))
	bssVars = make([]ins.Instruction, len(cg.bssVars))
	i := 0
	for _, v := range cg.bssVars {
		bssVars[i] = v
		i++
	}
	i = 0
	for _, v := range cg.funcs {
		instrs[i] = v
		i++
	}
	for _, v := range mainInstrs {
		instrs[i] = v
		i++
	}
	return bssVars, instrs
}

//VisitProgram visits AST node ast.Program
func (cg *CodeGenerator) VisitProgram(node ast.Program, ctx *ins.Context) ins.Instruction {
	prog := ins.Instructions(make([]ins.Instruction, 0))

	for _, st := range node.GetStructs() {
		prog = append(prog, cg.VisitUserType(*st, ctx))
	}

	for _, fn := range node.GetFuncs() {
		prog = append(prog, cg.VisitFunction(*fn, ctx))
	}

	return prog
}
