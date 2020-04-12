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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Run("WithFile", func(t *testing.T) {
		cases := []struct {
			args     []string
			filename string
			contents string
			expected string
		}{
			{
				filename: "normal",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2\n",
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				filename: "wo_newline",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2\n",
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				filename: "reverse",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2",
				args:     []string{"-r"},
				expected: "0.10.0\n0.2.0\n0.0.2\n0.0.1\n",
			},
			{
				filename: "v-prefix",
				contents: "v0.2.0\nv0.0.1\nv0.10.0\nv0.0.2\n",
				args:     []string{"-p", "v"},
				expected: "v0.0.1\nv0.0.2\nv0.2.0\nv0.10.0\n",
			},
		}

		for _, tt := range cases {
			t.Run(fmt.Sprintf("%q", append(tt.args, tt.filename)), func(t *testing.T) {
				file, err := createTempfile(tt.filename+"-*", tt.contents)
				if err != nil {
					t.Fatalf("cannot create tmpfile %s: %s", tt.filename, err)
				}
				defer os.Remove(file.Name())

				stdout := new(bytes.Buffer)
				stderr := new(bytes.Buffer)
				args := append(tt.args, file.Name())

				if err := Execute(new(bytes.Buffer), stdout, stderr, args); assert.NoError(t, err) {
					assert.Equal(t, tt.expected, stdout.String())
				}
			})
		}
	})

	t.Run("WithMultipleFiles", func(t *testing.T) {
		firstFile, err := createTempfile("first-*", "0.2.0\n0.0.1\n")
		if err != nil {
			t.Fatalf("cannot create tmpfile for first file: %s", err)
		}
		defer os.Remove(firstFile.Name())

		secondFile, err := createTempfile("second-*", "0.10.0\n0.0.2")
		if err != nil {
			t.Fatalf("cannot create tmpfile for second file: %s", err)
		}
		defer os.Remove(secondFile.Name())

		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)
		args := []string{firstFile.Name(), secondFile.Name()}
		if err := Execute(new(bytes.Buffer), stdout, stderr, args); assert.NoError(t, err) {
			expected := "0.0.1\n0.0.2\n0.2.0\n0.10.0\n"
			assert.Equal(t, expected, stdout.String())
		}
	})

	t.Run("WithStdin", func(t *testing.T) {
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
	})
}

func createTempfile(filename, contents string) (*os.File, error) {
	f, err := ioutil.TempFile("", filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if _, err = f.WriteString(contents); err != nil {
		return nil, err
	}

	return f, nil
}
