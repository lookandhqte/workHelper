package cache

import (
	"amocrm_golang/internal/entity"
	"sync"
	"time"
)

// Потокобезопасное in-memory хранилище
type Cache struct {
	data   map[int]interface{}
	tokens map[int]*entity.Token
	mu     sync.RWMutex
}

// Cоздает и возвращает новый экземпляр Cache
func NewCache() *Cache {
	return &Cache{
		data:   make(map[int]interface{}),
		tokens: make(map[int]*entity.Token),
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

// Добавляет значение в кэш с указанным временем жизни (TTL)
func (c *Cache) SetToken(key int, value *entity.Token, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.tokens[key] = value

	// Установка таймера для автоматического удаления записи по истечении TTL
	time.AfterFunc(ttl, func() {
		c.DeleteToken(key)
	})
}

// Получение значения из кэша по ключу
func (c *Cache) Get(key int) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.data[key]
	return val, ok
}

func (c *Cache) GetToken(key int) (*entity.Token, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	val, ok := c.tokens[key]
	return val, ok
}

// Удаление значения из кэша по ключу
func (c *Cache) Delete(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.data, key)
}

// Удаление значения из кэша по ключу
func (c *Cache) DeleteToken(key int) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.tokens, key)
}
