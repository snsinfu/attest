package command

import "io"

type closerT struct {
	c io.Closer
}

func closer(c io.Closer) *closerT {
	return &closerT{c: c}
}

func (c *closerT) Close() error {
	if c.c == nil {
		return nil
	}
	return c.c.Close()
}

func (c *closerT) NoClose() {
	c.c = nil
}
