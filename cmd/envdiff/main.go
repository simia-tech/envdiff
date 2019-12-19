package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/simia-tech/envdiff"
)

func main() {
	flag.Usage = func() {
		w := flag.CommandLine.Output()
		fmt.Fprintf(w, "Usage: %s [reference file] [binary path 1] [binary path 2] ...\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()
	if len(args) < 2 {
		fmt.Printf("expected a reference file and at least one binary path as command line arguments\n")
		return
	}

	if err := calculateDiff(args[0], args[1:]...); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func calculateDiff(referencePath string, processPaths ...string) error {
	referenceFile, err := os.Open(referencePath)
	if err != nil {
		return fmt.Errorf("open [%s]: %w", referencePath, err)
	}
	defer referenceFile.Close()

	processOutputs := []io.Reader{}
	for _, path := range processPaths {
		cmd := exec.Command(path, "-print-env", "-print-env-format", "short-bash")
		out, _ := cmd.CombinedOutput()
		processOutputs = append(processOutputs, bytes.NewReader(out))
	}

	if err := envdiff.Diff(os.Stdout, referenceFile, processOutputs...); err != nil {
		return fmt.Errorf("diff: %w", err)
	}

	return nil
}
