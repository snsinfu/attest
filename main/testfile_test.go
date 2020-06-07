package main

import "testing"

func TestParseTestCase_rejectsBadDocument(t *testing.T) {
	if _, _, ok := parseTestCase(""); ok {
		t.Errorf("incorrectly accepts empty document")
	}

	if _, _, ok := parseTestCase("\n"); ok {
		t.Errorf("incorrectly accepts blank document")
	}
}

func TestParseTestCase_acceptsEmptyInputOutput(t *testing.T) {
	input, output, ok := parseTestCase("---\n")
	if !ok {
		t.Fatal("incorrectly rejects empty input/output")
	}

	if input != "" {
		t.Errorf("input is nonempty: %q", input)
	}

	if output != "" {
		t.Errorf("output is nonempty: %q", output)
	}
}

func TestParseTestCase_parsesInputOutput(t *testing.T) {
	expectedInput := "1 2 3\n4 5 6\n"
	expectedOutput := "a b c\nd e f\n"
	doc := expectedInput + "---\n" + expectedOutput

	input, output, ok := parseTestCase(doc)
	if !ok {
		t.Fatal("incorrectly rejects a document")
	}

	if input != expectedInput {
		t.Errorf("input: actual %q expect %q", input, expectedInput)
	}

	if output != expectedOutput {
		t.Errorf("input: actual %q expect %q", output, expectedOutput)
	}
}
