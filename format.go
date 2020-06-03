package attest

import (
	"fmt"
	"time"

	"github.com/snsinfu/attest/colors"
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
