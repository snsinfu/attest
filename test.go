package attest

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/snsinfu/attest/command"
)

type testCase struct {
	Name    string
	Input   string
	Output  string
	Timeout time.Duration
}

const caseDelim = "\n---\n"

type testOutcome int

const (
	testPassed testOutcome = iota
	testFailed
	testTimeout
	testError
)

type testResult struct {
	Outcome testOutcome
	Stdout  string
	Stderr  string
}

// test runs a command and test its behavior against given test case.
func test(argv []string, tcase testCase) (testResult, error) {
	var r testResult

	cmd, err := command.Run(argv)
	if err != nil {
		r.Outcome = testError
		return r, err
	}

	cmd.Stdin.Write([]byte(tcase.Input))
	cmd.Stdin.Close()

	// We have feeded the input to the command. The command process should be
	// computing result now and reading stdout would block. So, handle timeout
	// here. If the command does not output result within a set time, we
	// forcifully terminate the command process and flag timedOut.
	timedOut := make(chan bool)

	var timer *time.Timer
	if tcase.Timeout != 0 {
		timer = time.AfterFunc(tcase.Timeout, func() {
			timer.Stop()
			cmd.Signal(syscall.SIGTERM)
			timedOut <- true
		})
		defer timer.Stop()
	}

	stdout, err := ioutil.ReadAll(cmd.Stdout)
	if err != nil {
		r.Outcome = testError
		return r, err
	}

	stderr, err := ioutil.ReadAll(cmd.Stderr)
	if err != nil {
		r.Outcome = testError
		return r, err
	}

	r.Stdout = string(stdout)
	r.Stderr = string(stderr)

	// Check exit status. Note that failed command is not our fault. It is a
	// valid observation of testError.
	if err := cmd.Wait(); err != nil {
		select {
		case <-timedOut:
			r.Outcome = testTimeout
			return r, nil
		default:
			r.Outcome = testError
			return r, nil
		}
	}

	// Test command output against expected one. We use token-wise comparison
	// so that spacing differences do not affect the validity of the output.
	// TODO: Correctly compare floating-point numbers.
	observed := strings.Fields(string(stdout))
	expected := strings.Fields(tcase.Output)

	if len(observed) != len(expected) {
		r.Outcome = testFailed
		return r, nil
	}

	for i := 0; i < len(observed); i++ {
		if observed[i] != expected[i] {
			r.Outcome = testFailed
			return r, nil
		}
	}

	r.Outcome = testPassed
	return r, nil
}

// makeTestCases interprets config and assembles test case objects.
func makeTestCases(config Config) ([]testCase, error) {
	testCases := []testCase{}

	for _, filename := range config.Tests {
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, err
		}

		text := string(data)
		pos := strings.Index(text, caseDelim)
		if pos == -1 {
			return nil, fmt.Errorf(
				"test file %s does not contain input and output delimited by ---", filename,
			)
		}
		inputEnd := pos + 1
		outputStart := pos + len(caseDelim)

		testCases = append(testCases, testCase{
			Name:    path.Base(filename),
			Input:   text[:inputEnd],
			Output:  text[outputStart:],
			Timeout: config.Timeout,
		})
	}

	return testCases, nil
}
