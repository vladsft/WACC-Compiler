package symboltable

import (
	"testing"
	"wacc_32/errors"
	"wacc_32/types"

	"github.com/stretchr/testify/assert"
)

var pos = errors.NewPosition(0, 0, 0, 0)

func TestAddDeclarationWorksInTopLevel(t *testing.T) {
	st := NewTopSymbolTable()

	intType := types.Integer

	err := st.AddDeclaration("wacc", intType, pos)
	metadata, ok := st.declarations["wacc"]

	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, metadata.wt, intType)
	assert.Equal(t, metadata.pos, pos)
}

func TestAddRepeatedDeclarationThrowsAnError(t *testing.T) {
	st := NewTopSymbolTable()

	intType := types.Integer

	st.AddDeclaration("wacc", intType, pos)
	err := st.AddDeclaration("wacc", intType, pos)

	assert.Error(t, err)
}

func TestCanAddDeclarationInLowerLevel(t *testing.T) {
	st := NewTopSymbolTable()
	st2 := NewSymbolTable(st)

	intType := types.Integer

	st.AddDeclaration("wacc", intType, pos)
	err := st2.AddDeclaration("wacc", intType, pos)

	assert.Nil(t, err)
}

func TestCanGetIdeantInCurrentLevel(t *testing.T) {
	st := NewTopSymbolTable()

	intType := types.Integer

	st.AddDeclaration("wacc", intType, pos)
	wt, err := st.GetType("wacc")

	assert.Nil(t, err)
	assert.Equal(t, wt, intType)
}

func TestCanGetIdentInHigherLevel(t *testing.T) {
	st := NewTopSymbolTable()
	st2 := NewSymbolTable(st)

	intType := types.Integer

	st.AddDeclaration("wacc", intType, pos)

	wt, err := st2.GetType("wacc")

	assert.Nil(t, err)
	assert.Equal(t, wt, intType)
}
func TestCannotGetIdentInParallelLevel(t *testing.T) {
	st := NewTopSymbolTable()
	st2 := NewSymbolTable(st)
	st3 := NewSymbolTable(st)

	intType := types.Integer

	st2.AddDeclaration("wacc", intType, pos)

	_, err := st3.GetType("wacc")

	assert.Error(t, err)
}
