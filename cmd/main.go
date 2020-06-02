package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/snsinfu/attest"
)

const usage = `
usage: attest [options] <command>...

options:
  -d <tests>    Directory containing test files [default: tests]
  -j <jobs>     Number of concurrent runs; 0 means maximum [default: 0]
  -t <timeout>  Timeout in seconds; 0 means no timeout [default: 0]
  -h            Show this message and exit

attest loads test files (*.txt) from test directory and examines command
behavior against input and output text written in each test file. The test
file must be formatted like this:

  input
  ---
  output

Namely, input lines and output lines are delimited by a line consisting of
three hyphens.
`

const testGlob = "*.txt"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run() error {
	parser := docopt.Parser{
		HelpHandler:  docopt.PrintHelpAndExit,
		OptionsFirst: true,
	}
	opts, err := parser.ParseArgs(usage, nil, "")
	if err != nil {
		return err
	}

	testsDir, err := opts.String("-d")
	if err != nil {
		return err
	}
	tests, err := filepath.Glob(filepath.Join(testsDir, testGlob))
	if err != nil {
		return err
	}

	maxJobs, err := opts.Int("-j")
	if err != nil {
		return err
	}
	if maxJobs == 0 {
		maxJobs = runtime.NumCPU()
	}
	if maxJobs <= 0 {
		return fmt.Errorf("concurrency (-j) cannot be negative")
	}

	timeoutSec, err := opts.Int("-t")
	if err != nil {
		return err
	}
	if timeoutSec < 0 {
		return fmt.Errorf("timeout (-t) cannot be negative")
	}

	config := attest.Config{
		Command: opts["<command>"].([]string),
		Tests:   tests,
		MaxJobs: maxJobs,
		Timeout: time.Duration(timeoutSec) * time.Second,
	}

	rc, err := attest.Run(config)
	if err != nil {
		return err
	}
	os.Exit(rc)

	panic("cannot reach")
}
