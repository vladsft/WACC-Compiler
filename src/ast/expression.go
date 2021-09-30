package ast

import (
	"fmt"
	"strings"
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

// Expression AST interface
type Expression interface {
	fmt.Stringer
	Check(ctx Context) (ok bool)
	EvalType(symboltable.SymbolTable) types.WaccType
	GetSymbolTable() *symboltable.SymbolTable
}

var _ Expression = &ArrayElem{}
var _ Expression = &Ident{}

//ArrayElem represents an Array Element
type ArrayElem struct {
	ast
	ident   *Ident
	indices []Expression
	pos     errors.Position
}

func NewArrayElem(ident *Ident, indices []Expression, pos errors.Position) *ArrayElem {
	return &ArrayElem{
		ident:   ident,
		indices: indices,
		pos:     pos,
	}
}

//GetIdent returns the identifier of the array element
func (a ArrayElem) GetIdent() *Ident {
	return a.ident
}

//GetIndices returns the indices of the array element
func (a ArrayElem) GetIndices() []Expression {
	return a.indices
}

//GetPos returns the position of the array elements in its array
func (a ArrayElem) GetPos() errors.Position {
	return a.pos
}

//Check ensures every index is a valid integer
func (a *ArrayElem) Check(ctx Context) bool {
	a.table = ctx.table
	a.ident.Check(ctx)
	ok := true
	t, err := ctx.table.GetType(a.ident.name)
	if err != nil {
		ctx.SemanticErrChan <- errors.NewUndefinedIdentifierError(a.pos, err)
		ok = false
	} else if !t.Is(types.Array) {
		ctx.SemanticErrChan <- errors.NewTypeError(a.pos, "\b", types.Array, t)
		ok = false
	}
	for _, expr := range a.indices {
		if !expr.Check(ctx) {
			ok = false
		} else if indexType := expr.EvalType(*ctx.table); indexType != types.Integer {
			ctx.SemanticErrChan <- errors.NewTypeError(a.pos, "array index", types.Integer, indexType)
		}
	}
	return ok
}

//String returns
// <type>[]
//   - idx1
//   ...
//   - idxn
func (a ArrayElem) String() string {
	s := make([]string, len(a.indices))
	for i, idx := range a.indices {
		s[i] = format("[]", idx.String())
	}
	return format(a.ident.String(), s...)
}

//EvalType returns the type of an array elem
//Assuming that the ArrayElem is semantically valid
func (a ArrayElem) EvalType(s symboltable.SymbolTable) types.WaccType {
	wt, _ := s.GetType(a.ident.name)
	subType := wt.GetChildren()[0]
	for i := 0; i < len(a.indices)-1; i++ {
		subType = subType.GetChildren()[0]
	}
	return subType
}

//Ident represents an identifier
type Ident struct {
	ast
	name       string
	namespaced bool
	imported   bool
	pos        errors.Position
}

//NewIdent creates a new Ident with default parameters
func NewIdent(name string, pos errors.Position) *Ident {
	namespaced := strings.Contains(name, "!")
	imported := strings.Contains(name, "::")
	return &Ident{
		name:       name,
		namespaced: namespaced,
		imported:   imported,
		pos:        pos,
	}
}

func (i Ident) IsNamespaced() bool {
	return i.namespaced
}

//String returns the ident name
//functions and structs/classes = 0name
//imported stuff				= dir$dir$dir$dir$file$0name
//field accesses 				= 1class!name
func (i Ident) String() string {
	if i.namespaced {
		return strings.Replace(i.name[1:], "!", " ", 1)
	}
	if i.imported {
		return strings.ReplaceAll(i.name, "$", "/")[1:]
	}
	if i.name[0] == '0' {
		return i.name[1:]
	}
	return i.name
}

//GetPos returns the position of the identifier
func (i Ident) GetPos() errors.Position {
	return i.pos
}

//GetName returns the name of the identifier
func (i Ident) GetName() string {
	return i.name
}

func (i Ident) GetNameComponents() []string {
	if i.namespaced {
		return strings.Split(i.name[1:], "!")
	}
	if i.imported {
		components := strings.Split(i.name, "$")
		return components[len(components)-1:]
	}
	if i.name[0] == '0' {
		return []string{i.name[1:]}
	}
	return []string{i.name}
}

//Check makes sure that the identifier exists
func (i *Ident) Check(ctx Context) bool {
	var scope *symboltable.SymbolTable
	var err error
	if i.namespaced {
		components := i.GetNameComponents()

		t, err1 := ctx.table.GetType(components[0])
		if err1 != nil {
			err = err1
			goto Error
		}
		for i := 1; i < len(components); i++ {
			fieldName := components[i]
			uType, err1 := LookupUserType(t.(types.UserType), *ctx.table)
			if err1 != nil {
				err = err1
				goto Error
			}
			t, err = uType.GetType(fieldName)
			if err != nil {
				goto Error
			}
		}
		scope = ctx.table
	} else {
		scope, err = ctx.table.Find(i.name)
	}

Error:
	if err != nil {
		ctx.SemanticErrChan <- errors.NewUndefinedIdentifierError(i.pos, err)
	}

	i.table = scope
	return err == nil
}

//EvalType returns the type of the ident in the symbol table
//It is assumed that the ident is a symbol
func (i Ident) EvalType(s symboltable.SymbolTable) types.WaccType {
	if i.namespaced {
		components := i.GetNameComponents()

		t, _ := s.GetType(components[0])
		for j := 1; j < len(components); j++ {
			fieldName := components[j]
			uType, _ := LookupUserType(t.(types.UserType), *i.table)
			t, _ = uType.GetType(fieldName)
		}
		return t
	}
	t, _ := s.GetType(i.name)
	return t
}
