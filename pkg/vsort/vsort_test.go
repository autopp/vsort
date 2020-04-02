package vsort

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparatorCompare(t *testing.T) {
	type Case struct {
		v1       string
		v2       string
		expected int
	}
	cases := []Case{
		{v1: "0.1.0", v2: "0.1.0", expected: 0},
		{v1: "0.1.1", v2: "0.1.0", expected: 1},
		{v1: "0.1.0", v2: "0.1.1", expected: -1},
	}

	genSubtestName := func(c Case) string {
		var op string
		switch c.expected {
		case 0:
			op = "="
		case 1:
			op = ">"
		case -1:
			op = "<"
		}
		return fmt.Sprintf("%q%s%q", c.v1, op, c.v2)
	}

	for _, tt := range cases {
		t.Run(genSubtestName(tt), func(t *testing.T) {
			var c Comparator
			actual, err := c.Compare(tt.v1, tt.v2)
			if assert.NoError(t, err) {
				assert.Equal(t, tt.expected, actual)
			}
		})
	}
}
