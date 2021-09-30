package ast

import (
	"fmt"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//AST represents an Abstract Syntax Tree
type AST interface {
	fmt.Stringer
	Check(ctx Context)
	GetSymbolTable() *symboltable.SymbolTable
}

type ast struct {
	table *symboltable.SymbolTable
}

func (a ast) GetSymbolTable() *symboltable.SymbolTable {
	return a.table
}

//Context is a struct to pass information between nodes during semantic analysis
type Context struct {
	SemanticErrChan chan<- error
	returnType      types.WaccType
	table           *symboltable.SymbolTable
	functionName    string
}
