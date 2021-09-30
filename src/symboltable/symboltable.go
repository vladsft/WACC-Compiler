package symboltable

import (
	"fmt"
	"wacc_32/errors"
	"wacc_32/types"
)

type stMetadata struct {
	wt     types.WaccType
	pos    errors.Position
	offset int
}

//newStMetadata creates a new stNodeMetadata
func newStMetadata(wt types.WaccType, pos errors.Position, offset int) *stMetadata {
	return &stMetadata{
		wt:     wt,
		pos:    pos,
		offset: offset,
	}
}

//SymbolTable manages variable and function declarations along with scoping information
type SymbolTable struct {
	parentScope    *SymbolTable
	declarations   map[string]*stMetadata
	objDeclaration map[string]*SymbolTable
	lastOffset     int
}

//NewTopSymbolTable returns a top level SymbolTable - one without any parents - like batman :(
func NewTopSymbolTable() *SymbolTable {
	return NewSymbolTable(nil)
}

//NewSymbolTable creates a new SymbolTable and returns a pointer to it
func NewSymbolTable(parent *SymbolTable) *SymbolTable {
	offset := 0
	if parent != nil {
		offset = parent.lastOffset
	}
	return &SymbolTable{
		parentScope:  parent,
		declarations: make(map[string]*stMetadata),
		lastOffset:   offset,
	}
}

//AddDeclaration adds a declaration to to the symbol table
func (st *SymbolTable) AddDeclaration(ident string, wt types.WaccType, pos errors.Position) error {
	err := 	st.AddDefinition(ident, wt, pos)
	if err != nil {
		return err
	}
	st.lastOffset += int(types.TypeSize(wt))
	return nil
}
func (st *SymbolTable) AddDefinition(ident string, wt types.WaccType, pos errors.Position) error {
	if val, ok := st.declarations[ident]; ok {
		return errors.NewIdentifierAlreadyInUseError(pos, ident, val.pos)
	}
	st.declarations[ident] = newStMetadata(wt, pos, st.lastOffset)

	return nil
}

func (st *SymbolTable) getMetadata(ident string) (*stMetadata, error) {
	scope, err := st.Find(ident)
	if err != nil {
		return nil, err
	}
	return scope.declarations[ident], nil
}

//Find returns the symbol that contains the ident declaration
func (st *SymbolTable) Find(ident string) (*SymbolTable, error) {
	_, ok := st.declarations[ident]
	if !ok {
		if st.parentScope == nil {
			return nil, fmt.Errorf("%s has not been declared in the current context", ident)
		}
		return st.parentScope.Find(ident)
	}
	return st, nil
}

//GetType returns a pointer to the WaccType of an identifier, or an error if it doesn't exist
func (st *SymbolTable) GetType(ident string) (types.WaccType, error) {
	metadata, err := st.getMetadata(ident)
	if err != nil {
		return nil, err
	}
	return metadata.wt, err
}

//SetOffset sets the offset field for the ident in the current context
//Assumes that the ident exists in the current context
func (st *SymbolTable) SetOffset(ident string, offset int) {
	metadata, err := st.getMetadata(ident)
	if err != nil {
		panic(fmt.Errorf("Ident %s not found", ident))
	}
	metadata.offset = offset
}

//GetOffset sets the offset field for the ident in the current context
//Assumes that the ident exists
func (st *SymbolTable) GetOffset(ident string) int {
	metadata, err := st.getMetadata(ident)
	if err != nil {
		panic("Ident not found")
	}
	return metadata.offset
}

//GetTotalOffset returns the stack offset expected for a function
func (st *SymbolTable) GetTotalOffset() int {
	return st.lastOffset
}

//SetTotalOffset sets the total offset for a symbol table
func (st *SymbolTable) SetTotalOffset(offset int) {
	st.lastOffset = offset
}

//PrintAllIdents prints all the identifiers in the symbol table
//DEBUG ONLY
func (st *SymbolTable) PrintAllIdents(recurse bool) {
	for f := range st.declarations {
		fmt.Println(f)
	}

	if recurse {
		if pst := st.parentScope; pst != nil {
			pst.PrintAllIdents(true)
		}
	}
}
