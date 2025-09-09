package ticker

import (
	"context"
	"time"
)

type TimeTicker struct {
	ticker   *time.Ticker
	period   time.Duration
	callback func()
	done     chan struct{}
	stopped  bool
}

func NewTimeTicker(period time.Duration, callback func()) TimeTicker {
	return TimeTicker{
		period:   period,
		callback: callback,
		done:     make(chan struct{}),
	}
}

func (t *TimeTicker) Start(ctx context.Context) {
	t.ticker = time.NewTicker(t.period)
	t.callback()
	for loop := true; loop; {
		select {
		case <-t.ticker.C:
			t.callback()
		case <-t.done:
			t.stopped = true
			loop = false
		case <-ctx.Done():
			t.stopped = true
			loop = false
		}
	}

	t.ticker.Stop()
	close(t.done)
}

func (t *TimeTicker) Stop() {
	if t.stopped {
		return
	}

	select {
	case t.done <- struct{}{}:
	default:
	}
}
