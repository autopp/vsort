// Copyright (C) 2020 Akira Tanimura (@autopp)
//
// Licensed under the Apache License, Version 2.0 (the “License”);
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an “AS IS” BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		{
			input: `0.2.0
0.0.1
0.10.0
0.0.2
`,
			expected: `0.0.1
0.0.2
0.2.0
0.10.0
`,
		},
		{
			input: `0.2.0
0.0.1
0.10.0
0.0.2`,
			args: []string{"-r"},
			expected: `0.10.0
0.2.0
0.0.2
0.0.1
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
