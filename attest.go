package attest

import (
	"fmt"
	"time"

	"github.com/snsinfu/attest/flyterm"
	"github.com/snsinfu/attest/periodic"
	"github.com/snsinfu/attest/test"
)

const spinnerInterval = time.Second / 10

// A Config specifies how tests should be run.
type Config struct {
	Command   []string
	TestCases []test.Case
	MaxJobs   int
	Verbose   bool
}

// Run runs command and tests its output against expected outcomes recorded in
// test files.
func Run(config Config) (int, error) {
	testCount := len(config.TestCases)
	term := flyterm.New(testCount, flyterm.Options{})

	type update struct {
		Index  int
		Result test.Result
	}
	updates := make(chan update, testCount)
	sem := make(chan bool, config.MaxJobs)

	for i := range config.TestCases {
		// Loop variable needs to be assigned to a local variable so that it
		// is correctly captured by the closure below.
		row := i
		tc := config.TestCases[i]

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
			r, _ := tc.Run(config.Command)
			p.Stop()

			elapsed := time.Now().Sub(start)
			term.Update(row, formatOutcome(tc.Name, elapsed, r.Outcome))

			<-sem
			updates <- update{row, r}
		}()
	}

	// This loop blocks until all the tests finish.
	failed := false
	results := make([]test.Result, testCount)
	for i := 0; i < testCount; i++ {
		up := <-updates
		results[up.Index] = up.Result
		if up.Result.Outcome != test.TestPassed {
			failed = true
		}
	}

	term.Stop()

	if config.Verbose {
		for i := 0; i < testCount; i++ {
			tc := config.TestCases[i]
			r := results[i]
			if r.Outcome == test.TestPassed {
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
