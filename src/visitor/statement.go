package visitor

import (
	"wacc_32/ast"
	"wacc_32/errors"
	"wacc_32/parser"
	"wacc_32/types"

	"github.com/antlr/antlr4/runtime/Go/antlr"
)

type hasPosition interface {
	GetStart() antlr.Token
	GetStop() antlr.Token
}

func getPos(ctx hasPosition) errors.Position {
	start := ctx.GetStart()
	stop := ctx.GetStop()

	return errors.NewPosition(start.GetLine(), start.GetColumn(), stop.GetLine(), stop.GetColumn())
}

//VisitStatSkip returns a StatSkip
func (w *WaccVisitor) VisitStatSkip(ctx *parser.StatSkipContext) interface{} {
	return ast.NewStatSkip()
}

//VisitStatFree returns a StatFree
func (w *WaccVisitor) VisitStatFree(ctx *parser.StatFreeContext) interface{} {
	freeExpr := ctx.Expr().Accept(w).(ast.Expression)
	pos := getPos(ctx)

	return ast.NewStatFree(freeExpr, pos)
}

//VisitStatIf returns a StatIf with correct condition and conditional statements
func (w *WaccVisitor) VisitStatIf(ctx *parser.StatIfContext) interface{} {
	cond := ctx.Expr().Accept(w).(ast.Expression)
	statsCtx := ctx.AllStat()
	ifStat := statsCtx[0].Accept(w).(ast.Statement)
	elseStat := statsCtx[1].Accept(w).(ast.Statement)
	pos := getPos(ctx)

	return ast.NewStatIf(cond, ifStat, elseStat, pos)
}

//VisitStatFor returns a StatFor with correct condition and body statements
func (w *WaccVisitor) VisitStatFor(ctx *parser.StatForContext) interface{} {
	initial := ctx.Newassign().Accept(w).(*ast.StatNewassign)
	cond := ctx.Expr().Accept(w).(ast.Expression)
	change := ctx.Assign().Accept(w).(*ast.StatAssign)
	bodyStat := ctx.Stat().Accept(w).(ast.Statement)

	return ast.NewStatFor(*initial, cond, *change, bodyStat)
}

//VisitStatWhile returns a StatWhile with correct condition and body statements
func (w *WaccVisitor) VisitStatWhile(ctx *parser.StatWhileContext) interface{} {
	cond := ctx.Expr().Accept(w).(ast.Expression)
	bodyStat := ctx.Stat().Accept(w).(ast.Statement)
	pos := getPos(ctx)
	
	return ast.NewStatWhile(cond, bodyStat, pos)
}

//VisitStatWhile returns a StatWhile with correct condition and body statements
func (w *WaccVisitor) VisitStatDoWhile(ctx *parser.StatDoWhileContext) interface{} {
	bodyStat := ctx.Stat().Accept(w).(ast.Statement)
	cond := ctx.Expr().Accept(w).(ast.Expression)

	return ast.NewStatDoWhile(bodyStat, cond)
}

//VisitStatBegin returns a StatBegin with the correct enclosed statement
func (w *WaccVisitor) VisitStatBegin(ctx *parser.StatBeginContext) interface{} {
	stat := ctx.Stat().Accept(w).(ast.Statement)

	return ast.NewStatBegin(stat)
}

//VisitStatPrintln returns a StatPrintln with the correct expression
func (w *WaccVisitor) VisitStatPrintln(ctx *parser.StatPrintlnContext) interface{} {
	exprToPrint := ctx.Expr().Accept(w).(ast.Expression)
	return ast.NewStatPrintln(exprToPrint)
}

//VisitStatPrint returns a StatPrint with the correct expression
func (w *WaccVisitor) VisitStatPrint(ctx *parser.StatPrintContext) interface{} {
	exprToPrint := ctx.Expr().Accept(w).(ast.Expression)
	return ast.NewStatPrint(exprToPrint)
}

