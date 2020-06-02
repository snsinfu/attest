package attest

import "time"

// A Config specifies how tests should be run.
type Config struct {
	Command []string
	Tests   []string
	MaxJobs int
	Timeout time.Duration
}
