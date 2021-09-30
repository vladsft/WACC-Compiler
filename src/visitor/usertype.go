package visitor

import (
	"strings"
	"wacc_32/ast"
	"wacc_32/parser"
	"wacc_32/types"
)

func (w *WaccVisitor) VisitUserType(ctx *parser.UserTypeContext) interface{} {
	ident := ctx.Ident().Accept(w).(*ast.Ident)
	ident = ast.NewIdent(w.importName+ident.GetName(), ident.GetPos())

	fieldsCtx := ctx.AllDeclaration()
	fields := make([]*ast.StatNewassign, len(fieldsCtx))

	for i, fieldCtx := range fieldsCtx {
		fields[i] = fieldCtx.Accept(w).(*ast.StatNewassign)
	}

	isClass := ctx.CLASS() != nil

	funcsCtx := ctx.AllFunction()
	funcs := make([]*ast.Function, len(funcsCtx))
	for i, funcCtx := range funcsCtx {
		fn := funcCtx.Accept(w).(*ast.Function)
		fn.MakeMethod(types.NewUserType(ident.GetName(), nil, nil, false))
		w.libMng.functions <- fn
		funcs[i] = fn
	}
	userType := ast.NewUserType(ident, fields, isClass, funcs)
	return userType
}

func (w *WaccVisitor) VisitDeclaration(ctx *parser.DeclarationContext) interface{} {
	ident := ctx.Ident().Accept(w).(*ast.Ident)

	t := ctx.Wacctype().Accept(w).(types.WaccType)
	pos := getPos(ctx)

	rhs := ast.NewDefaultLiteral(t, pos)
	return ast.NewStatNewassign(t, ident, rhs, pos)
}

func (w *WaccVisitor) VisitFieldident(ctx *parser.FieldidentContext) interface{} {
	name := ctx.Ident().GetText()
	if ctx.Fieldident() != nil {
		accessedName := ctx.Fieldident().GetText()
		name = "1" + name + "!" + strings.Replace(accessedName, ".", "!", -1)
	}
	return ast.NewIdent(name, getPos(ctx))
}
