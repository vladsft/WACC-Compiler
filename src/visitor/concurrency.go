package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
	"wacc_32/types"
)

//VisitStatWacc returns a function call in a new thread
func (w *WaccVisitor) VisitStatWacc(ctx *parser.StatWaccContext) interface{} {
	fName := ctx.Libident().Accept(w).(*ast.Ident)
	w.libMng.concurrentFunctions <- fName.GetName()
	arglist := make([]ast.Expression, 0)
	if args := ctx.Arglist(); args != nil {
		arglist = args.Accept(w).([]ast.Expression)
	}
	pos := getPos(ctx)
	if fName.IsNamespaced() {
		components := fName.GetNameComponents()
		thisStr := components[len(components)-2]
		thisArg := ast.NewIdent(thisStr, pos)
		arglist = append([]ast.Expression{thisArg}, arglist...)
	}
	return ast.NewWaccRoutine(ast.NewIdent("0"+fName.GetName(), fName.GetPos()), fName.IsNamespaced(), arglist, pos)
}

//VisitStatAcquire returns an acquire statement
func (w *WaccVisitor) VisitStatAcquire(ctx *parser.StatAcquireContext) interface{} {
	ident := ctx.Fieldident().Accept(w).(*ast.Ident)
	return ast.NewAcquire(ident, getPos(ctx))
}

//VisitStatRelease returns a release statement
func (w *WaccVisitor) VisitStatRelease(ctx *parser.StatReleaseContext) interface{} {
	ident := ctx.Fieldident().Accept(w).(*ast.Ident)
	return ast.NewRelease(ident, getPos(ctx))
}

func (w *WaccVisitor) VisitStatUp(ctx *parser.StatUpContext) interface{} {
	ident := ctx.Fieldident().Accept(w).(*ast.Ident)
	return ast.NewSemaUp(ident, getPos(ctx))
}

func (w *WaccVisitor) VisitStatDown(ctx *parser.StatDownContext) interface{} {
	ident := ctx.Fieldident().Accept(w).(*ast.Ident)
	return ast.NewSemaDown(ident, getPos(ctx))
}

func (w *WaccVisitor) VisitExprSemaLiter(ctx *parser.ExprSemaLiterContext) interface{} {
	return ctx.Semaliter().Accept(w)
}

func (w *WaccVisitor) VisitSemaliter(ctx *parser.SemaliterContext) interface{} {
	intLiter := ctx.Intliter().Accept(w).(*ast.Literal)
	return ast.NewLiteral(types.Sema, intLiter.GetValue(), getPos(ctx))
}
