package test

import (
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/snsinfu/attest/command"
)

type Case struct {
	Name    string
	Input   string
	Output  string
	Timeout time.Duration
	Digits  int
}

type Outcome int

const (
	TestPassed Outcome = iota
	TestFailed
	TestTimeout
	TestError
)

type Result struct {
	Outcome Outcome
	Stdout  string
	Stderr  string
}

// Run tests command specified by argv. Returns a Result object containing test
// outcome and observed output of the command. The outcome is TestError on any
// error.
func (c *Case) Run(argv []string) (Result, error) {
	r := Result{Outcome: TestError}

	cmd, err := command.Run(argv)
	if err != nil {
		return r, err
	}

	cmd.Stdin.Write([]byte(c.Input))
	cmd.Stdin.Close()

	// We have feeded the input to the command. The command process should be
	// computing result now and reading stdout would block. So, handle timeout
	// here. If the command does not output result within a set time, we
	// forcifully terminate the command process and flag timedOut.
	timedOut := make(chan bool)

	var timer *time.Timer
	if c.Timeout != 0 {
		timer = time.AfterFunc(c.Timeout, func() {
			timer.Stop()
			cmd.Signal(syscall.SIGTERM)
			timedOut <- true
		})
		defer timer.Stop()
	}

	stdout, err := ioutil.ReadAll(cmd.Stdout)
	if err != nil {
		return r, err
	}

	stderr, err := ioutil.ReadAll(cmd.Stderr)
	if err != nil {
		return r, err
	}

	r.Stdout = string(stdout)
	r.Stderr = string(stderr)

	// Check exit status. Note that failed command is not our fault. It is a
	// valid observation of testError.
	if err := cmd.Wait(); err != nil {
		select {
		case <-timedOut:
			r.Outcome = TestTimeout
			return r, nil
		default:
			r.Outcome = TestError
			return r, nil
		}
	}

	// Test command output against expected one. We use token-wise comparison
	// so that spacing differences do not affect the validity of the output.
	observed := strings.Fields(string(stdout))
	expected := strings.Fields(c.Output)

	if len(observed) != len(expected) {
		r.Outcome = TestFailed
		return r, nil
	}

	for i := 0; i < len(observed); i++ {
		if !c.match(observed[i], expected[i]) {
			r.Outcome = TestFailed
			return r, nil
		}
	}

	r.Outcome = TestPassed
	return r, nil
}

// match tests observed token obs against expected one exp for equivalence.
func (c *Case) match(obs, exp string) bool {
	// Round to a set number of decimal places if the token looks like a number.
	if c.Digits != 0 {
		obsNum, err1 := strconv.ParseFloat(obs, 64)
		expNum, err2 := strconv.ParseFloat(exp, 64)
		if err1 == nil && err2 == nil {
			obs = fmt.Sprintf("%.*f", c.Digits, obsNum)
			exp = fmt.Sprintf("%.*f", c.Digits, expNum)
		}
	}
	return obs == exp
}
