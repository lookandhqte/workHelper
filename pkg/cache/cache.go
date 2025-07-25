package cache

import (
	"sync"
	"time"
)

// Потокобезопасное in-memory хранилище
type Cache struct {
	data map[int]interface{}
	mu   sync.RWMutex
}

// Cоздает и возвращает новый экземпляр Cache
func NewCache() *Cache {
	return &Cache{
		data: make(map[int]interface{}),
	}
}

// Добавляет значение в кэш с указанным временем жизни (TTL)
func (c *Cache) Set(key int, value interface{}, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.data[key] = value

	// Установка таймера для автоматического удаления записи по истечении TTL
	time.AfterFunc(ttl, func() {
		c.Delete(key)
	})
}

// Получение значения из кэша по ключу
func (c *Cache) Get(key int) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[key]
	return val, ok
}

// Удаление значения из кэша по ключу
func (c *Cache) Delete(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}
