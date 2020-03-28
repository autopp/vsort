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
		{v1: "0.1.1", v2: "0.1.0"},
	}

	for _, tt := range cases {
		assert.Equal(t, tt.expected, Compare(tt.v1, tt.v2))
	}
}
