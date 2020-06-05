package attest

import (
	"fmt"
	"strings"
	"time"

	"github.com/snsinfu/attest/colors"
	"github.com/snsinfu/attest/test"
)

var spinner = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

func formatWait(name string) string {
	return fmt.Sprintf("%s  -:--  %s", colors.Gray("WAIT"), name)
}

func formatRun(name string, elapsed time.Duration, spin int) string {
	min, sec := extractMinSec(elapsed)

	return fmt.Sprintf(
		"%s  %d:%02d  %s",
		colors.Yellow("RUN"+spinner[spin%len(spinner)]),
		min,
		sec,
		name,
	)
}

func formatLabel(outcome test.Outcome) string {
	switch outcome {
	case test.TestPassed:
		return colors.Green("PASS")
	case test.TestFailed:
		return colors.Red("FAIL")
	case test.TestTimeout:
		return colors.Blue("TIME")
	case test.TestError:
		return colors.Magenta("DEAD")
	}
	panic("unexpected argument")
}

func formatOutcome(name string, elapsed time.Duration, outcome test.Outcome) string {
	label := formatLabel(outcome)
	min, sec := extractMinSec(elapsed)
	return fmt.Sprintf("%s  %d:%02d  %s", label, min, sec, name)
}

func extractMinSec(d time.Duration) (int, int) {
	sec := int(d.Seconds())
	min := sec / 60
	sec %= 60
	return min, sec
}

func formatResult(tc test.Case, r test.Result) string {
	label := formatLabel(r.Outcome)
	heading := fmt.Sprintf("%s  %s\n", label, tc.Name)

	tcSection := colors.Gray("IN:") + "\n"
	tcSection += endLine(tc.Input)
	tcSection += colors.Gray("OUT:") + "\n"
	tcSection += endLine(tc.Output)
	tcSection = colors.Gray("Test case") + "\n" + indent(tcSection, 2)

	rSection := colors.Gray("OUT:") + "\n"
	rSection += endLine(r.Stdout)
	if len(r.Stderr) != 0 {
		rSection += colors.Gray("LOG:") + "\n"
		rSection += endLine(r.Stderr)
	}
	rSection = colors.Gray("Program output") + "\n" + indent(rSection, 2)

	return heading + indent(tcSection, 2) + indent(rSection, 2)
}

func indent(s string, n int) string {
	lines := strings.Split(s, "\n")
	margin := strings.Repeat(" ", n)

	for i := range lines {
		lines[i] = margin + lines[i]
	}
	return strings.TrimRight(strings.Join(lines, "\n"), " ")
}

func endLine(s string) string {
	if !strings.HasSuffix(s, "\n") {
		s += "\n"
	}
	return s
}
