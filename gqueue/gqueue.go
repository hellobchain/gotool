package gqueue

import "sync"

type Queue struct {
	mu   sync.Mutex
	data []interface{}
}

func New() *Queue {
	return &Queue{}
}

func (q *Queue) Put(v interface{}) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.data = append(q.data, v)
}

func (q *Queue) Take() interface{} {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.data) == 0 {
		return nil
	}
	v := q.data[0]
	if len(q.data) == 1 {
		q.data = nil
	} else {
		q.data = q.data[1:]
	}
	return v
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.data)
}
