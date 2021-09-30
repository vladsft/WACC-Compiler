package builtins

import (
	ins "wacc_32/assembly/instructions"
	"wacc_32/types"
)

type RuntimeErrType int

//go:generate stringer -type=RuntimeErrType
const (
	ArrayIndexTooLargeError RuntimeErrType = iota + 1
	ArrayIndexNegativeError
	IntegerOverflowError
	DivideByZeroError
	NullPointerReferenceError
	InvalidThreadUnlockError
	SameThreadLockError
)

// Aliases for each of the runtime errors output string
const (
	arrayIndexTooLargeMsg   = `"ArrayIndexOutOfBoundsError: index too large\0"`
	arrayIndexNegativeMsg   = `"ArrayIndexOutOfBoundsError: negative index\0"`
	integerOverFlowMsg      = `"OverflowError: the result is too small/large to store in a 4-byte signed-integer.\0"`
	divideByZeroMsg         = `"DivideByZeroError: divide or modulo by zero\0"`
	nullPointerReferenceMsg = `"NullReferenceError: dereference a null reference\0"`
	invalidThreadUnlockMsg  = `"InvalidThreadUnlockError: can't release lock\0"`
	sameThreadLockMsg       = `"Deadlock: attempted to acquire an acquired lock in the same thread\0"`
)

var errorMsgs = []string{
	arrayIndexTooLargeMsg,
	arrayIndexNegativeMsg,
	integerOverFlowMsg,
	divideByZeroMsg,
	nullPointerReferenceMsg,
	invalidThreadUnlockMsg,
	sameThreadLockMsg,
}

// Aliases for the labels of the runtime error code
const (
	PrintErrorCheckLabel           = "p_print_error"
	DivideByZeroCheckLabel         = "p_check_divide_by_zero"
	IntegerOverflowCheckLabel      = "p_check_int_overflow"
	ArrayOutOfBoundsCheckLabel     = "p_check_array_index_out_of_bounds"
	NullPointerReferenceCheckLabel = "p_check_null_pointer"
	InvalidThreadUnlockCheckLabel  = "p_check_invalid_thread_unlock"
	SameThreadLockCheckLabel       = "p_check_same_thread_lock"
)

func (err RuntimeErrType) GetMsg() string {
	return errorMsgs[err-1]
}

func (err RuntimeErrType) GetErrMsgLabel() ins.Variable {
	return ins.Variable("err_" + err.String())
}

// Assembly code for printing runtime errors called by the variables down below
func PrintError() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(PrintErrorCheckLabel),
		ins.NewPush(lr),
		ins.NewLoad(ins.NewAddress(returnReg), arg1, types.Word),
		ins.NewAdd(returnReg, returnReg, ins.Immediate(types.Word)),
		ins.NewFunctionCall("printf"),
		ins.NewLoad(ins.NewDirective("streams"), returnReg, types.Word),
		ins.NewLoad(ins.NewAddress(returnReg), returnReg, types.Word),
		ins.NewFunctionCall("fflush"),
		ins.NewExit(ins.Immediate(-1)),
		ins.NewPop(pc),
	}
}

// Assembly code for the function that checks whether division by 0 occurs
// If this happens, the program throws a "divide by 0" error
func CheckDivideByZeroError() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(DivideByZeroCheckLabel),
		ins.NewPush(lr),
		ins.NewCompare(arg1, ins.Immediate(0)),
		ins.NewBranch("ok_divide", ins.NE),
		ins.NewLoad(DivideByZeroError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
		ins.NewLabel("ok_divide"),
		ins.NewPop(pc),
	}
}

// Assembly code for the function that checks whether integer overflow/underflow occurs
// If this happens, the program throws an overflow error
func CheckIntegerOverflowError() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(IntegerOverflowCheckLabel),
		ins.NewPush(lr),
		ins.NewBranch("not_ok", ins.VS),
		ins.NewPop(pc),
		ins.NewLabel("not_ok"),
		ins.NewLoad(IntegerOverflowError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
	}
}

// Assembly code for the function that checks whether a null variable is used/called
// If this happens, the program throws a Null Pointer Exception
func CheckNullPointerReferenceError() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(NullPointerReferenceCheckLabel),
		ins.NewPush(lr),
		ins.NewCompare(returnReg, ins.Immediate(0)),
		ins.NewBranch("ok_null", ins.NE),
		ins.NewLoad(NullPointerReferenceError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
		ins.NewLabel("ok_null"),
		ins.NewPop(pc),
	}
}

// Assembly code for the function that checks whether the code uses an index to access an array element out of the array's bounds
// If this happens, the program throws an  "Array Out of Bounds error"
func CheckArrayIndexOutOfBounds() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(ArrayOutOfBoundsCheckLabel),
		ins.NewPush(lr),
		ins.NewCompare(returnReg, ins.Immediate(0)),
		ins.NewBranch("negative_index_error", ins.LT),
		ins.NewCompare(returnReg, arg1),
		ins.NewBranch("index_too_large_error", ins.GE),
		ins.NewPop(pc),
		ins.NewLabel("negative_index_error"),
		ins.NewLoad(ArrayIndexNegativeError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
		ins.NewLabel("index_too_large_error"),
		ins.NewLoad(ArrayIndexTooLargeError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
		ins.NewPop(pc),
	}
}

//CheckInvalidThreadUnlock checks if a release returns EPERM
func CheckInvalidThreadUnlock() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(InvalidThreadUnlockCheckLabel),
		ins.NewPush(lr),
		ins.NewCompare(returnReg, ins.Immediate(1)), //EPERM
		ins.NewBranch(InvalidThreadUnlockCheckLabel+"_done", ins.NE),
		ins.NewLoad(InvalidThreadUnlockError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
		ins.NewLabel(InvalidThreadUnlockCheckLabel + "_done"),
		ins.NewPop(pc),
	}
}

//CheckSameThreadLock checks if a thread attempts to acquire a thread twice
func CheckSameThreadLock() ins.Instructions {
	return ins.Instructions{
		ins.NewLabel(SameThreadLockCheckLabel),
		ins.NewPush(lr),
		ins.NewCompare(returnReg, ins.Immediate(35)), //EDEADLK
		ins.NewBranch(SameThreadLockCheckLabel+"_done", ins.NE),
		ins.NewLoad(SameThreadLockError.GetErrMsgLabel(), returnReg, types.Word),
		ins.NewFunctionCall("p_print_error"),
		ins.NewLabel(SameThreadLockCheckLabel + "_done"),
		ins.NewPop(pc),
	}
}
