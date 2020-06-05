package test

import (
	"testing"
	"time"
)

func TestRun_testsSortCommand(t *testing.T) {
	// Just check if test can validate the sort command.
	tc := Case{
		Input:  "xyz\nabc\n123\n",
		Output: "123\nabc\nxyz\n",
	}

	r, err := tc.Run([]string{"sort"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.Outcome != TestPassed {
		t.Errorf("unexpected outcome: %v, expect %v", r.Outcome, TestPassed)
	}

	if r.Stdout != tc.Output {
		t.Errorf("unexpected stdout: %q, expect %q", r.Stdout, tc.Output)
	}

	if r.Stderr != "" {
		t.Errorf("stderr is not empty: %q", r.Stderr)
	}
}

func TestRun_testsTokens(t *testing.T) {
	// Test validates program output not literally but token-wise. That is,
	// extra spaces are ignored. We use the echo command with excessive spaces.
	tc := Case{
		Output: "hello world",
	}

	r, err := tc.Run([]string{"echo", "  hello\n", "\t world \n\n"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.Outcome != TestPassed {
		t.Errorf("unexpected outcome: %v, expect %v", r.Outcome, TestPassed)
	}
}

func TestRun_testsApproxNumber(t *testing.T) {
	// Test validates floating-point numbers to a set number of decilal places.
	// It does not treat non-parsable tokens as numbers.
	tc := Case{
		Output: "abc 123 -4.567890 2x",
		Digits: 4,
	}

	r, err := tc.Run([]string{"echo", "abc", "123.0", "-4.5679", "2x"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.Outcome != TestPassed {
		t.Errorf("unexpected outcome: %v, expect %v", r.Outcome, TestPassed)
	}
}

func TestRun_detectsTimeout(t *testing.T) {
	// Test treats program taking too long time as a failing one.
	tc := Case{
		Timeout: time.Millisecond,
	}

	r, err := tc.Run([]string{"sleep", "5"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.Outcome != TestTimeout {
		t.Errorf("unexpected outcome: %v, expect %v", r.Outcome, TestTimeout)
	}
}

func TestRun_detectsCrash(t *testing.T) {
	// Test detects program returning nonzero exit status. It is a test failure
	// and not an error of the Run() method itself.
	tc := Case{}

	r, err := tc.Run([]string{"false"})
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	if r.Outcome != TestError {
		t.Errorf("unexpected outcome: %v, expect %v", r.Outcome, TestError)
	}
}
