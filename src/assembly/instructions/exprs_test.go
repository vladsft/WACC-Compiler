package instructions

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCondOpposite(t *testing.T) {
	testCases := []struct {
		desc   string
		c1, c2 Cond
	}{
		{
			desc: "Equality",
			c1:   EQ,
			c2:   NE,
		},
		{
			desc: "Greater than",
			c1:   LE,
			c2:   GT,
		},
		{
			desc: "Less than",
			c1:   GE,
			c2:   LT,
		},
		{
			desc: "Always",
			c1:   AL,
			c2:   NV,
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			assert.Equal(t, tC.c2, tC.c1.opposite())
			assert.Equal(t, tC.c1, tC.c2.opposite())
		})
	}
}
