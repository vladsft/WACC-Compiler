package ast

import (
	"fmt"
	"strings"
	"wacc_32/errors"
	"wacc_32/symboltable"
	"wacc_32/types"
)

var _ Expression = &Literal{}

//Literal literally represents a typed literal
type Literal struct {
	ast
	t     types.WaccType
	value interface{}
	pos   errors.Position
}

//NewLiteral returns a new type
func NewLiteral(t types.WaccType, value interface{}, pos errors.Position) *Literal {
	return &Literal{
		t:     t,
		value: value,
		pos:   pos,
	}
}

//NewDefaultLiteral returns a new type
func NewDefaultLiteral(t types.WaccType, pos errors.Position) *Literal {
	return &Literal{
		t:     t,
		value: t.DefaultValue(),
		pos:   pos,
	}
}

//GetValue returns the literal's value
func (l Literal) GetValue() interface{} {
	return l.value
}

//GetValue returns the literal's value
func (l Literal) GetType() types.WaccType {
	return l.t
}

//Check returns nil as TypeNodes are always valid
func (l *Literal) Check(ctx Context) bool {
	if l.t == types.Array {
		typeSet := make(map[types.WaccType]struct{})
		for _, expr := range l.value.([]Expression) {
			if !expr.Check(ctx) {
				return false
			}
			typeSet[expr.EvalType(*ctx.table)] = struct{}{}
			if len(typeSet) > 1 {
				ctx.SemanticErrChan <- errors.NewArrayTypeError(l.pos)
				return false
			}
		}
		//Get the subtype of the array
		var subType types.WaccType = types.Integer
		for k := range typeSet {
			subType = k
			break
		}
		l.t = types.NewArray(subType, 1)
	}
	if l.t.Is(types.UserDefinedType) {
		lt, err := LookupUserType(l.t.(types.UserType), *ctx.table)
		if err != nil {
			ctx.SemanticErrChan <- errors.NewUndefinedIdentifierError(l.pos, err)
		}
		if l.value == 0 {
			return true
		}

		if l.value == nil {
			ctx.SemanticErrChan <- errors.NewUninitialisedUsertypeError(l.pos)
			return false
		}

		fieldTypes := lt.GetChildren()
		fieldValues := l.value.([]Expression)
		structName := l.t.(types.UserType).GetName()

		if len(fieldValues) != len(fieldTypes) {
			ctx.SemanticErrChan <- errors.NewArgCountError(l.pos, "constructor "+structName, len(fieldTypes), len(fieldValues))
			return false
		}

		for i, expr := range fieldValues {
			if !expr.Check(ctx) {
				return false
			}
			// fmt.Println(expr.EvalType(*ctx.table))
			if !expr.EvalType(*ctx.table).Is(fieldTypes[i]) {
				nodeName := fmt.Sprintf("%s constructor argument number %d", structName, i+1)
				ctx.SemanticErrChan <- errors.NewTypeError(l.pos, nodeName, fieldTypes[i], expr.EvalType(*ctx.table))
				return false
			}
		}

	}

	return true
}

//EvalType returns the type of the literal
//Literals can only have basetypes as we can't infer their auxiliary types at parse time
func (t Literal) EvalType(_ symboltable.SymbolTable) types.WaccType {
	return t.t
}

//String returns
// <type> LITERAL
//   - value
func (t Literal) String() string {
	if t.t == types.Array {
		typeStr := strings.ToUpper(t.t.String()) + " LITERAL"

		str := ""
		first := true
		for _, v := range t.value.([]Expression) {
			if first {
				first = false
			} else {
				str += header
			}
			str += fmt.Sprintf("%s\n", v)
		}

		return format(typeStr, str)
	}

	if t.value == nil {
		return "null"
	}

	return fmt.Sprint(t.value)
}
