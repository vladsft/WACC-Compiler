package visitor

import (
	"strconv"
	"wacc_32/ast"
	"wacc_32/parser"
	"wacc_32/types"
)

//VisitExprIdent returns the name of the identifier
func (w *WaccVisitor) VisitExprIdent(ctx *parser.ExprIdentContext) interface{} {
	return ctx.Fieldident().Accept(w)
}

//VisitExprCharLiter returns a TypeNode with the correct type and value
func (w *WaccVisitor) VisitExprCharLiter(ctx *parser.ExprCharLiterContext) interface{} {
	charLiter := ctx.CHAR_LITER().GetText()
	return ast.NewLiteral(types.Char, charLiter, getPos(ctx))
}

//VisitExprIntLiter returns a TypeNode with the correct type and value
func (w *WaccVisitor) VisitExprIntLiter(ctx *parser.ExprIntLiterContext) interface{} {
	return ctx.Intliter().Accept(w)
}

//VisitExprStringLiter returns a TypeNode with the correct type and value
func (w *WaccVisitor) VisitExprStringLiter(ctx *parser.ExprStringLiterContext) interface{} {
	strLiter := ctx.STRING_LITER().GetText()
	return ast.NewLiteral(types.Str, strLiter, getPos(ctx))
}

//VisitExprBoolLiter returns a TypeNode with the correct type and value
func (w *WaccVisitor) VisitExprBoolLiter(ctx *parser.ExprBoolLiterContext) interface{} {
	b, _ := strconv.ParseBool(ctx.BOOL_LITER().GetText())
	return ast.NewLiteral(types.Boolean, b, getPos(ctx))
}

//VisitExprPairLiter returns a pair with no type and value
func (w *WaccVisitor) VisitExprPairLiter(ctx *parser.ExprPairLiterContext) interface{} {
	return ctx.Pairliter().Accept(w)
}

//VisitIntliter returns an integer literal
func (w *WaccVisitor) VisitIntliter(ctx *parser.IntliterContext) interface{} {
	i, _ := strconv.Atoi(ctx.GetText())
	return ast.NewLiteral(types.Integer, i, getPos(ctx))
}

//VisitArrayliter covers all expressions of the form [...]
//If the array is empty, it doesn't assume a basetype, this should be done later
func (w *WaccVisitor) VisitArrayliter(ctx *parser.ArrayliterContext) interface{} {
	exprsCtx := ctx.AllExpr()
	values := make([]ast.Expression, len(exprsCtx))
	for i, element := range exprsCtx {
		values[i] = element.Accept(w).(ast.Expression)
	}
	return ast.NewLiteral(types.Array, values, getPos(ctx))
}

//VisitPairtype returns pair which is a WaccBasType
func (w *WaccVisitor) VisitPairtype(ctx *parser.PairtypeContext) interface{} {
	fst := ctx.Pairelemtype(0).Accept(w).(types.WaccType)
	snd := ctx.Pairelemtype(1).Accept(w).(types.WaccType)
	return types.NewPair(fst, snd)
}

//VisitPairliter returns an empty pair literal
func (w *WaccVisitor) VisitPairliter(ctx *parser.PairliterContext) interface{} {
	return ast.NewLiteral(types.Pair, nil, getPos(ctx))
}

//VisitPairelemtype returns a WaccBaseType
func (w *WaccVisitor) VisitPairelemtype(ctx *parser.PairelemtypeContext) interface{} {
	if bType := ctx.Basetype(); bType != nil {
		return bType.Accept(w)
	}
	if aType := ctx.Arraytype(); aType != nil {
		return aType.Accept(w)
	}
	return types.Pair
}

//VisitBasetype returns the correct WacBaseType
func (w *WaccVisitor) VisitBasetype(ctx *parser.BasetypeContext) interface{} {
	if intType := ctx.INT(); intType != nil {
		return types.Integer
	}
	if boolType := ctx.BOOL(); boolType != nil {
		return types.Boolean
	}
	if charType := ctx.CHAR(); charType != nil {
		return types.Char
	}
	if lockType := ctx.LOCK(); lockType != nil {
		return types.Lock
	}
	if semaType := ctx.SEMA(); semaType != nil {
		return types.Sema
	}
	return types.Str
}

//VisitArraytype returns a WaccBaseType
func (w *WaccVisitor) VisitArraytype(ctx *parser.ArraytypeContext) interface{} {
	dim := len(ctx.AllLBRACKET())
	if bType := ctx.Basetype(); bType != nil {
		return types.NewArray(bType.Accept(w).(types.WaccType), dim)
	}
	pType := ctx.Pairtype().Accept(w)

	return types.NewArray(pType.(types.WaccType), dim)
}

//VisitWacctype returns the correct WaccType
func (w *WaccVisitor) VisitWacctype(ctx *parser.WacctypeContext) interface{} {
	if libCtx := ctx.Libident(); libCtx != nil {
		libIdent := libCtx.Accept(w).(*ast.Ident)
		return types.NewUserType(libIdent.GetName(), nil, nil, false)
	}
	return w.VisitChildren(ctx).([]interface{})[0].(types.WaccType)
}
