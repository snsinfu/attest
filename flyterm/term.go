// Package flyterm provides functions for live-updating terminal like flight
// information display.
package flyterm

import (
	"fmt"
	"time"

	"github.com/gosuri/uilive"
)

const (
	// uilive refreshes display extremely frequently, which casues flickering.
	// It does not support disabling the auto-refresh feature, so we set the
	// refresh interval very long as a workaround.
	uiliveAutoInterval = 1000 * time.Second

	defaultQueuePerRow = 10
	defaultInterval    = time.Second / 10
)

type Term struct {
	rows     int
	updates  chan update
	quit     chan bool
	interval time.Duration
}

type update struct {
	Row  int
	Text string
}

type Options struct {
	Queue    int
	Interval time.Duration
}

func New(rows int, opts Options) *Term {
	if opts.Queue == 0 {
		opts.Queue = defaultQueuePerRow * rows
	}
	if opts.Interval == 0 {
		opts.Interval = defaultInterval
	}

	t := &Term{
		rows:     rows,
		updates:  make(chan update, opts.Queue),
		quit:     make(chan bool),
		interval: opts.Interval,
	}
	go t.start()

	return t
}

func (t *Term) Update(row int, text string) {
	t.updates <- update{
		Row:  row,
		Text: text,
	}
}

func (t *Term) start() {
	w := uilive.New()
	w.RefreshInterval = uiliveAutoInterval
	w.Start()
	defer w.Stop()

	ticker := time.NewTicker(t.interval)
	defer ticker.Stop()

	rows := make([]string, t.rows)
	end := false

	for !end {
		select {
		case <-t.quit:
			// Caller likely sends a last update right before calling Stop().
			// So, we make sure all updates are processed before quitting.
			end = true
		case <-ticker.C:
		}

	pump:
		for {
			select {
			case up := <-t.updates:
				rows[up.Row] = up.Text
			default:
				break pump
			}
		}

		for _, row := range rows {
			fmt.Fprintln(w, row)
		}
		w.Flush()
	}
	close(t.quit)
}

func (t *Term) Stop() {
	t.quit <- true
	<-t.quit
}
