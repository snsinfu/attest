package periodic

import "time"

// A Periodic manages periodic invocation of a function.
type Periodic struct {
	quit chan bool
}

// New starts a goroutine periodically calls f. The periodic call is stopped by
// calling Stop() on the returned Periodic object.
func New(d time.Duration, f func()) *Periodic {
	quit := make(chan bool)
	go func() {
		ticker := time.NewTicker(d)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				f()
			case <-quit:
				close(quit)
				return
			}
		}
	}()
	return &Periodic{quit}
}

// Stop stops the periodic call started on p.
func (p *Periodic) Stop() {
	p.quit <- true
	<-p.quit
}
