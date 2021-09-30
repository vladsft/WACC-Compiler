package visitor

import (
	"strings"
	"wacc_32/ast"
	"wacc_32/parser"
	"wacc_32/types"
)

//VisitRightExpr returns an expression
func (w *WaccVisitor) VisitRightExpr(ctx *parser.RightExprContext) interface{} {
	expr := ctx.Expr().Accept(w).(ast.Expression)

	return expr
}

//VisitRightNewPair returns a new pair expression
func (w *WaccVisitor) VisitRightNewPair(ctx *parser.RightNewPairContext) interface{} {
	exprs := ctx.AllExpr()
	fst, snd := exprs[0].Accept(w), exprs[1].Accept(w)
	return ast.NewRHSNewPair(fst.(ast.Expression), snd.(ast.Expression))
}

//VisitRightPairElem returns a new pair element
func (w *WaccVisitor) VisitRightPairElem(ctx *parser.RightPairElemContext) interface{} {
	return ctx.Pairelem().Accept(w)
}

//VisitPairelem returns a new pair element
func (w *WaccVisitor) VisitPairelem(ctx *parser.PairelemContext) interface{} {
	pos := ast.FST
	if ctx.FST() == nil {
		pos = ast.SND
	}
	expr := ctx.Expr().Accept(w).(ast.Expression)
	return ast.NewPairElem(pos, expr)
}

//VisitRightFunctionCall returns a function call
func (w *WaccVisitor) VisitRightFunctionCall(ctx *parser.RightFunctionCallContext) interface{} {
	fName := ctx.Libident().Accept(w).(*ast.Ident)
	w.libMng.calledFunctions <- fName.GetName()

	arglist := make([]ast.Expression, 0)
	if args := ctx.Arglist(); args != nil {
		arglist = args.Accept(w).([]ast.Expression)
	}
	isMethod := false
	pos := getPos(ctx)
	if fName.IsNamespaced() {
		isMethod = true
		components := fName.GetNameComponents()
		thisStr := strings.Join(components[:len(components)-1], "!")
		if len(components) > 2 {
			thisStr = "1" + thisStr
		}
		thisArg := ast.NewIdent(thisStr, pos)
		arglist = append(arglist, thisArg)
	}

	return ast.NewRHSFunctionCall(ast.NewIdent("0"+fName.GetName(), fName.GetPos()), arglist, isMethod, pos)
}

//VisitRightArrayLiter return an array literal
func (w *WaccVisitor) VisitRightArrayLiter(ctx *parser.RightArrayLiterContext) interface{} {
	return ctx.Arrayliter().Accept(w)
}

func (w *WaccVisitor) VisitRightNewUserType(ctx *parser.RightNewUserTypeContext) interface{} {

	var values interface{}
	if valuesCtx := ctx.Arglist(); valuesCtx != nil {
		values = valuesCtx.Accept(w)
	}

	ident := ctx.Libident().Accept(w).(*ast.Ident)
	pos := getPos(ctx)

	utName := ident.GetName()
	utType := types.NewUserType(utName, nil, nil, false)

	return ast.NewLiteral(utType, values, pos)
}

func (w *WaccVisitor) VisitMake(ctx *parser.MakeContext) interface{} {
	makeType := ctx.Wacctype().Accept(w).(types.WaccType)
	lengthExpr := ctx.Expr().Accept(w).(ast.Expression)
	return ast.NewMake(makeType, lengthExpr)
}
