package gfuture

import (
	"context"
	"sync"
	"time"
)

type Future struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

func Go(fn func() (interface{}, error)) *Future {
	f := &Future{}
	f.wg.Add(1)
	go func() {
		defer f.wg.Done()
		f.val, f.err = fn()
	}()
	return f
}

func (f *Future) Get(ctx context.Context) (interface{}, error) {
	done := make(chan struct{})
	go func() {
		f.wg.Wait()
		close(done)
	}()
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-done:
		return f.val, f.err
	}
}
func (f *Future) GetTimeout(timeout time.Duration) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return f.Get(ctx)
}
