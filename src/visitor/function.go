package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
	"wacc_32/types"
)

//VisitFunction retuns a Function with the correct signature and enclosed body statements
func (w *WaccVisitor) VisitFunction(ctx *parser.FunctionContext) interface{} {
	retType := ctx.Wacctype().Accept(w).(types.WaccType)
	ident := ctx.Ident().Accept(w).(*ast.Ident)

	params := ast.NewEmptyParamList()
	if paramsCtx := ctx.Paramlist(); paramsCtx != nil {
		params = paramsCtx.Accept(w).(ast.ParamList)
	}

	stats := ctx.Funcbody().Accept(w).(ast.Statement)
	pos := getPos(ctx)

	return ast.NewFunction(retType, ast.NewIdent("0"+w.importName+ident.String(), ident.GetPos()), params, stats, pos)
}

//VisitReturnable returns an return or exit statement based
func (w *WaccVisitor) VisitReturnable(ctx *parser.ReturnableContext) interface{} {
	expr := ctx.Expr().Accept(w).(ast.Expression)
	pos := getPos(ctx)
	if ctx.RETURN() != nil {
		return ast.NewStatReturn(expr, pos)
	}

	return ast.NewStatExit(expr, pos)
}

//VisitFuncbody returns a list of statements
func (w *WaccVisitor) VisitFuncbody(ctx *parser.FuncbodyContext) interface{} {
	statCtx := ctx.Stat()
	stats := make(ast.StatMultiple, 0)

	if statCtx != nil {
		stat := statCtx.Accept(w).(ast.Statement)
		stats = append(stats, stat)
	}

	returnableCtx := ctx.Returnable()
	if returnableCtx != nil {
		ret := returnableCtx.Accept(w).(ast.Statement)
		stats = append(stats, ret)
	} else { // if statement
		cond := ctx.Expr().Accept(w).(ast.Expression)

		fBodysCtx := ctx.AllFuncbody()
		ifFBodyCtx, elseFBodyCtx := fBodysCtx[0], fBodysCtx[1]

		ifStatList := ifFBodyCtx.Accept(w).(ast.StatMultiple)
		elseStatList := elseFBodyCtx.Accept(w).(ast.StatMultiple)

		ifStat := ast.NewStatMultiple(ifStatList)
		elseStat := ast.NewStatMultiple(elseStatList)
		pos := getPos(ifFBodyCtx)

		statIf := ast.NewStatIf(cond, ifStat, elseStat, pos)
		stats = append(stats, statIf)
	}

	return stats
}

//VisitArglist returns an an ArgList
func (w *WaccVisitor) VisitArglist(ctx *parser.ArglistContext) interface{} {
	exprsCtx := ctx.AllExpr()
	exprs := make([]ast.Expression, 0)

	for _, exprCtx := range exprsCtx {
		expr := exprCtx.Accept(w).(ast.Expression)

		exprs = append(exprs, expr)
	}

	return exprs
}

//VisitParamlist returns a ParamList
func (w *WaccVisitor) VisitParamlist(ctx *parser.ParamlistContext) interface{} {
	paramCtxs := ctx.AllParam()
	params := make([]*ast.Param, 0)

	for _, paramCtx := range paramCtxs {
		params = append(params, paramCtx.Accept(w).(*ast.Param))
	}
	return ast.NewParamList(params)
}

//VisitParam returns a Param
func (w *WaccVisitor) VisitParamNormal(ctx *parser.ParamNormalContext) interface{} {
	wType := ctx.Wacctype().Accept(w).(types.WaccType)
	ident := ctx.Ident().Accept(w).(*ast.Ident)
	pos := getPos(ctx)

	return ast.NewParam(wType, ident, pos)
}

//VisitParam returns a Param
func (w *WaccVisitor) VisitParamUserType(ctx *parser.ParamUserTypeContext) interface{} {
	libIdent := ctx.Libident().Accept(w).(*ast.Ident)
	wType := types.NewUserType(libIdent.GetName(), nil, nil, false)
	ident := ctx.Ident().Accept(w).(*ast.Ident)
	pos := getPos(ctx)

	return ast.NewParam(wType, ident, pos)
}
