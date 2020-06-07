package command

import (
	"io"
	"os"
	"os/exec"
)

// A Cmd is a running process where stdin, stdout and stderr are piped as Writer
// and Reader instances.
type Cmd struct {
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
	cmd    *exec.Cmd
}

// Run creates a process executing the given command line.
func Run(argv []string) (*Cmd, error) {
	cmd := exec.Command(argv[0], argv[1:]...)

	// No need to close the pipes. See os/exec documentation.
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	c := &Cmd{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
		cmd:    cmd,
	}
	return c, nil
}

// Wait indefinitely waits until the command ends.
func (c *Cmd) Wait() error {
	return c.cmd.Wait()
}

// Signal sends an OS signal to the running command.
func (c *Cmd) Signal(s os.Signal) error {
	return c.cmd.Process.Signal(s)
}

// Status returns the exit status of the command.
func (c *Cmd) Status() *os.ProcessState {
	return c.cmd.ProcessState
}
