package ast

import (
	"strings"
	"sync"
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

// Program AST
type Program struct {
	ast
	userTypes []*UserType
	funcs     []*Function
	pos       errors.Position
}

const (
	header = "- "
	footer = "\n  "
)

var _ AST = &Program{}

func format(name string, children ...string) string {
	var str string
	if name != "" {
		str = name + footer
	}
	for _, child := range children {
		if child != "" {
			str += header + strings.Replace(child, "\n", footer, -1) + footer
		}
	}
	return strings.TrimSpace(str)
}

// NewProgram constructor
func NewProgram(structs []*UserType, funcs []*Function, pos errors.Position) *Program {
	return &Program{
		userTypes: structs,
		funcs:     funcs,
		pos:       pos,
	}
}

//GetFuncs returns the program's functions
func (prog Program) GetFuncs() []*Function {
	return prog.funcs
}

//GetStructs returns the program's structures
func (prog Program) GetStructs() []*UserType {
	return prog.userTypes
}

//String returns
// Program
//   - func1
//   - ...
//   - funcn
func (prog Program) String() string {
	funcs := make([]string, len(prog.funcs))
	for i, f := range prog.funcs {
		funcs[i] = f.String()
	}

	structs := make([]string, len(prog.userTypes))
	for i, s := range prog.userTypes {
		structs[i] = s.String()
	}
	children := append(structs, funcs...)
	return format("Program", children...)
}

//Check semantic validity of Program
//The context is ignored as it creates its own
func (prog *Program) Check(ctx Context) {
	st := symboltable.NewTopSymbolTable()
	utCtx := Context{
		functionName:    "",
		table:           st,
		SemanticErrChan: ctx.SemanticErrChan,
	}

	for _, ut := range prog.userTypes {
		//Declare Class/Struct
		utName := ut.GetName()
		err := utCtx.table.AddDefinition("2"+utName, ut.EvalType(), prog.pos)
		if err != nil {
			ctx.SemanticErrChan <- err
		}
	}

	fCtx := Context{
		functionName:    "",
		table:           st,
		SemanticErrChan: ctx.SemanticErrChan,
	}

	prog.table = ctx.table
	for _, fn := range prog.funcs {
		//Declare function
		paramTypes := make([]types.WaccType, len(fn.params))
		for i, param := range fn.params {
			paramTypes[i] = param.t
		}
		fType := types.NewFunction(fn.retType, paramTypes)
		err := fCtx.table.AddDefinition(fn.ident.name, fType, prog.pos)
		if err != nil {
			ctx.SemanticErrChan <- err
		}
	}

	//Concurrent Semantic Analysis of structures
	var wg sync.WaitGroup
	for _, ut := range prog.userTypes {
		wg.Add(1)
		go func(ut *UserType) {
			defer wg.Done()
			ut.Check(utCtx)
		}(ut)
	}
	wg.Wait()

	//Concurrent Semantic Analysis of functions
	for _, fn := range prog.funcs {
		wg.Add(1)
		go func(fn *Function) {
			defer wg.Done()
			fn.Check(fCtx)
			fCtx.table.SetOffset(fn.ident.name, fn.table.GetTotalOffset())
		}(fn)
	}
	wg.Wait()

	close(ctx.SemanticErrChan)
}

func LookupUserType(uType types.UserType, table symboltable.SymbolTable) (types.UserType, error) {
	if uType.GetFieldTypes() != nil {
		return uType, nil
	}
	wt, err := table.GetType("2" + uType.GetName())
	if err != nil {
		return types.UserType{}, err
	}
	return wt.(types.UserType), nil
}
