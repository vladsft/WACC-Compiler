package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
)

//VisitLeftIdent returns an Ident
func (w *WaccVisitor) VisitLeftIdent(ctx *parser.LeftIdentContext) interface{} {
	return ctx.Fieldident().Accept(w)
}

//VisitLeftArrayElem returns an ArrayElem
func (w *WaccVisitor) VisitLeftArrayElem(ctx *parser.LeftArrayElemContext) interface{} {
	return ctx.Arrayelem().Accept(w)
}

//VisitLeftPairType return a PairType
func (w *WaccVisitor) VisitLeftPairType(ctx *parser.LeftPairTypeContext) interface{} {
	return ctx.Pairtype().Accept(w)
}

//VisitLeftPairElem returns a PairElem
func (w *WaccVisitor) VisitLeftPairElem(ctx *parser.LeftPairElemContext) interface{} {
	return ctx.Pairelem().Accept(w)
}

//VisitSetlhs returns a PairElem
func (w *WaccVisitor) VisitSetlhs(ctx *parser.SetlhsContext) interface{} {
	if ident := ctx.Libident(); ident != nil {
		if exprsCtx := ctx.AllExpr(); len(exprsCtx) != 0 {
			exprs := make([]ast.Expression, len(exprsCtx))

			for i, exprCtx := range exprsCtx {
				exprs[i] = exprCtx.Accept(w).(ast.Expression)
			}
			pos := getPos(ctx)
			return ast.NewArrayElem(ident.Accept(w).(*ast.Ident), exprs, pos)
		}
		return ident.Accept(w)
	}
	if pairType := ctx.Pairtype(); pairType != nil {
		return pairType.Accept(w)
	}
	return ctx.Pairelem().Accept(w)
}
