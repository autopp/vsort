package cmd

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	cases := []struct {
		stdin    string
		args     []string
		expected string
	}{}

	for _, tt := range cases {
		t.Run(fmt.Sprintf("%q", tt.args), func(t *testing.T) {
			stdin := bytes.NewBufferString(tt.stdin)
			stdout := new(bytes.Buffer)
			stderr := new(bytes.Buffer)

			if err := Execute(stdin, stdout, stderr, tt.args); assert.NoError(t, err) {
				assert.Equal(t, tt.expected, stdout.String())
			}
		})
	}
}
