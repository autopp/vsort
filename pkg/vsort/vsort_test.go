package vsort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
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
		actual, err := Compare(tt.v1, tt.v2)
		if assert.NoError(t, err) {
			assert.Equal(t, tt.expected, actual)
		}
	}
}
