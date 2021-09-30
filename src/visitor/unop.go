package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
	"wacc_32/types"
)

//VisitUnaryoper returns a UnaryOper
func (w *WaccVisitor) VisitUnaryoper(ctx *parser.UnaryoperContext) interface{} {
	if op := ctx.NOT(); op != nil {
		return ast.Not
	}
	if op := ctx.LEN(); op != nil {
		return ast.Len
	}
	if op := ctx.ORD(); op != nil {
		return ast.Ord
	}
	if op := ctx.CHR(); op != nil {
		return ast.Chr
	}
	if op := ctx.TRYLOCK(); op != nil {
		return ast.TryLock
	}
	return ast.Neg
}

//VisitStatTrailingUnOp returns the assignment of type i++
func (w *WaccVisitor) VisitStatTrailingUnOp(ctx *parser.StatTrailingUnOpContext) interface{} {
	ident := ctx.Fieldident().Accept(w).(*ast.Ident)
	pos := getPos(ctx)

	var binopType ast.BinopType
	if ctx.INC() != nil {
		binopType = ast.Plus
	} else {
		binopType = ast.Minus
	}

	oneLiter := ast.NewLiteral(types.Integer, 1, pos)
	rhs := ast.NewBinOp(binopType, ident, oneLiter, pos)

	return ast.NewStatAssign(ident, rhs, pos)
}
