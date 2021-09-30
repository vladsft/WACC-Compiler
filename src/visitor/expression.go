package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
)

//VisitExprBracketed reduces to the inner expression
func (w *WaccVisitor) VisitExprBracketed(ctx *parser.ExprBracketedContext) interface{} {
	return ctx.Expr().Accept(w)
}

//VisitExprArrayElem returns an ArrayElem
func (w *WaccVisitor) VisitExprArrayElem(ctx *parser.ExprArrayElemContext) interface{} {
	return ctx.Arrayelem().Accept(w)
}

//VisitArrayelem returns an NewArrayElem
func (w *WaccVisitor) VisitArrayelem(ctx *parser.ArrayelemContext) interface{} {
	ident := ctx.Fieldident().Accept(w).(*ast.Ident)
	indicesCtx := ctx.AllExpr()
	indices := make([]ast.Expression, len(indicesCtx))

	for i, iCtx := range indicesCtx {
		indices[i] = iCtx.Accept(w).(ast.Expression)
	}
	pos := getPos(ctx)
	return ast.NewArrayElem(ident, indices, pos)
}

//VisitExprUnaryOp constructs a unary operator
func (w *WaccVisitor) VisitExprUnaryOp(ctx *parser.ExprUnaryOpContext) interface{} {
	op := ctx.Unaryoper().Accept(w).(ast.UnopType)
	expr := ctx.Expr().Accept(w).(ast.Expression)
	pos := getPos(ctx)
	return ast.NewUnOp(op, expr, pos)
}