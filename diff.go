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
	processFields := make([]map[string]string, len(processOutputs))
	for index, processOutput := range processOutputs {
		processFields[index] = parseOutput(processOutput)
	}

	s := bufio.NewScanner(referenceFile)
	for s.Scan() {
		line := s.Text()

		field, _, ok := parseLine(line)
		key := make([]byte, len(processFields))
		for index, pf := range processFields {
			if !ok {
				key[index] = ' '
				continue
			}
			if _, ok := pf[field]; ok {
				key[index] = fmt.Sprintf("%1d", index+1)[0]
				delete(pf, field)
				continue
			}
			key[index] = '-'
		}
		fmt.Fprintln(output, string(key)+"| "+line)
	}
	if err := s.Err(); err != nil {
		return fmt.Errorf("scanner: %w", err)
	}
	fmt.Fprintln(output)

	for index, pf := range processFields {
		fmt.Fprintf(output, "missed by %d:\n", index+1)
		for field, value := range pf {
			fmt.Fprintf(output, "%s=%s\n", field, value)
		}
		fmt.Fprintln(output)
	}

	return nil
}

func parseOutput(r io.Reader) map[string]string {
	fields := map[string]string{}

	s := bufio.NewScanner(r)
	for s.Scan() {
		if field, value, ok := parseLine(s.Text()); ok {
			fields[field] = value
		}
	}

	return fields
}

func parseLine(line string) (string, string, bool) {
	if index := strings.Index(line, "="); index > -1 {
		field := strings.TrimSpace(line[:index])
		value := line[index+1:]
		return field, value, true
	}
	return "", "", false
}
