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

type testStatus int

const (
	testPassed testStatus = iota
	testFailed
	testTimeout
	testError
)

const caseDelim = "\n---\n"

// test runs a command and test its behavior against given test case.
func test(argv []string, tcase testCase) (testStatus, error) {
	cmd, err := command.Run(argv)
	if err != nil {
		return testError, err
	}
	defer cmd.Close()

	cmd.Write([]byte(tcase.Input))
	cmd.WriteEnd()

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

	stdout, err := ioutil.ReadAll(cmd)
	if err != nil {
		return testError, err
	}

	// Check exit status. Note that failed command is not our fault. It is a
	// valid observation of testError.
	if err := cmd.Wait(); err != nil {
		select {
		case <-timedOut:
			return testTimeout, nil
		default:
			return testError, nil
		}
	}

	// Test command output against expected one. We use token-wise comparison.
	// TODO: Correctly compare floating-point numbers.
	observed := strings.Fields(string(stdout))
	expected := strings.Fields(tcase.Output)

	if len(observed) != len(expected) {
		return testFailed, nil
	}

	for i := 0; i < len(observed); i++ {
		if observed[i] != expected[i] {
			return testFailed, nil
		}
	}

	return 0, nil
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
