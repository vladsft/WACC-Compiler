package types

import (
	"fmt"
)

var _ WaccType = UserType{}

type UserType struct {
	name       string
	fieldTypes []WaccType
	fieldNames []string
	IsClass    bool
	funcTypes  []WaccType
}

//NewUserType creates a new userType type
func NewUserType(name string, fieldNames []string, fieldTypes []WaccType, isClass bool) UserType {
	return UserType{
		name:       name,
		fieldTypes: fieldTypes,
		fieldNames: fieldNames,
		IsClass:    isClass,
	}
}

func (s UserType) GetName() string {
	return s.name
}

func (s UserType) GetFieldNames() []string {
	return s.fieldNames
}

func (s UserType) GetType(name string) (WaccType, error) {
	i := findString(name, s.fieldNames)
	if i == -1 {
		return nil, fmt.Errorf("%s does not contain field %s", s.name, name)
	}
	return s.fieldTypes[i], nil
}

func findString(s string, strs []string) int {
	for i, v := range strs {
		if v == s {
			return i
		}
	}
	return -1
}

//DefaultValue should not be called for this
func (s UserType) DefaultValue() interface{} {
	return 0
}

//GetChildren returns a list of the WaccTypes of the fields of the structs
func (s UserType) GetFieldTypes() []WaccType {
	return s.fieldTypes
}

//GetChildren returns a list of the WaccTypes of the fields of the structs
func (s UserType) GetChildren() []WaccType {
	return s.fieldTypes
}

func (s UserType) GetFormatString() string {
	return "%p"
}

func (s UserType) Is(wt WaccType) bool {
	switch w := wt.(type) {
	case waccBaseType:
		return UserDefinedType == w
	case UserType:
		return s.name == w.name
	default:
		return false
	}
}

func (s UserType) String() string {
	if s.IsClass {
		return "struct"
	} else {
		return "class"
	}
}
