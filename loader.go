package hotconf

import (
	"os"
	"sync"
)

type Config[T any] struct {
	mu  sync.RWMutex
	val T
}

func Load[T any](path string, unmarshal func([]byte, any) error) (T, error) {
	var v T
	var zero T

	bytes, err := os.ReadFile(path)
	if err != nil {
		return zero, err
	}

	if err := unmarshal(bytes, &v); err != nil {
		return zero, err
	}

	return v, nil
}

func (c *Config[T]) Get() T {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.val
}

func (c *Config[T]) Set(v T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.val = v
}
