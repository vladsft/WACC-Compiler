package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
)

//VisitExprBinop visits a binary operator
func (w *WaccVisitor) VisitExprBinop(ctx *parser.ExprBinopContext) interface{} {
	l := ctx.GetLeft().Accept(w).(ast.Expression)
	r := ctx.GetRight().Accept(w).(ast.Expression)
	opStr := ctx.GetOp().GetText()
	pos := getPos(ctx)
	return ast.NewBinOp(stringToBinop[opStr], l, r, pos)
}

var stringToBinop = map[string]ast.BinopType{
	"*":  ast.Star,
	"/":  ast.Div,
	"%":  ast.Mod,
	"+":  ast.Plus,
	"-":  ast.Minus,
	">":  ast.Greater,
	">=": ast.GreaterEq,
	"<":  ast.Less,
	"<=": ast.LessEq,
	"==": ast.Equal,
	"!=": ast.NotEq,
	"&&": ast.And,
	"||": ast.Or}
