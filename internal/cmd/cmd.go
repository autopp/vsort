package cmd

import (
	"bufio"
	"fmt"
	"io"

	"github.com/autopp/vsort/pkg/vsort"
)

// Execute execute main logic
func Execute(stdin io.Reader, stdout, stderr io.Writer, args []string) error {
	scanner := bufio.NewScanner(stdin)
	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(stderr, err)
		return err
	}

	vsort.Sort(lines)
	for _, line := range lines {
		fmt.Fprintln(stdout, line)
	}

	return nil
}
