package attest

import (
	"sync"
	"time"

	"github.com/snsinfu/attest/flyterm"
	"github.com/snsinfu/attest/periodic"
)

const spinnerInterval = time.Second / 10

// A Config specifies how tests should be run.
type Config struct {
	Command []string
	Tests   []string
	MaxJobs int
	Timeout time.Duration
}

// Run runs command and tests its output against expected outcomes recorded in
// test files.
func Run(config Config) (int, error) {
	testCases, err := makeTestCases(config)
	if err != nil {
		return 0, err
	}

	sem := make(chan bool, config.MaxJobs)
	wg := sync.WaitGroup{}
	wg.Add(len(testCases))

	term := flyterm.New(len(testCases), flyterm.Options{})
	defer term.Stop()

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
			stat, _ := test(config.Command, tc)
			p.Stop()

			elapsed := time.Now().Sub(start)
			term.Update(row, formatResult(tc.Name, elapsed, stat))

			<-sem
			wg.Done()
		}()
	}

	wg.Wait()

	return 0, nil
}
