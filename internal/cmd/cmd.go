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

			order := vsort.Asc
			if reverse {
				order = vsort.Desc
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
