package ast

import (
	"fmt"
	"strings"
	"sync"
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

//Statement representing different statements in the Wacc language
type Statement AST

var (
	_ Statement = &StatSkip{}
	_ Statement = &StatRead{}
	_ Statement = &StatFree{}
	_ Statement = &StatNewassign{}
	_ Statement = &StatPrint{}
	_ Statement = &StatPrintln{}
	_ Statement = &StatExit{}
	_ Statement = &StatFor{}
	_ Statement = &StatWhile{}
	_ Statement = &StatDoWhile{}
	_ Statement = &StatBegin{}
	_ Statement = &StatAssign{}
	_ Statement = &StatReturn{}
	_ Statement = &StatIf{}
	_ Statement = &WaccRoutine{}
	_ Statement = StatMultiple{}
)

//StatSkip represents the skip statement
type StatSkip struct {
}

//NewStatSkip creates a new skip statement
func NewStatSkip() *StatSkip {
	return &StatSkip{}
}

//GetSymbolTable returns nil as skips don't need symbol tables
func (s StatSkip) GetSymbolTable() *symboltable.SymbolTable {
	return nil
}

//String returns
// SKIP
func (s StatSkip) String() string {
	return "SKIP"
}

//Check does nothing
func (s *StatSkip) Check(_ Context) {}

//LHS represents anything that can be on the Left Hand Side
type LHS Expression

//StatRead represents a read statement
type StatRead struct {
	ast
	toRead LHS
	pos    errors.Position
}

func (s StatRead) GetToRead() LHS {
	return s.toRead
}

//NewStatRead constructor
func NewStatRead(toRead LHS, pos errors.Position) *StatRead {
	return &StatRead{
		toRead: toRead,
		pos:    pos,
	}
}

//String returns
// READ
//   - <expr>
func (s StatRead) String() string {
	return format("READ", s.toRead.String())
}

//Check makes sure that it's LHS is either an integer or a character
func (s *StatRead) Check(ctx Context) {
	s.table = ctx.table
	if s.toRead.Check(ctx) {
		rType := s.toRead.EvalType(*ctx.table)
		if rType != types.Integer && rType != types.Char && rType != types.Str {
			ctx.SemanticErrChan <- errors.NewMultiTypeError(s.pos, "read", rType, types.Integer, types.Char, types.Str)
		}
	}
}

//StatFree represents the free statement
type StatFree struct {
	ast
	expr Expression
	pos  errors.Position
}

//NewStatFree creates a new free statement
func NewStatFree(expr Expression, pos errors.Position) *StatFree {
	return &StatFree{
		expr: expr,
		pos:  pos,
	}
}

//GetExpression returns the expression to free
func (s StatFree) GetExpression() Expression {
	return s.expr
}

//String returns
// FREE
//   - expr
func (s StatFree) String() string {
	return format("FREE", s.expr.String())
}

//Check checks the validity of the expression being freed and if the expression is a pair
func (s *StatFree) Check(ctx Context) {
	s.table = ctx.table
	if s.expr.Check(ctx) {
		freeType := s.expr.EvalType(*ctx.table)
		if !freeType.Is(types.Pair) && !freeType.Is(types.Array) && !freeType.Is(types.Lock) {
			ctx.SemanticErrChan <- errors.NewMultiTypeError(s.pos, "free", freeType, types.Pair, types.Array, types.Lock)
		}
	}
}

//StatNewassign represents a variable declaration
type StatNewassign struct {
	ast
	t     types.WaccType
	ident *Ident
	rhs   Expression
	pos   errors.Position
}

//NewStatNewassign creates a new StatNewassign
func NewStatNewassign(t types.WaccType, ident *Ident, rhs Expression,
	pos errors.Position) *StatNewassign {
	return &StatNewassign{
		t:     t,
		ident: ident,
		rhs:   rhs,
		pos:   pos,
	}
}

//GetType returns the type of the variable being declared
func (s StatNewassign) GetType() types.WaccType {
	return s.t
}

//GetRHS returns the right hand side of the declaration
func (s StatNewassign) GetRHS() RHS {
	return s.rhs
}

//GetName returns the name of the declaration
func (s StatNewassign) GetName() string {
	return s.ident.name
}

//String returns
// DECLARE
//   - TYPE
//       - <type>
//   - LHS
//       - ident
//   - RHS
//       - rhs
func (s StatNewassign) String() string {
	lhsStr := format("LHS", s.ident.String())
	rhsStr := format("RHS", s.rhs.String())
	return format("DECLARE", format("TYPE", s.t.String()), lhsStr, rhsStr)
}

//Check makes sure that the assignment can be added
//And that the types match
//And that the rhs is valid
func (s *StatNewassign) Check(ctx Context) {
	s.table = ctx.table
	s.ident.table = ctx.table

	if s.rhs.Check(ctx) {
		rType := s.rhs.EvalType(*ctx.table)
		if !s.t.Is(rType) {
			ctx.SemanticErrChan <- errors.NewTypeError(s.pos, "variable "+s.ident.name, s.t, rType)
		}
	}

	err := ctx.table.AddDeclaration(s.ident.name, s.t, s.pos)
	if err != nil {
		ctx.SemanticErrChan <- err
	}
}

//StatPrint represents a print statement
type StatPrint struct {
	ast
	exprToPrint Expression
	newLine     bool
}

//GetExprToPrint returns the expression that has to be printed
func (s StatPrint) GetExprToPrint() Expression {
	return s.exprToPrint
}

//NewStatPrint creates a new print statement
func NewStatPrint(exprToPrint Expression) *StatPrint {
	return &StatPrint{
		exprToPrint: exprToPrint,
		newLine:     false,
	}
}

//String returns
// PRINT
//   - expr
func (s StatPrint) String() string {
	return format("PRINT", s.exprToPrint.String())
}

//Check makes sure that the printed expression is valid
func (s *StatPrint) Check(ctx Context) {
	s.table = ctx.table
	s.exprToPrint.Check(ctx)
}

//StatPrintln is just a println statement
type StatPrintln StatPrint

//NewStatPrintln creates a new println statement
func NewStatPrintln(exprToPrint Expression) *StatPrintln {
	return &StatPrintln{
		exprToPrint: exprToPrint,
		newLine:     true,
	}
}

//GetExprToPrint returns the expression that has to be printed
func (s StatPrintln) GetExprToPrint() Expression {
	return s.exprToPrint
}

//String returns
// PRINTln
//   - expr
func (s StatPrintln) String() string {
	return format("PRINTLN", s.exprToPrint.String())
}

//Check makes sure that the printed expression is valid
func (s *StatPrintln) Check(ctx Context) {
	s.table = ctx.table
	s.exprToPrint.Check(ctx)
}

//StatExit represents an exit statement
type StatExit struct {
	ast
	exitCode Expression
	pos      errors.Position
}

//NewStatExit creates a new exit statement
func NewStatExit(ExitCode Expression, pos errors.Position) *StatExit {
	return &StatExit{
		exitCode: ExitCode,
		pos:      pos,
	}
}

//GetCode returns the exitcode expression
func (exit StatExit) GetCode() Expression {
	return exit.exitCode
}

//Check checks if the exit code is an integer and if it would evaluate correctly
func (exit *StatExit) Check(ctx Context) {
	exit.table = ctx.table
	if exit.exitCode.Check(ctx) {
		codeType := exit.exitCode.EvalType(*ctx.table)
		if codeType != types.Integer {
			ctx.SemanticErrChan <- errors.NewTypeError(exit.pos, "exit", types.Integer, codeType)
		}
	}
}

//String returns
// EXIT
//   - ExitCode
func (exit StatExit) String() string {
	return format("EXIT", exit.exitCode.String())
}

//StatFor represents a for loop
type StatFor struct {
	ast
	initial  StatNewassign
	cond     Expression
	change   StatAssign
	bodyStat Statement
}

//GetInitial returns the initial state of the iterator from the loop
func (s StatFor) GetInitial() StatNewassign {
	return s.initial
}

//GetCond returns the condition expression of the for loop
func (s StatFor) GetCond() Expression {
	return s.cond
}

//GetChange returns the statement which updates the iterator
func (s StatFor) GetChange() StatAssign {
	return s.change
}

//GetBody returns the body statements of the for loop
func (s StatFor) GetBody() Statement {
	return s.bodyStat
}

//NewStatFor creates a new for statement
func NewStatFor(initial StatNewassign, cond Expression, change StatAssign, bodyStat Statement) *StatFor {
	return &StatFor{
		initial:  initial,
		cond:     cond,
		change:   change,
		bodyStat: bodyStat,
	}
}

//String returns
// FOR
//	 - INITIAL VALUE
//		 - initial
//   - CONDITION
//       - cond
//	 - UPDATED VALUE
//   	 - change
//   - DO
//       - bodyStat
func (s StatFor) String() string {
	initialStr := format("INITIAL VALUE", s.initial.String())
	condStr := format("CONDITION", s.cond.String())
	changeStr := format("UPDATED VALUE", s.change.String())
	bodyStr := format("DO", s.bodyStat.String())
	str := format("FOR", initialStr, condStr, changeStr, bodyStr)
	return str
}

//Check ensures that the for parameters are valid and that the body is valid.
func (s *StatFor) Check(ctx Context) {
	s.table = ctx.table

	//Check ensures that the initial statement is valid
	s.initial.Check(ctx)

	//Check ensures that the condition a valid boolean expression
	if s.cond.Check(ctx) {
		boolean := types.Boolean
		if s.cond.EvalType(*ctx.table) != boolean {
			ctx.SemanticErrChan <- fmt.Errorf("Condition should be a boolean type")
		}
	}

	//Check ensures that the change statement is valid
	s.change.Check(ctx)

	forCtx := Context{
		SemanticErrChan: ctx.SemanticErrChan,
		functionName:    ctx.functionName,
		table:           symboltable.NewSymbolTable(ctx.table),
		returnType:      ctx.returnType,
	}

	//Check ensures that the body is valid
	s.bodyStat.Check(forCtx)
	s.table.SetTotalOffset(forCtx.table.GetTotalOffset())
}

//StatWhile represents a while loop
type StatWhile struct {
	ast
	cond     Expression
	bodyStat Statement
	pos      errors.Position
}

//GetCond returns the condition expression of the while loop
func (s StatWhile) GetCond() Expression {
	return s.cond
}

//GetBody returns the body statements of the while loop
func (s StatWhile) GetBody() Statement {
	return s.bodyStat
}

//NewStatWhile creates a new while statement
func NewStatWhile(cond Expression, bodyStat Statement, pos errors.Position) *StatWhile {
	return &StatWhile{
		cond:     cond,
		bodyStat: bodyStat,
		pos:      pos,
	}
}

//String returns
// LOOP
//   - CONDITION
//       - cond
//   - DO
//       - bodyStat
func (s StatWhile) String() string {
	condStr := format("CONDITION", s.cond.String())
	bodyStr := format("DO", s.bodyStat.String())
	str := format("LOOP", condStr, bodyStr)
	return str
}

//Check ensures that the condition is a valid boolean expression and that the body is valid
func (s *StatWhile) Check(ctx Context) {
	s.table = ctx.table
	if s.cond.Check(ctx) {
		boolean := types.Boolean
		condType := s.cond.EvalType(*ctx.table)
		if condType != boolean {
			ctx.SemanticErrChan <- errors.NewTypeError(s.pos, "condition", boolean, condType)
		}
	}

	whileCtx := Context{
		SemanticErrChan: ctx.SemanticErrChan,
		functionName:    ctx.functionName,
		table:           symboltable.NewSymbolTable(ctx.table),
		returnType:      ctx.returnType,
	}

	s.bodyStat.Check(whileCtx)
	s.table.SetTotalOffset(whileCtx.table.GetTotalOffset())
}

type StatDoWhile struct {
	ast
	bodyStat Statement
	cond     Expression
}

//GetCond returns the condition expression of the do while loop
func (s StatDoWhile) GetCond() Expression {
	return s.cond
}

//GetBody returns the body statements of the do while loop
func (s StatDoWhile) GetBody() Statement {
	return s.bodyStat
}

//NewStatWhile creates a new while statement
func NewStatDoWhile(bodyStat Statement, cond Expression) *StatDoWhile {
	return &StatDoWhile{
		bodyStat: bodyStat,
		cond:     cond,
	}
}

//String returns
// LOOP
//   - DO
//       - bodyStat
//   - CONDITION
//       - cond
func (s StatDoWhile) String() string {
	bodyStr := format("DO", s.bodyStat.String())
	condStr := format("CONDITION", s.cond.String())
	str := format("LOOP", bodyStr, condStr)
	return str
}

//Check ensures that the condition is a valid boolean expression and that the body is valid
func (s *StatDoWhile) Check(ctx Context) {
	s.table = ctx.table
	doWhileCtx := Context{
		SemanticErrChan: ctx.SemanticErrChan,
		functionName:    ctx.functionName,
		table:           symboltable.NewSymbolTable(ctx.table),
		returnType:      ctx.returnType,
	}

	s.bodyStat.Check(doWhileCtx)
	if s.cond.Check(ctx) {
		boolean := types.Boolean
		if s.cond.EvalType(*ctx.table) != boolean {
			ctx.SemanticErrChan <- fmt.Errorf("condition should be a boolean type")
		}
	}

	s.table.SetTotalOffset(doWhileCtx.table.GetTotalOffset())
}

//StatBegin represents a begin-end statement
type StatBegin struct {
	ast
	stat Statement
}

//GetStat returns the statement in the local scope
func (s StatBegin) GetStat() Statement {
	return s.stat
}

//NewStatBegin creates a new scope
func NewStatBegin(stat Statement) *StatBegin {
	return &StatBegin{
		stat: stat,
	}
}

//String returns
// BEGIN
//   - stat
// END
func (s StatBegin) String() string {
	return format("SCOPE", s.stat.String())
}

//Check - When the programs begin, make a symbol table with the top level
//symbol table as parent
func (s *StatBegin) Check(ctx Context) {
	s.table = ctx.table
	st := symboltable.NewSymbolTable(ctx.table)
	s.stat.Check(Context{
		functionName:    ctx.functionName,
		table:           st,
		SemanticErrChan: ctx.SemanticErrChan,
	})
}

//StatAssign represents an assignment statement
type StatAssign struct {
	ast
	lhs LHS
	rhs RHS
	pos errors.Position
}

//NewStatAssign creates a new StatAssign
func NewStatAssign(lhs LHS, rhs RHS, pos errors.Position) *StatAssign {
	return &StatAssign{
		lhs: lhs,
		rhs: rhs,
		pos: pos,
	}
}

//GetLHS returns the variable being assigned
func (s StatAssign) GetLHS() Expression {
	return s.lhs
}

//GetRHS returns the expression used to assign
func (s StatAssign) GetRHS() Expression {
	return s.rhs
}

//String returns
// ASSIGN
//   - LHS
//       - ident
//   - RHS
//       - rhs
func (s StatAssign) String() string {
	lhsStr := format("LHS", s.lhs.String())
	rhsStr := format("RHS", s.rhs.String())
	return format("ASSIGNMENT", lhsStr, rhsStr)
}

//Check lhs and rhs compatibility
func (s *StatAssign) Check(ctx Context) {
	s.table = ctx.table
	if s.lhs.Check(ctx) && s.rhs.Check(ctx) {
		lType := s.lhs.EvalType(*ctx.table)
		rType := s.rhs.EvalType(*ctx.table)
		if !lType.Is(rType) {
			ctx.SemanticErrChan <- errors.NewTypeError(s.pos, "assignment", lType, rType)
		}
	}
}

//StatReturn represents a return statement
type StatReturn struct {
	ast
	caller   string
	retValue Expression
	pos      errors.Position
}

//NewStatReturn creates a new return statement
func NewStatReturn(retValue Expression, pos errors.Position) *StatReturn {
	return &StatReturn{
		retValue: retValue,
		pos:      pos,
	}
}

func (s StatReturn) GetFunctionName() string {
	return s.caller
}

//GetReturnExpr returns an expression
func (s StatReturn) GetReturnExpr() Expression {
	return s.retValue
}

//String returns
// RETURN
//   - retValue
func (s StatReturn) String() string {
	return format("RETURN", s.retValue.String())
}

// Check makes sure that the return is in the correct context and that the returned expression is
// correct
func (s *StatReturn) Check(ctx Context) {
	s.table = ctx.table
	s.caller = ctx.functionName
	if ctx.functionName == "" {
		ctx.SemanticErrChan <- errors.NewReturnError(s.pos)
		return
	}
	if s.retValue.Check(ctx) {
		retType := s.retValue.EvalType(*ctx.table)
		if !ctx.returnType.Is(retType) {
			ctx.SemanticErrChan <- errors.NewTypeError(s.pos, "return", ctx.returnType, retType)
		}
	}
}

//StatIf represents an If statement
type StatIf struct {
	ast
	cond     Expression
	ifStat   Statement
	elseStat Statement
	pos      errors.Position
}

//GetCondition gets the if statement's condition
func (s StatIf) GetCondition() Expression {
	return s.cond
}

//GetIfStat gets the statement found in the line with the "if" instruction
func (s StatIf) GetIfStat() Statement {
	return s.ifStat
}

//GetElseStat gets the statement found in the line with the "else" instruction
func (s StatIf) GetElseStat() Statement {
	return s.elseStat
}

//NewStatIf creates an if statement
func NewStatIf(cond Expression, ifStat Statement, elseStat Statement, pos errors.Position) *StatIf {
	return &StatIf{
		cond:     cond,
		ifStat:   ifStat,
		elseStat: elseStat,
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
func (s StatIf) String() string {
	condStr := format("CONDITION", s.cond.String())
	ifStatStr := format("THEN", s.ifStat.String())
	elseStatStr := format("ELSE", s.elseStat.String())

	return format("IF", condStr, ifStatStr, elseStatStr)
}

//Check ensures that the condition is a boolean and that the statements are valid
func (s *StatIf) Check(ctx Context) {
	s.table = ctx.table
	if s.cond.Check(ctx) {
		condType := s.cond.EvalType(*ctx.table)
		if condType != types.Boolean {
			ctx.SemanticErrChan <- errors.NewTypeError(s.pos, "if condition", types.Boolean, condType)
		}
	}
	branches := [2]Statement{s.ifStat, s.elseStat}
	offsets := [2]int{}

	//Concurrent Semantic analysis of the if else branches
	var wg sync.WaitGroup
	for i, branch := range branches {
		wg.Add(1)
		go func(i int, branch Statement) {
			defer wg.Done()
			branchCtx := Context{
				functionName:    ctx.functionName,
				returnType:      ctx.returnType,
				table:           symboltable.NewSymbolTable(ctx.table),
				SemanticErrChan: ctx.SemanticErrChan,
			}
			branch.Check(branchCtx)
			offsets[i] = branchCtx.table.GetTotalOffset()
		}(i, branch)
	}
	wg.Wait()
	offset := offsets[1]
	if offsets[0] > offsets[1] {
		offset = offsets[0]
	}
	s.table.SetTotalOffset(offset)
}

//A WaccRoutine is just a wrapper around a RHSFunctionCall
type WaccRoutine struct {
	*RHSFunctionCall
}

//NewWaccRoutine creates a FunctionCall which will be executed in a new thread
func NewWaccRoutine(fName *Ident, isMethod bool, args []Expression, pos errors.Position) *WaccRoutine {
	call := NewRHSFunctionCall(fName, args, isMethod, pos)
	call.concurrent = true
	return &WaccRoutine{call}
}

//Check makes sure the underlying WaccRoutine is semantically correct
func (wr *WaccRoutine) Check(ctx Context) {
	wr.RHSFunctionCall.Check(ctx)
}

//StatMultiple represents multiple statements
type StatMultiple []Statement

//String returns
// - stat1
// ...
// - statn
func (s StatMultiple) String() string {
	var str string
	first := true
	for _, sts := range s {
		if first {
			first = false
		} else {
			str += header
		}
		str += sts.String() + "\n"
	}
	return strings.TrimSpace(str)
}

//NewStatMultiple constructor
func NewStatMultiple(stats []Statement) StatMultiple {
	return StatMultiple(stats)
}

//Check checks each statement
func (s StatMultiple) Check(ctx Context) {
	for _, child := range s {
		child.Check(ctx)
	}
}

//GetSymbolTable returns nil
func (s StatMultiple) GetSymbolTable() *symboltable.SymbolTable {
	return nil
}

type EnhancedBinopType int

//Enum representing the different binary operators
const (
	EnhPlus EnhancedBinopType = iota + 1
	EnhMinus
	EnhStar
	EnhDiv
	EnhMod
)

var enhancedBinopStrings = []string{"+=", "-=", "*=", "/=", "%="}

//String returns the encoding of a binop
func (e EnhancedBinopType) String() string {
	return binopStrings[e-1]
}

//StatEnhOperator represents an accumulator statement
type StatEnhancedAssign struct {
	ast
	lhs LHS
	op  EnhancedBinopType
	rhs Expression
	pos errors.Position
}

//NewStatEnhOperator creates a new StatEnhOperator
func NewStatEnhancedAssign(lhs LHS, op EnhancedBinopType, rhs RHS, pos errors.Position) *StatEnhancedAssign {
	return &StatEnhancedAssign{
		lhs: lhs,
		op:  op,
		rhs: rhs,
		pos: pos,
	}
}

//GetLHS returns the variable being assigned
func (s StatEnhancedAssign) GetIdent() Expression {
	return s.lhs
}

//GetRHS returns the expression used to assign
func (s StatEnhancedAssign) GetRHS() Expression {
	return s.rhs
}

//GetOp returns the operatorof the accumulating assignment
func (s StatEnhancedAssign) GetOp() EnhancedBinopType {
	return s.op
}

//String returns
// ENHANCED OPERATOR
//   - LHS
//       - lhs
//   - RHS
//       - rhs
func (s StatEnhancedAssign) String() string {
	return format(s.op.String(), s.lhs.String(), s.rhs.String())
}

//Check lhs and rhs compatibility
func (s *StatEnhancedAssign) Check(ctx Context) {
	s.table = ctx.table
	if s.lhs.Check(ctx) && s.rhs.Check(ctx) {

		lType := s.lhs.EvalType(*ctx.table)
		rType := s.rhs.EvalType(*ctx.table)

		if !(rType.Is(types.Integer) && lType.Is(rType)) {
			ctx.SemanticErrChan <- errors.NewTypeError(s.pos, "accumulator assignment", types.Integer, rType)
		}
	}
}
