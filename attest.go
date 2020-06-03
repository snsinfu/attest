package attest

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/snsinfu/attest/colors"
	"github.com/snsinfu/attest/flyterm"
	"github.com/snsinfu/attest/periodic"
)

const caseDelim = "\n---\n"
var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

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

			p := periodic.New(time.Second / 10, func() {
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

func formatWait(name string) string {
	return fmt.Sprintf("%s  -:--  %s", colors.Gray("WAIT"), name)
}

func formatRun(name string, elapsed time.Duration, spin int) string {
	min, sec := extractMinSec(elapsed)

	return fmt.Sprintf(
		"%s  %d:%02d  %s",
		colors.Yellow("RUN" + spinner[spin % len(spinner)]),
		min,
		sec,
		name,
	)
}

func formatResult(name string, elapsed time.Duration, stat testStatus) string {
	var label string
	switch stat {
	case testPassed:
		label = colors.Green("PASS")
	case testFailed:
		label = colors.Red("FAIL")
	case testTimeout:
		label = colors.Blue("TIME")
	case testError:
		label = colors.Magenta("DEAD")
	}
	min, sec := extractMinSec(elapsed)

	return fmt.Sprintf("%s  %d:%02d  %s", label, min, sec, name)
}

func extractMinSec(d time.Duration) (int, int) {
	sec := int(d.Seconds())
	min := sec / 60
	sec %= 60
	return min, sec
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
