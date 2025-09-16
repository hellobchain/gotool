package gcron

import (
	"sync"
	"time"
)

type Cron struct {
	ticker *time.Ticker
	stop   chan struct{}
	wg     sync.WaitGroup
}

func Every(interval time.Duration, fn func()) *Cron {
	c := &Cron{stop: make(chan struct{})}
	c.ticker = time.NewTicker(interval)
	c.wg.Add(1)
	go func() {
		defer c.wg.Done()
		for {
			select {
			case <-c.ticker.C:
				fn()
			case <-c.stop:
				return
			}
		}
	}()
	return c
}

func (c *Cron) Stop() {
	c.ticker.Stop()
	close(c.stop)
	c.wg.Wait()
}
