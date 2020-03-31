package vsort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComparatorCompare(t *testing.T) {
	cases := []struct {
		v1       string
		v2       string
		expected int
	}{
		{v1: "0.1.0", v2: "0.1.0", expected: 0},
		{v1: "0.1.1", v2: "0.1.0", expected: 1},
		{v1: "0.1.0", v2: "0.1.1", expected: -1},
	}

	for _, tt := range cases {
		var c Comparator
		actual, err := c.Compare(tt.v1, tt.v2)
		if assert.NoError(t, err) {
			assert.Equal(t, tt.expected, actual)
		}
	}
}
