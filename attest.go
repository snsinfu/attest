package attest

import (
	"fmt"
	"time"

	"github.com/snsinfu/attest/flyterm"
	"github.com/snsinfu/attest/periodic"
)

const spinnerInterval = time.Second / 10

// A Config specifies how tests should be run.
type Config struct {
	Command []string
	Tests   []string
	Digits  int
	MaxJobs int
	Timeout time.Duration
	Verbose bool
}

// Run runs command and tests its output against expected outcomes recorded in
// test files.
func Run(config Config) (int, error) {
	testCases, err := makeTestCases(config)
	if err != nil {
		return 0, err
	}

	sem := make(chan bool, config.MaxJobs)

	term := flyterm.New(len(testCases), flyterm.Options{})

	type update struct {
		Index  int
		Result testResult
	}
	updates := make(chan update, len(testCases))

	for i := range testCases {
		row := i
		tc := testCases[i]

		go func() {
			term.Update(row, formatWait(tc.Name))

			sem <- true

			start := time.Now()
			spin := 0

			p := periodic.New(spinnerInterval, func() {
				elapsed := time.Now().Sub(start)
				term.Update(row, formatRun(tc.Name, elapsed, spin))
				spin++
			})
			r, _ := test(config.Command, tc)
			p.Stop()

			elapsed := time.Now().Sub(start)
			term.Update(row, formatOutcome(tc.Name, elapsed, r.Outcome))

			<-sem
			updates <- update{row, r}
		}()
	}

	// This loop blocks until all the tests finish.
	failed := false
	results := make([]testResult, len(testCases))
	for i := 0; i < len(testCases); i++ {
		up := <-updates
		results[up.Index] = up.Result
		if up.Result.Outcome != testPassed {
			failed = true
		}
	}

	term.Stop()

	if config.Verbose {
		for i := 0; i < len(testCases); i++ {
			tc := testCases[i]
			r := results[i]
			if r.Outcome == testPassed {
				continue
			}
			fmt.Print("\n")
			fmt.Print(formatResult(tc, r))
		}
	}

	if failed {
		return 1, nil
	}
	return 0, nil
}
