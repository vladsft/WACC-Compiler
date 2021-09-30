package ast

import (
	"fmt"
	"strings"
	"sync"
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//RHS represents expressions that can only occur on the right hand side
type RHS Expression

var (
	_ RHS = &RHSNewPair{}
	_ RHS = &RHSFunctionCall{}
	_ RHS = &PairElem{}
	_ RHS = &Make{}
)

//RHSNewPair represents a new symboltable.Pair
type RHSNewPair struct {
	ast
	fst, snd Expression
}

//NewRHSNewPair creates a new RHSNewPair
func NewRHSNewPair(fst, snd Expression) *RHSNewPair {
	return &RHSNewPair{
		fst: fst,
		snd: snd,
	}
}

//GetExpr returns either the first or snd expression of the rhs new pair
func (np RHSNewPair) GetExpr(i int) Expression {
	if i == 0 {
		return np.fst
	}
	return np.snd
}

//String returns:
// NewPair
//   - fst.String()
//   - snd.String()
func (np RHSNewPair) String() string {
	fstStr := format("FST", np.fst.String())
	sndStr := format("SND", np.snd.String())
	return format("NEW_PAIR", fstStr, sndStr)
}

//Check semantic validity of new pairs
//A pair is valid iff the fst and snd are valid
func (np *RHSNewPair) Check(ctx Context) bool {
	np.table = ctx.table
	return concurrentCheck(ctx, np.fst, np.snd)
}

func concurrentCheck(ctx Context, exprs ...Expression) bool {
	okChan := make(chan bool, len(exprs))
	var wg sync.WaitGroup //WaitGroup for synchronisation like a semaphore

	for _, expr := range exprs {
		wg.Add(1)
		go func(expr Expression) {
			defer wg.Done()
			okChan <- expr.Check(ctx)
		}(expr)
	}

	//The WaitGroup Wait ensures that all the checks have been done
	wg.Wait()
	close(okChan)
	for b := range okChan {
		if !b {
			return false
		}
	}
	return true
}

// EvalType returns a pair with the fst and snd types
func (np RHSNewPair) EvalType(s symboltable.SymbolTable) types.WaccType {
	return types.NewPair(np.fst.EvalType(s), np.snd.EvalType(s))
}

//RHSFunctionCall represents a function call
type RHSFunctionCall struct {
	ast
	fName      *Ident
	args       []Expression
	pos        errors.Position
	concurrent bool
	isMethod   bool
}

//GetName returns the name of the function being called
func (fnc RHSFunctionCall) GetName() string {
	return fnc.fName.String()
}

//GetName returns the name of the function being called
func (fnc RHSFunctionCall) GetInternalName() string {
	return fnc.fName.GetName()
}

//GetArgs returns the arguments of the function being called
func (fnc RHSFunctionCall) GetArgs() []Expression {
	return fnc.args
}

//NewRHSFunctionCall creates a new FunctionCall
func NewRHSFunctionCall(fName *Ident, args []Expression, isMethod bool, pos errors.Position) *RHSFunctionCall {
	return &RHSFunctionCall{
		fName:      fName,
		args:       args,
		pos:        pos,
		concurrent: false,
		isMethod:   isMethod,
	}
}

//String returns
// CALL
//   - <fName>
//   - arg1
//   ...
//   - argn
func (fnc RHSFunctionCall) String() string {
	str := fnc.fName.String() + "\n"
	args := argListString(fnc.args)
	if args != "" {
		str += header + args
	}
	return str
}

//String returns
// - arg1
// ...
// - argn
func argListString(aList []Expression) string {
	if len(aList) == 0 {
		return ""
	}
	str := aList[0].String() + footer
	for _, ex := range aList[1:] {
		str += header + ex.String() + footer
	}
	return strings.TrimSpace(str)
}

func (fnc *RHSFunctionCall) FormatName() (string, error) {
	sym := *fnc.table
	toLookUp := fnc.fName.GetName()

	if fnc.isMethod {
		components := fnc.fName.GetNameComponents()
		classType, err := sym.GetType(components[0][1:])
		if err != nil {
			return "", errors.NewUndefinedIdentifierError(fnc.pos, err)
		}
		t, err := LookupUserType(classType.(types.UserType), *fnc.table)
		if err != nil {
			return "", errors.NewUndefinedIdentifierError(fnc.pos, err)
		}
		for i := 1; i < len(components)-1; i++ {
			fieldName := components[i]
			for j, field := range t.GetFieldNames() {
				if field == fieldName {
					t = t.GetFieldTypes()[j].(types.UserType)
					break
				}
			}
		}

		className := t.GetName()
		toLookUp = "0" + components[len(components)-1] + "_" + className
	}

	return toLookUp, nil
}

//Check ensures that the function exists, and its arguments are valid
func (fnc *RHSFunctionCall) Check(ctx Context) bool {
	ok := true
	fnc.table = ctx.table

	toLookup, err := fnc.FormatName()
	if err != nil {
		ctx.SemanticErrChan <- err
		return false
	}
	fType, _ := ctx.table.GetType(toLookup)
	paramTypes := fType.GetChildren()
	expArgLength := len(paramTypes) - 1
	actArgLength := len(fnc.args)
	if expArgLength != actArgLength {
		ctx.SemanticErrChan <- errors.NewArgCountError(fnc.pos, fnc.fName.String(), expArgLength, actArgLength)
		ok = false
	}

	for i, arg := range fnc.args {
		if arg.Check(ctx) {
			argType := arg.EvalType(*ctx.table)
			if !paramTypes[i].Is(argType) {
				nodeName := fmt.Sprintf("function %s", fnc.fName.name)
				ctx.SemanticErrChan <- errors.NewTypeError(fnc.pos, nodeName, paramTypes[i], argType)
				ok = false
			}
		}
	}

	return ok
}

//EvalType returns the type of a function call
//Assumes the function call has already passed semantic checks
func (fnc RHSFunctionCall) EvalType(s symboltable.SymbolTable) types.WaccType {
	toLookup, _ := fnc.FormatName()
	returnType, _ := s.GetType(toLookup)
	types := returnType.GetChildren()
	return types[len(types)-1]
}

//PairElemPos enum
type PairElemPos int

//Enum representing the 2 elements of a pair
const (
	FST PairElemPos = iota + 1
	SND
)

//PairElem represents either the first or the second symboltable.Pair element
type PairElem struct {
	ast
	pElemPos PairElemPos
	value    Expression
}

//GetPairElemPos returns either FST or SND
func (p PairElem) GetPairElemPos() PairElemPos {
	return p.pElemPos
}

//GetValue returns a pair
func (p PairElem) GetValue() Expression {
	return p.value
}

//NewPairElem creates a new PairElem
func NewPairElem(pElemType PairElemPos, value Expression) *PairElem {
	return &PairElem{
		pElemPos: pElemType,
		value:    value,
	}
}

// Check semantic validity of FunctionCall
// No need to check pElemPos, its type gives already either FST or SND
func (p *PairElem) Check(ctx Context) bool {
	p.table = ctx.table
	return p.value.Check(ctx)
}

// EvalType returns the type of the value
func (p PairElem) EvalType(s symboltable.SymbolTable) types.WaccType {
	children := p.value.EvalType(s).GetChildren()
	if len(children) != 2 {
		return types.None
	}
	return children[int(p.pElemPos)-1]
}

//String returns
// (FST|SND)
//   - value
func (p PairElem) String() string {
	str := "FST"
	if p.pElemPos == SND {
		str = "SND"
	}
	return format(str, p.value.String())
}

type Make struct {
	ast
	t      types.WaccType
	length Expression
}

func NewMake(t types.WaccType, length Expression) *Make {
	return &Make{
		t:      t,
		length: length,
	}
}

func (m Make) GetLengthExpression() Expression {
	return m.length
}

//Check semantic validity of Make
//Checks if length is an integer
func (m *Make) Check(ctx Context) bool {
	m.table = ctx.table
	if !m.length.Check(ctx) {
		return false
	}
	return m.length.EvalType(*ctx.table).Is(types.Integer)
}

// EvalType returns an array of m.t
func (m Make) EvalType(_ symboltable.SymbolTable) types.WaccType {
	return types.NewArray(m.t, 1)
}

//String returns
// MAKE
//   - type
//   - length
func (m Make) String() string {
	return format("MAKE", m.t.String(), m.length.String())
}
