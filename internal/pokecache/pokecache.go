package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}
type Cache struct {
	entries  map[string]cacheEntry
	mu       sync.Mutex
	interval time.Duration
}

func NewCache(interval time.Duration) Cache {
	c := Cache{}
	c.entries = make(map[string]cacheEntry)
	c.interval = interval
	go c.reapLoop()
	return c
}
func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	c.entries[key] = cacheEntry{createdAt: time.Now(), val: val}
	defer c.mu.Unlock()
}
func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	value, ok := c.entries[key]
	if ok {
		return value.val, true
	}
	return nil, false
}
func (c *Cache) reapLoop() {
	tock := time.NewTicker(c.interval)
	for range tock.C {
		c.mu.Lock()
		defer c.mu.Unlock()
		for key, entry := range c.entries {
			age := time.Since(entry.createdAt)
			if age > c.interval {
				delete(c.entries, key)
			}
		}
	}
}