//VisitStatAssign returns a StatAssign with the correct rhs ad lhs
func (w *WaccVisitor) VisitStatAssign(ctx *parser.StatAssignContext) interface{} {
	return ctx.Assign().Accept(w).(*ast.StatAssign)
}

func (w *WaccVisitor) VisitAssign(ctx *parser.AssignContext) interface{} {
	lhs := ctx.Setlhs().Accept(w).(ast.LHS)
	rhs := ctx.Assignrhs().Accept(w).(ast.RHS)
	pos := getPos(ctx)

	return ast.NewStatAssign(lhs, rhs, pos)
}

func (w *WaccVisitor) VisitStatDeclaration(ctx *parser.StatDeclarationContext) interface{} {
	return ctx.Declaration().Accept(w)
}

func (w *WaccVisitor) VisitNewassign(ctx *parser.NewassignContext) interface{} {
	ident := ctx.Ident().Accept(w).(*ast.Ident)
	t := ctx.Wacctype().Accept(w).(types.WaccType)
	pos := getPos(ctx)

	rhs := ctx.Assignrhs().Accept(w).(ast.Expression)

	return ast.NewStatNewassign(t, ident, rhs, pos)
}

//VisitStatNewassign returns a StatNewAssign with the correct expression
func (w *WaccVisitor) VisitStatNewassign(ctx *parser.StatNewassignContext) interface{} {
	return ctx.Newassign().Accept(w)
}

//VisitStatExit returns a StatExit with the correct exit code
func (w *WaccVisitor) VisitStatExit(ctx *parser.StatExitContext) interface{} {
	exitCode := ctx.Expr().Accept(w).(ast.Expression)
	pos := getPos(ctx)
	return ast.NewStatExit(exitCode, pos)
}

//VisitStatMultiple returns a StatMultiple with the correct exit code
func (w *WaccVisitor) VisitStatMultiple(ctx *parser.StatMultipleContext) interface{} {
	statsCtx := ctx.AllStat()
	stats := make([]ast.Statement, 0)

	for _, statCtx := range statsCtx {
		stat := statCtx.Accept(w).(ast.Statement)
		switch s := stat.(type) {
		case ast.StatMultiple:
			for _, st := range s {
				stats = append(stats, st)
			}
		default:
			stats = append(stats, stat)
		}
	}

	return ast.NewStatMultiple(stats)
}

//VisitStatReturn returns a StatReturn with the correct return value
func (w *WaccVisitor) VisitStatReturn(ctx *parser.StatReturnContext) interface{} {
	retValue := ctx.Expr().Accept(w).(ast.Expression)
	pos := getPos(ctx)

	return ast.NewStatReturn(retValue, pos)
}

//VisitStatRead returns a StatRead with the correct value to be read
func (w *WaccVisitor) VisitStatRead(ctx *parser.StatReadContext) interface{} {
	lhs := ctx.Assignlhs().Accept(w).(ast.LHS)
	pos := getPos(ctx)

	return ast.NewStatRead(lhs, pos)
}

//VisitStatEnhancedAssign returns the accumulating assignment
func (w *WaccVisitor) VisitStatEnhancedAssign(ctx *parser.StatEnhancedAssignContext) interface{} {

	lhs := ctx.Assignlhs().Accept(w).(ast.LHS)
	accRHS := ctx.Assignrhs().Accept(w).(ast.RHS)
	pos := getPos(ctx)
	var binopType ast.BinopType
	if op := ctx.ENH_PLUS(); op != nil {
		binopType = ast.Plus
	}
	if op := ctx.ENH_MINUS(); op != nil {
		binopType = ast.Minus
	}
	if op := ctx.ENH_STAR(); op != nil {
		binopType = ast.Star
	}
	if op := ctx.ENH_DIV(); op != nil {
		binopType = ast.Div
	}
	if op := ctx.ENH_MOD(); op != nil {
		binopType = ast.Mod
	}

	rhs := ast.NewBinOp(binopType, lhs, accRHS, pos)

	return ast.NewStatAssign(lhs, rhs, pos)
}
