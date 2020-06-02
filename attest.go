package attest

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/gosuri/uilive"
	"github.com/snsinfu/attest/colors"
)

const caseDelim = "\n---\n"
var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}


type update struct {
	Row  int
	Text string
}

func task(row int, updates chan<- update, argv []string, tcase testCase) {

	end := make(chan bool)

	go func() {
		tick := time.Tick(time.Second / 10)

		for i := 0; ; i++ {
			updates <- update{
				Row:  row,
				Text: fmt.Sprintf(
					"%s  %s",
					colors.Yellow("RUN" + spinner[i % len(spinner)]),
					tcase.Name,
				),
			}

			select {
			case <-tick:
			case <-end:
				return
			}
		}
	}()

	stat, err := test(argv, tcase)
	end <- true
	if err != nil {
		return
	}

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

	updates <- update{
		Row:  row,
		Text: fmt.Sprintf("%s  %s", label, tcase.Name),
	}
}


// Run runs command and tests its output against expected outcomes recorded in
// test files.
func Run(config Config) (int, error) {
	testCases, err := makeTestCases(config)
	if err != nil {
		return 0, err
	}

	sem := make(chan bool, config.MaxJobs)
	updates := make(chan update, len(testCases) * 10)
	wg := sync.WaitGroup{}
	wg.Add(len(testCases))

	for i := range testCases {
		row := i
		tcase := testCases[i]
		go func() {
			updates <- update{
				Row:  row,
				Text: fmt.Sprintf("%s  %s", colors.Gray("WAIT"), tcase.Name),
			}

			sem <- true
			task(row, updates, config.Command, tcase)
			<-sem

			wg.Done()
		}()
	}

	end := make(chan bool)
	endOK := make(chan bool)

	go func() {
		w := uilive.New()
		w.RefreshInterval = 1000*time.Second
		w.Start()

		rows := make([]string, len(testCases))

	render:
		for range time.Tick(100 * time.Millisecond) {
		pump:
			for {
				select {
				case up := <-updates:
					rows[up.Row] = up.Text
				default:
					break pump
				}
			}

			for _, row := range rows {
				fmt.Fprintln(w, row)
			}
			w.Flush()

			select {
			case <-end:
				break render
			default:
			}
		}

		w.Stop()
		endOK <- true
	}()

	wg.Wait()
	end <- true
	<-endOK

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
