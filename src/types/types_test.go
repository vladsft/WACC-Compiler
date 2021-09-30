package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBaseTypeStringsMatch(t *testing.T) {
	for i := Integer; i < Function+1; i++ {
		assert.Equal(t, i.String(), typeStrings[i-1])
	}
}

func TestWaccTypeArrayStrings(t *testing.T) {
	for i := Integer; i < Array; i++ {
		arr := NewArray(i, 1)
		assert.Equal(t, i.String()+"[]", arr.String())
	}
}

func TestWaccTypeNestedArrayString(t *testing.T) {
	arr := NewArray(Integer, 1)
	arr1 := NewArray(arr, 1)

	assert.Equal(t, arr1.String(), "int[][]")
}

func TestWaccTypeEqualityHolds(t *testing.T) {
	base := Integer
	arr := NewArray(base, 1)
	arr1 := NewArray(arr, 1)
	arr2 := NewArray(arr1, 1)

	xbase := Integer
	xarr := NewArray(xbase, 1)
	xarr1 := NewArray(xarr, 1)
	xarr2 := NewArray(xarr1, 1)

	assert.True(t, base == xbase)
	assert.True(t, arr == xarr)
	assert.True(t, arr1 == xarr1)
	assert.True(t, arr2 == xarr2)

	assert.False(t, base == xarr)
	assert.False(t, arr == xarr1)
	assert.False(t, arr1 == xarr2)
	assert.False(t, arr2 == xbase)
}
