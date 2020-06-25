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
usage: attest [-v ...] [options] <command>...

options:
  -d <tests>    Directory containing test files [default: tests]
  -f <digits>   Test numbers for specified number of decimal places
  -j <jobs>     Number of concurrent runs; 0 means maximum [default: 0]
  -t <timeout>  Timeout in seconds; 0 means no timeout [default: 0]
  -v            Display detailed test results; -v for only failed tests and -vv
                for all tests
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

	// Load test files.
	testsDir, err := opts.String("-d")
	if err != nil {
		return err
	}
	testFiles, err := filepath.Glob(filepath.Join(testsDir, testGlob))
	if err != nil {
		return err
	}
	testCases, err := makeTestCases(testFiles)
	if err != nil {
		return err
	}

	// Set rounded floating-point comparison mode.
	digits := 0
	if opt, ok := opts["-f"]; ok && opt != nil {
		digits, err = opts.Int("-f")
		if err != nil {
			return err
		}
	}

	// Configure concurrency.
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

	// Set timeout of single test run.
	timeoutSec, err := opts.Int("-t")
	if err != nil {
		return err
	}
	if timeoutSec < 0 {
		return fmt.Errorf("timeout (-t) cannot be negative")
	}

	// Get verbosity as the number of -v flags passed. docopt-go does not allow
	// calling Int("-v") on counted flags so we resort to manual cast here.
	verbose, ok := opts["-v"].(int)
	if !ok {
		return fmt.Errorf("failed to get -v option")
	}

	// Cook test cases.
	for i := range testCases {
		testCases[i].Timeout = time.Duration(timeoutSec) * time.Second
		testCases[i].Digits = digits
	}

	config := attest.Config{
		Command:   opts["<command>"].([]string),
		TestCases: testCases,
		MaxJobs:   maxJobs,
		Verbose:   verbose,
	}

	rc, err := attest.Run(config)
	if err != nil {
		return err
	}
	os.Exit(rc)

	panic("cannot reach")
}
