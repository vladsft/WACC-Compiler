package ast

import (
	"fmt"
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

var (
	_ AST = &Function{}
	_ AST = &ParamList{}
	_ AST = &Param{}
)

// Function AST
type Function struct {
	ast
	retType      types.WaccType
	ident        *Ident
	params       ParamList
	stats        StatMultiple
	pos          errors.Position
	isConcurrent bool
	isMethod     bool
}

//NewFunction returns a function
func NewFunction(retType types.WaccType, ident *Ident, params ParamList,
	stats Statement, pos errors.Position) *Function {
	fn := &Function{
		retType: retType,
		ident:   ident,
		params:  params,
		pos:     pos,
	}
	switch st := stats.(type) {
	case StatMultiple:
		fn.stats = st
	default:
		fn.stats = StatMultiple{st}
	}
	return fn
}

// NewMainFunction returns the function "int main()"
func NewMainFunction(stat Statement, pos errors.Position) *Function {
	return NewFunction(types.Integer, NewIdent("0main", pos), NewEmptyParamList(), stat, pos)
}

//SetConcurrent marks the function as concurrent
func (f *Function) SetConcurrent() {
	f.isConcurrent = true
}

func (f *Function) IsConcurrent() bool {
	return f.isConcurrent
}

//GetName returns the function's name
func (f Function) GetName() string {
	return f.ident.String()
}

//GetParams returns the parameters the function takes
func (f Function) GetParams() ParamList {
	return f.params
}

//GetStats returns the function's statements
func (f Function) GetStats() []Statement {
	return f.stats
}

//String returns
// <return_type> <function_name>(<params>)
//   - stat0
//   ...
//   - statn
func (f Function) String() string {
	s := make([]string, len(f.stats))
	for i, st := range f.stats {
		s[i] = st.String()
	}
	return format(fmt.Sprintf("%s %s(%s)", f.retType, f.ident.name, f.params.String()), s...)
}

func (f *Function) MakeMethod(ut types.UserType) {
	this := NewIdent("this", f.pos)
	f.ident = NewIdent(f.ident.GetName()+"_"+ut.GetName(), f.pos)
	f.params = append(f.params, NewParam(ut, this, f.pos))
	f.isMethod = true
}

func (f Function) IsMethod() bool {
	return f.isMethod
}

func (f *Function) Check(ctx Context) {
	var fCtx Context
	if f.GetName() != "main" {
		fCtx = Context{
			functionName:    f.ident.name,
			returnType:      f.retType,
			table:           symboltable.NewSymbolTable(ctx.table),
			SemanticErrChan: ctx.SemanticErrChan,
		}
	} else {
		fCtx = Context{
			table:           symboltable.NewSymbolTable(ctx.table),
			SemanticErrChan: ctx.SemanticErrChan,
		}
	}

	f.params.Check(fCtx)

	//Check all statements with new context
	for _, stat := range f.stats {
		stat.Check(fCtx)
	}
	f.table = fCtx.table
}

// ParamList is a list of parameters
type ParamList []*Param

// NewParamList constructor
func NewParamList(params []*Param) ParamList {
	return ParamList(params)
}

//NewEmptyParamList creates an empty parameter list
func NewEmptyParamList() ParamList {
	return ParamList{}
}

//String returns
// param1,...,paramn
func (pList ParamList) String() string {
	var str string
	for _, p := range pList {
		str += ", " + p.String()
	}
	if len(str) == 0 {
		return ""
	}
	return str[2:]
}

//Check that no duplicate parameters exist
func (pList ParamList) Check(ctx Context) {
	for _, param := range pList {
		param.Check(ctx)
	}
}

//GetSymbolTable returns nil
func (pList *ParamList) GetSymbolTable() *symboltable.SymbolTable {
	return nil
}

// Param is an ident with a type
type Param struct {
	ast
	t     types.WaccType
	ident *Ident
	pos   errors.Position
}

// NewParam constructor
func NewParam(t types.WaccType, ident *Ident, pos errors.Position) *Param {
	return &Param{
		t:     t,
		ident: ident,
		pos:   pos,
	}
}

func (param Param) GetName() string {
	return param.ident.name
}

func (param Param) GetType() types.WaccType {
	return param.t
}

//String returns
// <type> <ident>
func (param Param) String() string {
	return param.t.String() + " " + param.ident.name
}

//Check for param
func (param *Param) Check(ctx Context) {
	param.table = ctx.table
	name := param.ident.name

	if param.t.Is(types.UserDefinedType) {
		paramUT := param.t.(types.UserType)
		_, err := LookupUserType(paramUT, *ctx.table)

		if err != nil {
			ctx.SemanticErrChan <- errors.NewUndefinedIdentifierError(param.pos, err)
		}
	}
	if err := ctx.table.AddDeclaration(param.ident.name, param.t, param.pos); err != nil {
		ctx.SemanticErrChan <- errors.NewParamError(param.pos, name)
	}
}
