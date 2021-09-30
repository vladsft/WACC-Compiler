package ast

import (
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//BinopType enum
type BinopType int

//Enum representing the different binary operators
const (
	Star BinopType = iota + 1
	Div
	Mod
	Plus
	Minus
	Greater
	GreaterEq
	Less
	LessEq
	Equal
	NotEq
	And
	Or
)

var binopStrings = []string{"*", "/", "%", "+", "-", ">", ">=", "<", "<=", "==", "!=", "&&", "||"}

//String returns the encoding of a binop
func (b BinopType) String() string {
	return binopStrings[b-1]
}

var _ Expression = &BinOp{}

//BinOp represents a binary operator
type BinOp struct {
	ast
	op          BinopType
	left, right Expression
	pos         errors.Position
}

//NewBinOp creates a new BinOp
func NewBinOp(op BinopType, l, r Expression, pos errors.Position) *BinOp {
	return &BinOp{
		op:    op,
		left:  l,
		right: r,
		pos:   pos,
	}
}

//GetLeftExpr returns the left expression of the binary operator
func (b BinOp) GetLeftExpr() Expression {
	return b.left
}

//GetRightExpr returns the right expression of the binary operator
func (b BinOp) GetRightExpr() Expression {
	return b.right
}

//GetOpType returns the binary operator
func (b BinOp) GetOpType() BinopType {
	return b.op
}

//Check makes sure both sides have the same and correct type
func (b *BinOp) Check(ctx Context) bool {
	b.table = ctx.table
	if !concurrentCheck(ctx, b.left, b.right) {
		return false
	}
	lt := b.left.EvalType(*ctx.table)
	rt := b.right.EvalType(*ctx.table)

	ts := [2]types.WaccType{lt, rt}

	intType := types.Integer
	charType := types.Char
	boolType := types.Boolean
	pairType := types.Pair

	op := b.op.String()
	ok := true
	switch b.op {
	case Star:
		fallthrough
	case Div:
		fallthrough
	case Mod:
		fallthrough
	case Plus:
		fallthrough
	case Minus:
		for _, t := range ts {
			if !t.Is(intType) {
				ctx.SemanticErrChan <- errors.NewTypeError(b.pos, op, t, intType)
				ok = false
			}
		}
	case Greater:
		fallthrough
	case GreaterEq:
		fallthrough
	case Less:
		fallthrough
	case LessEq:
		for _, t := range ts {
			if !t.Is(intType) && !t.Is(charType) {
				ctx.SemanticErrChan <- errors.NewMultiTypeError(b.pos, op, t, intType)
				ok = false
			}
		}
	case Equal:
		fallthrough
	case NotEq:
		for _, t := range ts {
			if t.Is(types.Array) {
				ctx.SemanticErrChan <- errors.NewMultiTypeError(b.pos, "equality operators", t, intType, boolType, charType, pairType)
				ok = false
			}
		}
	case And:
		fallthrough
	case Or:
		for _, t := range ts {
			if !t.Is(boolType) {
				ctx.SemanticErrChan <- errors.NewTypeError(b.pos, op, boolType, t)
				ok = false
			}
		}
	}

	//This check is done after the switch case to allow for better error messaging
	if !lt.Is(rt) {
		ctx.SemanticErrChan <- errors.NewSameTypeError(b.pos, op)
		ok = false
	}

	return ok
}

//String returns
// <BinopType>
//   - left_expr
//   - right_expr
func (b BinOp) String() string {
	return format(b.op.String(), b.left.String(), b.right.String())
}

//EvalType retuns integer or boolean depending on the binop
func (b BinOp) EvalType(_ symboltable.SymbolTable) types.WaccType {
	if b.op <= Minus {
		return types.Integer
	}
	return types.Boolean
}
