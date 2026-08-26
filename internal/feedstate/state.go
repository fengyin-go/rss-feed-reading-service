package feedstate

import "sync"

type Cache struct {
	mu    sync.RWMutex
	value []byte
}

func (c *Cache) Save(value []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	cp := make([]byte, len(value))
	copy(cp, value)
	c.value = cp
}

func (c *Cache) Load() []byte {
	c.mu.RLock()
	defer c.mu.RUnlock()
	cp := make([]byte, len(c.value))
	copy(cp, c.value)
	return cp
}
