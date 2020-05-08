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
	"errors"
	"io"
	"os"

	"encoding/json"
	"io/ioutil"

	"fmt"

	"github.com/autopp/vsort/pkg/vsort"
	"github.com/spf13/cobra"
)

// Execute execute main logic
func Execute(version string, stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	const (
		linesInput = "lines"
		jsonInput  = "json"
	)

	const (
		linesOutput = "lines"
		jsonOutput  = "json"
	)

	cmd := &cobra.Command{
		Use:           "vsort",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Get --version and process if given
			if showVersion, err := cmd.Flags().GetBool("version"); err != nil {
				return err
			} else if showVersion {
				cmd.Printf("vsort version %s\n", version)
				return nil
			}

			// Get --input
			input, err := cmd.Flags().GetString("input")
			if err != nil {
				return err
			}

			inputFunc, ok := map[string]func(io.Reader) ([]string, error){
				linesInput: readLines, jsonInput: readJSON,
			}[input]
			if !ok {
				return fmt.Errorf("unknown input format: %q (expected %q or %q)", input, linesInput, jsonInput)
			}

			// Get --output
			output, err := cmd.Flags().GetString("output")
			if err != nil {
				return err
			}

			outputFunc, ok := map[string]func(*cobra.Command, []string) error{
				linesOutput: outputLines, jsonOutput: outputJSON,
			}[output]
			if !ok {
				return fmt.Errorf("unknown output format: %q (expected %q or %q)", output, linesOutput, jsonOutput)
			}

			// Get --reverse
			reverse, err := cmd.Flags().GetBool("reverse")
			if err != nil {
				return err
			}

			// Get --prefix
			prefix, err := cmd.Flags().GetString("prefix")
			if err != nil {
				return err
			}

			// Get --suffix
			suffix, err := cmd.Flags().GetString("suffix")
			if err != nil {
				return err
			}

			// Get --level
			level, err := cmd.Flags().GetInt("level")
			if err != nil {
				return err
			}

			// Get --strict
			strict, err := cmd.Flags().GetBool("strict")
			if err != nil {
				return err
			}

			type inputStream struct {
				name string
				r    io.Reader
			}
			var is []inputStream

			if len(args) == 0 {
				is = []inputStream{{"<stdin>", cmd.InOrStdin()}}
			} else {
				is = make([]inputStream, len(args))

				for i, path := range args {
					f, err := os.Open(path)
					if err != nil {
						return err
					}
					defer f.Close()
					is[i].name = path
					is[i].r = f
				}
			}

			var versions []string
			for _, i := range is {
				vs, err := inputFunc(i.r)
				if err != nil {
					return fmt.Errorf("cannot read from %s: %w", i.name, err)
				}
				versions = append(versions, vs...)
			}

			order := vsort.WithOrder(vsort.Asc)
			if reverse {
				order = vsort.WithOrder(vsort.Desc)
			}

			options := []vsort.Option{order, vsort.WithPrefix(prefix), vsort.WithLevel(level)}
			if suffix != "" {
				options = append(options, vsort.WithSuffix(suffix))
			}
			s, err := vsort.NewSorter(options...)
			if err != nil {
				return err
			}

			// validate inputs
			validated := make([]string, 0, len(versions))
			for _, v := range versions {
				if s.IsValid(v) {
					validated = append(validated, v)
				} else if strict {
					msg := fmt.Sprintf("invalid version is contained: %s\n", v)
					cmd.PrintErrln(msg)
					return errors.New(msg)
				}
			}

			s.Sort(validated)

			if err := outputFunc(cmd, validated); err != nil {
				return err
			}

			return nil
		},
	}

	cmd.Flags().BoolP("version", "v", false, "Print the version and silently exits.")
	cmd.Flags().StringP("input", "i", linesInput, `Specify input format. Accepted values are "lines" or "json" (default: "lines").`)
	cmd.Flags().StringP("output", "o", linesOutput, `Specify output format. Accepted values are "lines" or "json" (default: "lines").`)
	cmd.Flags().BoolP("reverse", "r", false, "Sort in reverse order.")
	cmd.Flags().StringP("prefix", "p", "", "Expected prefix of version string.")
	cmd.Flags().StringP("suffix", "s", "", "Expected suffix pattern of version string.")
	cmd.Flags().IntP("level", "L", -1, "Expected version level")
	cmd.Flags().Bool("strict", false, "Make error when invalid version is contained.")

	cmd.SetIn(stdin)
	cmd.SetOut(stdout)
	cmd.SetErr(stderr)
	cmd.SetArgs(args)
	return cmd.Execute()
}

func readLines(r io.Reader) ([]string, error) {
	lines := make([]string, 0)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func readJSON(r io.Reader) ([]string, error) {
	j, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	var a []string
	if err := json.Unmarshal(j, &a); err != nil {
		return nil, err
	}
	return a, nil
}

func outputLines(cmd *cobra.Command, versions []string) error {
	for _, v := range versions {
		cmd.Println(v)
	}

	return nil
}

func outputJSON(cmd *cobra.Command, versions []string) error {
	b, err := json.Marshal(versions)
	if err != nil {
		return err
	}
	cmd.Print(string(b))

	return nil
}
