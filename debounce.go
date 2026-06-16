package hotconf

import (
	"sync"
	"time"
)

type Debouncer struct {
	wait  time.Duration
	timer *time.Timer
	mu    sync.Mutex
}

func NewDebouncer(wait time.Duration) *Debouncer {
	return &Debouncer{
		wait: wait,
		mu:   sync.Mutex{},
	}
}

func (d *Debouncer) Do(fn func()) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.timer != nil {
		d.timer.Stop()
	}
	d.timer = time.AfterFunc(d.wait, fn)
}
