package gcache

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key        string
	value      interface{}
	expireTime time.Time
}

type Cache struct {
	cap int
	ll  *list.List
	mp  map[string]*list.Element
	mu  sync.Mutex
}

func New(capacity int) *Cache {
	return &Cache{
		cap: capacity,
		ll:  list.New(),
		mp:  make(map[string]*list.Element),
	}
}

func (c *Cache) Set(key string, val interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ee, ok := c.mp[key]; ok {
		c.ll.MoveToFront(ee)
		en := ee.Value.(*entry)
		en.value = val
		en.expireTime = time.Now().Add(ttl)
		return
	}
	en := &entry{key: key, value: val, expireTime: time.Now().Add(ttl)}
	ele := c.ll.PushFront(en)
	c.mp[key] = ele
	if c.ll.Len() > c.cap {
		c.removeOldest()
	}
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if ele, ok := c.mp[key]; ok {
		en := ele.Value.(*entry)
		if time.Now().After(en.expireTime) {
			c.removeElement(ele)
			return nil, false
		}
		c.ll.MoveToFront(ele)
		return en.value, true
	}
	return nil, false
}

func (c *Cache) removeOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.removeElement(ele)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.ll.Remove(e)
	en := e.Value.(*entry)
	delete(c.mp, en.key)
}
