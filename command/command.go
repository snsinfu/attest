package command

import (
	"io"
	"os"
	"os/exec"
)

// A Cmd is a running process where stdin and stdout are accessible as
// Writer and Reader interfaces.
type Cmd struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout io.ReadCloser
}

// Run creates a process executing the given command line.
func Run(argv []string) (*Cmd, error) {
	cmd := exec.Command(argv[0], argv[1:]...)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdinC := closer(stdin)
	defer stdinC.Close()

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	stdoutC := closer(stdout)
	defer stdoutC.Close()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	c := &Cmd{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
	}
	stdinC.NoClose()
	stdoutC.NoClose()

	return c, nil
}

// Write writes given content to the stdin of the command.
func (c *Cmd) Write(buf []byte) (int, error) {
	return c.stdin.Write(buf)
}

// WriteEnd closes the stdin of the command.
func (c *Cmd) WriteEnd() error {
	return c.stdin.Close()
}

// Read reads content from the stdout of the command.
func (c *Cmd) Read(buf []byte) (int, error) {
	return c.stdout.Read(buf)
}

// ReadEnd closes the stdout of the command.
func (c *Cmd) ReadEnd() error {
	return c.stdout.Close()
}

// Close releases all resources associated to the command execution. It blocks
// unntil the command ends.
func (c *Cmd) Close() error {
	if err := c.WriteEnd(); err != nil {
		return err
	}
	if err := c.ReadEnd(); err != nil {
		return err
	}
	return c.Wait()
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
