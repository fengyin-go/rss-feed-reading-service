package feedstate

import "sync"

type Cache struct {
	mu    sync.RWMutex
	value []byte
}

func (c *Cache) Save(value []byte) { c.mu.Lock(); defer c.mu.Unlock(); c.value = value }
func (c *Cache) Load() []byte      { c.mu.RLock(); defer c.mu.RUnlock(); return c.value }
