package attest

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/snsinfu/attest/test"
)

const caseDelim = "\n---\n"

func makeTestCases(c Config) ([]test.Case, error) {
	testCases := []test.Case{}

	for _, filename := range c.Tests {
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

		testCases = append(testCases, test.Case{
			Name:    path.Base(filename),
			Input:   text[:inputEnd],
			Output:  text[outputStart:],
			Timeout: c.Timeout,
			Digits:  c.Digits,
		})
	}
	return testCases, nil
}
