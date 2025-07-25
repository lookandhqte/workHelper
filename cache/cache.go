package cache

import (
	"sync"
	"time"
)

type Cache struct {
	data map[int]interface{}
	mu   sync.RWMutex
}

func NewCache() *Cache {
	return &Cache{
		data: make(map[int]interface{}),
	}
}

func (c *Cache) Set(key int, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[key] = value

	time.AfterFunc(ttl, func() {
		c.Delete(key)
	})
}

func (c *Cache) Get(key int) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) Delete(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, key)
}
