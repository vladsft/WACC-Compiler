package ast

import (
	"wacc_32/symboltable"
	"wacc_32/types"
)

var _ AST = &UserType{}

// UserType AST
type UserType struct {
	ast
	ident     *Ident
	fields    []*StatNewassign
	IsClass   bool
	functions []*Function
}

func NewUserType(ident *Ident, fields []*StatNewassign, isClass bool, functions []*Function) *UserType {
	return &UserType{
		ident:     ident,
		fields:    fields,
		IsClass:   isClass,
		functions: functions,
	}
}

//GetName returns the name of the userType
func (ut UserType) GetName() string {
	return ut.ident.GetName()
}

//FieldNameList returns a ordered list of the field names to their respective WaccTypes
func (ut UserType) FieldNameList() []string {
	namesList := make([]string, len(ut.fields))
	for i, field := range ut.fields {
		namesList[i] = field.ident.name
	}
	return namesList
}

//FieldTypeList returns a ordered list of the field names to their respective WaccTypes
func (ut UserType) FieldTypeList() []types.WaccType {
	fieldsList := make([]types.WaccType, len(ut.fields))

	for i, field := range ut.fields {
		fieldsList[i] = field.t
	}
	return fieldsList
}

//FieldTypeList returns a ordered list of the field names to their respective WaccTypes
func (ut UserType) FieldsMap() map[string]types.WaccType {
	fieldsMap := make(map[string]types.WaccType)

	for _, field := range ut.fields {
		fieldsMap[field.GetName()] = field.t
	}
	return fieldsMap
}

func (ut UserType) String() string {
	//String representations of fields
	fieldStrs := make([]string, len(ut.fields))
	for i, field := range ut.fields {
		fieldStrs[i] = field.String()
	}

	//String representations of functions
	funcStrs := make([]string, len(ut.functions))
	for i, fn := range ut.functions {
		funcStrs[i] = fn.String()
	}

	//String representation of class/Struct
	var userType string
	if ut.IsClass {
		userType = "CLASS "
	} else {
		userType = "STRUCT "
	}

	children := append(fieldStrs, funcStrs...)

	return format(userType+ut.ident.name, children...)
}

func (ut *UserType) Check(ctx Context) {
	utCtx := Context{
		table:           symboltable.NewSymbolTable(ctx.table),
		SemanticErrChan: ctx.SemanticErrChan,
	}

	for _, field := range ut.fields {
		field.Check(utCtx)
	}

	for _, fn := range ut.functions {
		fn.Check(utCtx)
	}
	ut.table = utCtx.table
}

func (ut UserType) EvalType() types.WaccType {
	return types.NewUserType(ut.ident.name, ut.FieldNameList(), ut.FieldTypeList(), ut.IsClass)
}
