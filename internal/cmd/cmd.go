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
	"fmt"
	"io"

	"github.com/autopp/vsort/pkg/vsort"
	"github.com/spf13/cobra"
)

// Execute execute main logic
func Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	cmd := &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			reverse, err := cmd.Flags().GetBool("reverse")
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(stdin)
			lines := make([]string, 0)

			for scanner.Scan() {
				lines = append(lines, scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				fmt.Fprintln(stderr, err)
				return err
			}

			order := vsort.WithOrder(vsort.Asc)
			if reverse {
				order = vsort.WithOrder(vsort.Desc)
			}
			vsort.Sort(lines, order)
			for _, line := range lines {
				fmt.Fprintln(stdout, line)
			}

			return nil
		},
	}

	cmd.Flags().BoolP("reverse", "r", false, "Sort in reverse order.")
	cmd.SetArgs(args)
	return cmd.Execute()
}
