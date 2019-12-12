package envdiff

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

// Diff calculates the diff between the provided reference file and process outputs and writes
// the result to the provided output.
func Diff(output io.Writer, referenceFile io.Reader, processOutputs ...io.Reader) error {
	processFields := make([][]string, len(processOutputs))
	for index, processOutput := range processOutputs {
		processFields[index] = parseFields(processOutput)
	}

	s := bufio.NewScanner(referenceFile)
	for s.Scan() {
		line := s.Text()

		field, ok := extractField(line)
		for index, pf := range processFields {
			if !ok {
				fmt.Fprint(output, " ")
				continue
			}
			if contains(pf, field) {
				fmt.Fprintf(output, "%1d", index+1)
				continue
			}
			fmt.Fprint(output, "-")
		}
		fmt.Fprint(output, "| ")

		fmt.Fprintln(output, line)
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}

	return nil
}

func parseFields(r io.Reader) []string {
	fields := []string{}

	s := bufio.NewScanner(r)
	for s.Scan() {
		if field, ok := extractField(s.Text()); ok {
			fields = append(fields, field)
		}
	}

	return fields
}

func extractField(line string) (string, bool) {
	if index := strings.Index(line, "="); index > -1 {
		return strings.TrimSpace(line[:index]), true
	}
	return "", false
}

func contains(items []string, item string) bool {
	for _, i := range items {
		if i == item {
			return true
		}
	}
	return false
}
