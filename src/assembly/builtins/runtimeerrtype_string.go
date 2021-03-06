// Code generated by "stringer -type=RuntimeErrType"; DO NOT EDIT.

package builtins

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ArrayIndexTooLargeError-1]
	_ = x[ArrayIndexNegativeError-2]
	_ = x[IntegerOverflowError-3]
	_ = x[DivideByZeroError-4]
	_ = x[NullPointerReferenceError-5]
	_ = x[InvalidThreadUnlockError-6]
	_ = x[SameThreadLockError-7]
}

const _RuntimeErrType_name = "ArrayIndexTooLargeErrorArrayIndexNegativeErrorIntegerOverflowErrorDivideByZeroErrorNullPointerReferenceErrorInvalidThreadUnlockErrorSameThreadLockError"

var _RuntimeErrType_index = [...]uint8{0, 23, 46, 66, 83, 108, 132, 151}

func (i RuntimeErrType) String() string {
	i -= 1
	if i < 0 || i >= RuntimeErrType(len(_RuntimeErrType_index)-1) {
		return "RuntimeErrType(" + strconv.FormatInt(int64(i+1), 10) + ")"
	}
	return _RuntimeErrType_name[_RuntimeErrType_index[i]:_RuntimeErrType_index[i+1]]
}
