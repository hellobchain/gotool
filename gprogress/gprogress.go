package gprogress

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Bar struct {
	current, total int64
	start          time.Time
	mu             sync.Mutex
}

func New(total int64) *Bar {
	return &Bar{total: total, start: time.Now()}
}

func (b *Bar) Add(n int64) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current += n
	if b.current > b.total {
		b.current = b.total
	}
	b.render()
}

func (b *Bar) Done() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.current = b.total
	b.render()
	println()
}

func (b *Bar) render() {
	rate := float64(b.current) / float64(b.total)
	percent := int(rate * 100)
	filled := int(rate * 50)
	bar := strings.Repeat("█", filled) + strings.Repeat("░", 50-filled)
	elapsed := time.Since(b.start).Seconds()
	eta := ""
	if b.current > 0 {
		etaSec := (float64(b.total-b.current) / float64(b.current)) * elapsed
		eta = fmt.Sprintf(" ETA %.0fs", etaSec)
	}
	fmt.Printf("\r[%s] %d%% %s", bar, percent, eta)
}
