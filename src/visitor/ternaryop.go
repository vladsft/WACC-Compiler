package visitor

import (
	"wacc_32/ast"
	"wacc_32/parser"
)

func (w *WaccVisitor) VisitExprTernaryOp(ctx *parser.ExprTernaryOpContext) interface{} {
	exprs := ctx.AllExpr()
	cond, ifExpr, elseExpr := exprs[0].Accept(w).(ast.Expression), exprs[1].Accept(w).(ast.Expression), exprs[2].Accept(w).(ast.Expression)
	/*cond := ctx.GetExpr(0).Accept(w).(ast.Expression)
	ifexpr := ctx.GetExpr(1).Accept(w).(ast.Expression)
	elseexpr := ctx.GetExpr(2).Accept(w).(ast.Expression)*/
	pos := getPos(ctx)

	return ast.NewTernaryOp(cond, ifExpr, elseExpr, pos)

}

// //VisitStatIf returns a StatIf with correct condition and conditional statements
// func (w *WaccVisitor) VisitStatIf(ctx *parser.StatIfContext) interface{} {
// 	cond := ctx.Expr().Accept(w).(ast.Expression)
// 	statsCtx := ctx.AllStat()
// 	ifStat := statsCtx[0].Accept(w).(ast.Statement)
/// 	elseStat := statsCtx[1].Accept(w).(ast.Statement)
// 	pos := getPos(ctx)

// 	return ast.NewStatIf(cond, ifStat, elseStat, pos)
// }
