package ast

import (
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

var _ Expression = &TernaryOp{}

type TernaryOp struct {
	ast
	cond, ifExpr, elseExpr Expression
	pos                    errors.Position
}

//GetCondition gets the ternary's statement condition
func (t TernaryOp) GetCondition() Expression {
	return t.cond
}

//GetIfExpr gets the expression of the "else" branch
func (t TernaryOp) GetIfExpr() Expression {
	return t.ifExpr
}

//GetElseExpr gets the expression of the "if" branch
func (t TernaryOp) GetElseExpr() Expression {
	return t.elseExpr
}

//NewTernarynOp creates a new TernaryOp
func NewTernaryOp(cond, ifExpr, elseExpr Expression, pos errors.Position) *TernaryOp {
	return &TernaryOp{
		cond:     cond,
		ifExpr:   ifExpr,
		elseExpr: elseExpr,
		pos:      pos,
	}
}

//String returns
// IF
//   - CONDITION
//       - expr
//   - THEN
//       - stat
//   - ELSE
//       - stat
func (t TernaryOp) String() string {

	condStr := format("CONDITION", t.cond.String())
	ifExprStr := format("THEN", t.ifExpr.String())
	elseExprStr := format("ELSE", t.elseExpr.String())

	return format("IF", condStr, ifExprStr, elseExprStr)
}

func (t TernaryOp) EvalType(s symboltable.SymbolTable) types.WaccType {
	//check truth value of cond - how??
	//if true, return ifExpr type
	//if false, return elseExpr type
	/*	condType := t.cond.EvalType(s)
		switch condType {
		case *ast.Literal:
			if condType.GetValue().(bool) {
				return t.ifExpr.EvalType(s)
			} else {
				return t.elseExpr.EvalType(s)
			}
		}*/
	return t.ifExpr.EvalType(s)
}

func (t *TernaryOp) Check(ctx Context) bool {
	t.table = ctx.table
	ok := true
	if t.cond.Check(ctx) {
		condType := t.cond.EvalType(*ctx.table)
		if condType != types.Boolean {
			ctx.SemanticErrChan <- errors.NewTypeError(t.pos, "if condition", types.Boolean, condType)
			ok = false
		}
	}
	// check if both are the same type
	if !concurrentCheck(ctx, t.ifExpr, t.elseExpr) {
		return false
	}
	if !t.ifExpr.EvalType(*ctx.table).Is(t.elseExpr.EvalType(*ctx.table)) {
		ctx.SemanticErrChan <- errors.NewTernaryError(t.pos, t.ifExpr.EvalType(*ctx.table), t.elseExpr.EvalType(*ctx.table))
		ok = false
	}
	return ok
}

//rt := t.elseExpr.EvalType(*ctx.table)

// branches := [2]Expression{t.ifExpr, t.elseExpr}
// offsets := [2]int{}

// //Concurrent Semantic analysis of the if else branches
// var wg sync.WaitGroup
// for i, branch := range branches {
// 	wg.Add(1)
// 	go func(i int, branch Statement) {
// 		defer wg.Done()
// 		branchCtx := Context{
// 			functionName:    ctx.functionName,
// 			returnType:      ctx.returnType,
// 			table:           symboltable.NewSymbolTable(ctx.table),
// 			SemanticErrChan: ctx.SemanticErrChan,
// 		}
// 		branch.Check(branchCtx)
// 		offsets[i] = branchCtx.table.GetTotalOffset()
// 	}(i, branch)
// }
// wg.Wait()
// offset := offsets[1]
// if offsets[0] > offsets[1] {
// 	offset = offsets[0]
// }
// t.table.SetTotalOffset(offset)
