package gpool

type Pool struct {
	work chan func()
}

func New(size int) *Pool {
	p := &Pool{work: make(chan func())}
	for i := 0; i < size; i++ {
		go func() {
			for f := range p.work {
				f()
			}
		}()
	}
	return p
}

func (p *Pool) Submit(fn func()) {
	p.work <- fn
}

func (p *Pool) Release() {
	close(p.work)
}
