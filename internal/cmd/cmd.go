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
	"bufio"
	"io"
	"os"

	"encoding/json"
	"io/ioutil"

	"github.com/autopp/vsort/pkg/vsort"
	"github.com/spf13/cobra"
)

// Execute execute main logic
func Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			jsonInput, err := cmd.Flags().GetBool("json-input")
			if err != nil {
				return err
			}

			reverse, err := cmd.Flags().GetBool("reverse")
			if err != nil {
				return err
			}

			prefix, err := cmd.Flags().GetString("prefix")
			if err != nil {
				return err
			}

			var rs []io.Reader
			if len(args) == 0 {
				rs = []io.Reader{cmd.InOrStdin()}
			} else {
				rs = make([]io.Reader, len(args))
				for i, path := range args {
					f, err := os.Open(path)
					if err != nil {
						return err
					}
					defer f.Close()
					rs[i] = f
				}
			}

			var lines []string
			if jsonInput {
				lines, err = readJSONFromStreams(rs)
			} else {
				lines, err = readLinesFromStreams(rs)
			}
			if err != nil {
				return err
			}

			order := vsort.WithOrder(vsort.Asc)
			if reverse {
				order = vsort.WithOrder(vsort.Desc)
			}
			s := vsort.NewSorter(order, vsort.WithPrefix(prefix))
			s.Sort(lines)
			for _, line := range lines {
				cmd.Println(line)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("json-input", "j", false, "assume that input is json array")
	cmd.Flags().BoolP("reverse", "r", false, "Sort in reverse order.")
	cmd.Flags().StringP("prefix", "p", "", "Expected prefix of version string")

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func readJSONFromStreams(rs []io.Reader) ([]string, error) {
	lines := make([]string, 0)
	for _, r := range rs {
		j, err := ioutil.ReadAll(r)
		if err != nil {
			return nil, err
		}

		var a []string
		if err := json.Unmarshal(j, &a); err != nil {
			return nil, err
		}

		lines = append(lines, a...)
	}
	return lines, nil
}

func readLinesFromStreams(rs []io.Reader) ([]string, error) {
	lines := make([]string, 0)
	for _, r := range rs {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	return lines, nil
}
