package errors

import (
	"fmt"
	"wacc_32/types"
)

type semanticError int

const (
	typeError = iota + 1
	paramError
	returnError
	undefinedIdentifierError
	identifierAlreadyInUseError
	argCountError
	importError
	uninitialisedUserTypeError
	invalidFieldAccessError
)

var semanticErrors = []string{"TypeError", "ParamError", "ReturnError", "UndefinedIdentifierError", "IdentifierAlreadyInUseError", "ArgCountError", "ImportError", "UninitialisedUserTypeError", "InvalidFieldAccess"}

func (s semanticError) String() string {
	return red(semanticErrors[s-1])
}

func newError(p Position, errType semanticError, template string, args ...interface{}) error {
	msg := fmt.Sprintf(template, args...)
	err := fmt.Errorf("%s %s: %s\t👎", p.String(), errType, msg)
	return err
}

func red(s string) string {
	return "\u001b[1m\u001b[31m" + s + "\u001b[0m"
}

/* **************************** FACTORY METHODS **************************** */

//NewTypeError returns
// Line [s:e-s:e] TypeError: <nodeName> expected <expected> not <not>
func NewTypeError(p Position, nodeName string, expected, not types.WaccType) error {
	return newError(p, typeError, "%s expected %s not %s", nodeName, expected, not)
}

//NewTernaryError
func NewTernaryError(p Position, ifType, elseType types.WaccType) error {
	return newError(p, typeError, "ternary operator expected types of both branches to be the same. '?' branch had type %s and  ':' branch had type %s ", ifType, elseType)
}

// Line [s:e-s:e] TypeError: <nodeName> expected one of [<expected>+] not <not>
func NewMultiTypeError(p Position, nodeName string,
	not types.WaccType, expected ...types.WaccType) error {
	var str string
	last := len(expected) - 1
	for _, exp := range expected[:last] {
		str += exp.String() + ", "
	}
	str += expected[last].String()
	return newError(p, typeError, "%s expected one of [%s] not %s", nodeName, str, not)
}

//NewParamError returns
// Line [s:e-s:e] ParamError: param <pName> duplicated
func NewParamError(p Position, pName string) error {
	return newError(p, paramError, "parameter %s duplicated", pName)
}

//NewReturnError returns
// Line [s:e-s:e] ReturnError: cannot return outside of function
func NewReturnError(p Position) error {
	return newError(p, returnError, "cannot return outside of function")
}

//NewUndefinedIdentifierError returns
// Line [s:e-s:e] UndefinedIdentifierError: <ident> has not been defined in the current context
func NewUndefinedIdentifierError(p Position, err error) error {
	return newError(p, undefinedIdentifierError, err.Error())
}

//NewIdentifierAlreadyInUseError returns
// Line [s:e-s:e] IdentifierAlreadyInUseError: <name> already declared at <original_position>
func NewIdentifierAlreadyInUseError(p Position, name string, firstPos Position) error {
	return newError(p, identifierAlreadyInUseError, "%s already declared at %s", name, firstPos)
}

//NewArgCountError returns
// Line [s:e-s:e] ArgCountError: wrong number of arguments for function <function>, expected <n> got
// <m>
func NewArgCountError(p Position, fname string, n, m int) error {
	return newError(
		p,
		argCountError,
		"wrong number of arguments for function %s, expected %d got %d",
		fname,
		n,
		m,
	)
}

//NewSameTypeError returns
// Line [s:e-s:e] TypeError: <op> requires both arguments to have the same type
func NewSameTypeError(p Position, op string) error {
	return newError(p, typeError, "%s requires both arguments to have the same type", op)
}

//NewArrayTypeError returns
// Line [s:e-s:e] TypeError: all array elements must have the same type
func NewArrayTypeError(p Position) error {
	return newError(p, typeError, "all array elements must have the same type")
}

//NewLibraryNotFoundError returns
// Line [s:e-s:e] ImportError: library <lib> not found
func NewLibraryNotFoundError(lib string, p Position) error {
	return newError(p, importError, "library %s not found", lib)
}

//NewFunctionNotDefinedError
// Line [s:e-s:e] ImportError: function <func> not defined
func NewFunctionNotDefinedError(function string, p Position) error {
	return newError(p, importError, "function %s not defined", function)
}

//NewLibraryNotImportedError
// Line [s:e-s:e] ImportError: library <lib> not imported
func NewLibraryNotImportedError(lib string, p Position) error {
	return newError(p, importError, "library %s not imported", lib)
}

//NewInvalidImportPathError
func NewInvalidImportPathError(p Position) error {
	return newError(p, importError, "import paths cannot contain '$' symbols")
}

//NewUninitialisedStructError returns
// Line [s:e-s:e] uninitialised objects not allowed
func NewUninitialisedUsertypeError(p Position) error {
	return newError(p, uninitialisedUserTypeError, "uninitialised objects not allowed")
}

//NewInvalidFieldAccessError returns
// Line [s:e-s:e] field <fieldName> is not present in struct <structName>
func NewInvalidFieldAccessError(p Position, fieldName, structName string) error {
	return newError(p, invalidFieldAccessError, "field %s is not present in struct %s", fieldName, structName)
}
