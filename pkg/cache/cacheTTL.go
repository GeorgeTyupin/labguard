package cache

import (
	"errors"
	"sync"
	"time"
)

type elem[V any] struct {
	value   V
	expDate time.Time
}

type CacheWithTTL[K comparable, V any] struct {
	cache map[K]elem[V]
	TTL   time.Duration
	done  chan struct{}
	mu    *sync.RWMutex
	once  *sync.Once
}

func NewCacheWithTTL[K comparable, V any](ttl time.Duration) CacheWithTTL[K, V] {
	cache := CacheWithTTL[K, V]{
		cache: make(map[K]elem[V]),
		TTL:   ttl,
		done:  make(chan struct{}),
		mu:    &sync.RWMutex{},
		once:  &sync.Once{},
	}

	cache.clearByTTL()

	return cache
}

func (c *CacheWithTTL[K, V]) Set(key K, value V) {
	el := elem[V]{
		value:   value,
		expDate: time.Now().Add(c.TTL),
	}

	c.mu.Lock()
	c.cache[key] = el
	c.mu.Unlock()
}

func (c *CacheWithTTL[K, V]) Get(key K) (V, error) {
	c.mu.RLock()
	el, ok := c.cache[key]
	c.mu.RUnlock()

	var zero V
	if !ok {
		return zero, errors.New("данного элемента нет в кеше")
	}

	if el.expDate.Before(time.Now()) {
		c.Delete(key)
		return zero, errors.New("время жизни данного элемента истекло")
	}

	return el.value, nil
}

func (c *CacheWithTTL[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.cache, key)
}

func (c *CacheWithTTL[K, V]) clearByTTL() {
	ticker := time.NewTicker(time.Hour)
	defer ticker.Stop()

	go func() {
		for {
			select {
			case <-ticker.C:
				c.clear()
			case <-c.done:
				return
			}
		}
	}()
}

func (c *CacheWithTTL[K, V]) clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for key, el := range c.cache {
		if el.expDate.Before(time.Now()) {
			delete(c.cache, key)
		}
	}
}
func (c *CacheWithTTL[K, V]) Stop() {
	c.once.Do(func() {
		close(c.done)
	})
}
