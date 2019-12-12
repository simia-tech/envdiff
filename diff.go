package envdiff

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/fatih/color"
)

// Diff calculates the diff between the provided reference file and process outputs and writes
// the result to the provided output.
func Diff(noColor bool, output io.Writer, referenceFile io.Reader, processOutputs ...io.Reader) error {
	processFields := make([][]string, len(processOutputs))
	for index, processOutput := range processOutputs {
		processFields[index] = parseFields(processOutput)
	}

	s := bufio.NewScanner(referenceFile)
	for s.Scan() {
		line := s.Text()

		field, ok := extractField(line)
		key := make([]byte, len(processFields))
		for index, pf := range processFields {
			if !ok {
				key[index] = ' '
				continue
			}
			if contains(pf, field) {
				key[index] = fmt.Sprintf("%1d", index+1)[0]
				continue
			}
			key[index] = '-'
		}
		if noColor {
			fmt.Fprintln(output, string(key)+"| "+line)
		} else {
			selectColor(key).Fprintln(output, line)
		}
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

func selectColor(key []byte) *color.Color {
	c := color.New(color.FgWhite)
	for _, k := range key {
		switch k {
		case '-':
			c = c.Add(color.FgBlue)
		}
	}
	return c
}
