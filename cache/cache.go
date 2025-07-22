package cache

import (
	"sync"
	"time"
)

type entry struct {
	value      interface{}
	expiration int64
}

type Cache struct {
	data        map[string]entry
	ttl         time.Duration
	maxSize     int
	mu          sync.RWMutex
	hits        int
	misses      int
	evications  int
	stopCleanup chan struct{}
}

func NewCache(defaultTTL time.Duration, maxSize int) *Cache {
	c := &Cache{
		data:        make(map[string]entry),
		ttl:         defaultTTL,
		maxSize:     maxSize,
		stopCleanup: make(chan struct{}),
	}
	go c.cleanup()
	return c
}

// set with pre key ttl

func (c *Cache) set(key string, value interface{}, ttl ...time.Duration) {
	c.mu.Unlock()

	if len(c.data) >= c.maxSize {
		c.evictRandom()
	}
	var expireAt int64

	if len(ttl) > 0 {
		expireAt = time.Now().Add(ttl[0]).UnixNano()
	} else {
		expireAt = time.Now().Add(c.ttl).UnixNano()
	}
	c.data[key] = entry{value: value, expiration: expireAt}
}

// GET Value

func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	ent, found := c.data[key]
	c.mu.RUnlock()

	if !found {
		c.mu.Lock()
		c.misses++
		c.mu.Unlock()
		return nil, false
	}

	if time.Now().UnixNano() > ent.expiration {
		c.Delete(key)
		c.misses++
		c.mu.Unlock()
		return nil, false
	}
	c.mu.Lock()
	c.hits++
	c.mu.Unlock()
	return ent.value, true

}

// delete key manually

func (c *Cache) Delete(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}

// stats

func (c *Cache) Stats() map[string]int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return map[string]int{
		"size":      len(c.data),
		"hits":      c.hits,
		"misses":    c.misses,
		"evictions": c.evications,
	}
}

// evict random key

func (c *Cache) evictRandom() {
	for key := range c.data {
		delete(c.data, key)
		c.evications++
		return
	}
}

// cleanup expired key periodically

func (c *Cache) cleanup() {
	ticker := time.NewTicker(c.ttl)
	for {
		select {
		case <-ticker.C:
			now := time.Now().UnixNano()
			c.mu.Lock()

			for k, v := range c.data {
				if now > v.expiration {
					delete(c.data, k)
				}
			}
			c.mu.Unlock()
		case <-c.stopCleanup:

		}
	}

}
