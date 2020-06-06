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
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	t.Run("WithVersionFlag", func(t *testing.T) {
		version := "v0.1.0"
		stdout := new(bytes.Buffer)
		stderr := new(bytes.Buffer)

		if assertSuccessWithNoStderr(t, version, new(bytes.Buffer), stdout, stderr, []string{"--version"}) {
			assert.Contains(t, stdout.String(), "vsort version "+version)
		}
	})

	t.Run("WithFile", func(t *testing.T) {
		cases := []struct {
			args     []string
			filename string
			contents string
			success  bool
			expected string
		}{
			{
				filename: "normal",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2\n",
				success:  true,
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				filename: "wo_newline",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2\n",
				success:  true,
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				filename: "reverse",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2",
				args:     []string{"-r"},
				success:  true,
				expected: "0.10.0\n0.2.0\n0.0.2\n0.0.1\n",
			},
			{
				filename: "v-prefix",
				contents: "v0.2.0\nv0.0.1\nv0.10.0\nv0.0.2\n",
				args:     []string{"-p", "v"},
				success:  true,
				expected: "v0.0.1\nv0.0.2\nv0.2.0\nv0.10.0\n",
			},
			{
				filename: "alphas-prefix",
				contents: "version-0.2.0\nv-0.0.1\nversion-0.10.0\nv-0.0.2\n",
				args:     []string{"-p", "[a-z]+-"},
				success:  true,
				expected: "v-0.0.1\nv-0.0.2\nversion-0.2.0\nversion-0.10.0\n",
			},
			{
				filename: "release-suffix",
				contents: "0.2.0-1\n0.0.1-2\n0.10.0-3\n0.0.2-4\n",
				args:     []string{"-s", `-\d+`},
				success:  true,
				expected: "0.0.1-2\n0.0.2-4\n0.2.0-1\n0.10.0-3\n",
			},
			{
				filename: "json-input",
				contents: `["0.2.0", "0.0.1", "0.10.0", "0.0.2"]`,
				args:     []string{"-i", "json"},
				success:  true,
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				filename: "json-output",
				contents: "0.2.0\n0.0.1\n0.10.0\n0.0.2\n",
				args:     []string{"-o", "json"},
				success:  true,
				expected: `["0.0.1","0.0.2","0.2.0","0.10.0"]`,
			},
			{
				filename: "level2",
				contents: "2.0\n0.1\n10.0\n0.2\n",
				args:     []string{"-L", "2"},
				success:  true,
				expected: "0.1\n0.2\n2.0\n10.0\n",
			},
			{
				filename: "with-invalid",
				contents: "0.2.0\nv0.3.0\n0.0.1\n0.10.0\n1.0.0-a\n0.0.2\n",
				success:  true,
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				filename: "with-invalid-strict",
				contents: "0.2.0\nv0.3.0\n0.0.1\n0.10.0\n1.0.0-a\n0.0.2\n",
				args:     []string{"--strict"},
				success:  false,
			},
		}

		for _, tt := range cases {
			t.Run(fmt.Sprintf("%q", append(tt.args, tt.filename)), func(t *testing.T) {
				file, err := createTempfile(tt.filename+"-*", tt.contents)
				if err != nil {
					t.Fatalf("cannot create tmpfile %s: %s", tt.filename, err)
				}
				defer os.Remove(file.Name())

				version := "HEAD"
				stdin := new(bytes.Buffer)
				stdout := new(bytes.Buffer)
				stderr := new(bytes.Buffer)
				args := append(tt.args, file.Name())

				if tt.success {
					if assertSuccessWithNoStderr(t, version, stdin, stdout, new(bytes.Buffer), args) {
						assert.Equal(t, tt.expected, stdout.String())
					}
				} else {
					assert.Error(t, Execute(version, stdin, stdout, stderr, args))
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
		args := []string{firstFile.Name(), secondFile.Name()}

		if assertSuccessWithNoStderr(t, "HEAD", new(bytes.Buffer), stdout, new(bytes.Buffer), args) {
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
				input:    "0.2.0\n0.0.1\n0.10.0\n0.0.2",
				expected: "0.0.1\n0.0.2\n0.2.0\n0.10.0\n",
			},
			{
				input:    "0.2.0\n0.0.1\n0.10.0\n0.0.2",
				args:     []string{"-r"},
				expected: "0.10.0\n0.2.0\n0.0.2\n0.0.1\n",
			},
		}

		for _, tt := range cases {
			t.Run(fmt.Sprintf("%q", tt.args), func(t *testing.T) {
				stdin := bytes.NewBufferString(tt.input)
				stdout := new(bytes.Buffer)

				if assertSuccessWithNoStderr(t, "HEAD", stdin, stdout, new(bytes.Buffer), tt.args) {
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

func assertSuccessWithNoStderr(t *testing.T, version string, stdin io.Reader, stdout, stderr *bytes.Buffer, args []string) bool {
	return assert.NoError(t, Execute(version, stdin, stdout, stderr, args)) && assert.Empty(t, stderr.String())
}
