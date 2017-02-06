package memory

import "time"

func NewTicker(d time.Duration) *Ticker {
	t := &Ticker{
		time.NewTicker(d),
		make(chan *time.Time),
		make(chan struct{}),
	}

	go t.loop()

	return t
}

type Ticker struct {
	ticker  *time.Ticker
	C       chan *time.Time
	closeCh chan struct{}
}

func (t *Ticker) Stop() {
	t.ticker.Stop()
	close(t.closeCh)
}

func (t *Ticker) loop() {
	for {
		select {
		case now := <-t.ticker.C:
			t.C <- &now
		case <-t.closeCh:
			t.C <- nil
		}
	}
}
