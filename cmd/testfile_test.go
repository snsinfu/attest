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
