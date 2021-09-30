package ast

import (
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//UnopType enum
type UnopType int

//All unary operators
const (
	Not UnopType = iota + 1
	Len
	Ord
	Chr
	Neg
	TryLock
)

var _ Expression = &UnOp{}

var unopStrings = []string{"!", "len", "ord", "chr", "-", "try_lock"}

//UnOp represents unary operators
type UnOp struct {
	ast
	op   UnopType
	expr Expression
	pos  errors.Position
}

//GetOpType returns the type of the operation
func (u UnOp) GetOpType() UnopType {
	return u.op
}

//GetExpr returns the expression
func (u UnOp) GetExpr() Expression {
	return u.expr
}

//NewUnOp creates a new UnOp
func NewUnOp(op UnopType, expr Expression, pos errors.Position) *UnOp {
	return &UnOp{
		op:   op,
		expr: expr,
		pos:  pos,
	}
}

//String returns
// (! | len | ord | chr | - | try_lock)
func (op UnopType) String() string {
	return unopStrings[int(op)-1]
}

//EvalType returns
// a bool for try_lock
// a boolean for !
// a character for Chr
// an integer otherwise
func (u *UnOp) EvalType(s symboltable.SymbolTable) types.WaccType {
	switch u.op {
	case TryLock:
		fallthrough
	case Not:
		return types.Boolean
	case Chr:
		return types.Char
	default:
		return types.Integer
	}
}

//String returns
// op
//   - expr
func (u UnOp) String() string {
	return format(u.op.String(), u.expr.String())
}

//Check ensures that the types match and that the subexpression is correct
func (u *UnOp) Check(ctx Context) bool {
	u.table = ctx.table
	if !u.expr.Check(ctx) {
		return false
	}
	exprT := u.expr.EvalType(*ctx.table)

	integer := types.Integer
	boolean := types.Boolean
	character := types.Char
	lock := types.Lock
	str := types.Str
	op := u.op

	var expectedType types.WaccType
	if op == Not && exprT != boolean {
		expectedType = types.Boolean
	} else if op == Len && exprT != str && !exprT.Is(types.Array) {
		expectedType = types.Str
	} else if op == Neg && exprT != integer {
		expectedType = types.Integer
	} else if op == Ord && exprT != character {
		expectedType = types.Char
	} else if op == Chr && exprT != integer {
		expectedType = types.Integer
	} else if op == TryLock && exprT != lock {
		expectedType = types.Boolean
	} else {
		return true
	}
	ctx.SemanticErrChan <- errors.NewTypeError(u.pos, op.String(), expectedType, exprT)
	return false
}
