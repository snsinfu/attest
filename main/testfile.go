package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"

	"github.com/snsinfu/attest/test"
)

const caseDelim = "---\n"

func makeTestCases(files []string) ([]test.Case, error) {
	testCases := []test.Case{}

	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}

		input, output, ok := parseTestCase(string(data))
		if !ok {
			return nil, fmt.Errorf("test file %s is not formatted correctly", file)
		}

		testCases = append(testCases, test.Case{
			Name:   path.Base(file),
			Input:  input,
			Output: output,
		})
	}
	return testCases, nil
}

func parseTestCase(s string) (string, string, bool) {
	pos := strings.Index(s, caseDelim)
	if pos == -1 || (pos > 0 && s[pos-1] != '\n') {
		return "", "", false
	}
	input := s[:pos]
	output := s[pos+len(caseDelim):]
	return input, output, true
}
