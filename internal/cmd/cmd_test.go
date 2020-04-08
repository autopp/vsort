package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		input    string
		args     []string
		expected string
	}{
		{
			input: `0.2.0
0.0.1
0.10.0
0.0.2`,
			expected: `0.0.1
0.0.2
0.2.0
0.10.0
`,
		},
	}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%q", tt.args), func(t *testing.T) {
			stdin := bytes.NewBufferString(tt.input)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			if err := Execute(stdin, stdout, stderr, tt.args); assert.NoError(t, err) {
				assert.Equal(t, tt.expected, stdout.String())
			}
		})
	}
}
