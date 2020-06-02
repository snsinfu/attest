package attest

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/snsinfu/attest/colors"
)

const caseDelim = "\n---\n"

// Run runs command and tests its output against expected outcomes recorded in
// test files.
func Run(config Config) (int, error) {
	testCases, err := makeTestCases(config)
	if err != nil {
		return 0, err
	}

	for _, tcase := range testCases {
		stat, err := test(config.Command, tcase)
		if err != nil {
			return 0, err
		}

		var label string
		switch stat {
		case testPassed:
			label = colors.Green("PASS")
		case testFailed:
			label = colors.Red("FAIL")
		case testTimeout:
			label = colors.Yellow("TIME")
		case testError:
			label = colors.Magenta("DEAD")
		}
		fmt.Println(label, tcase.Name)
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
