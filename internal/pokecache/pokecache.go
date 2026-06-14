package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	content map[string]cacheEntry
	mu      sync.RWMutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		content: make(map[string]cacheEntry),
	}
	go (&cache).reapLoop(interval)
	return &cache
}

func (cache *Cache) Add(key string, value []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.content[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	entry, ok := cache.content[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (cache *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		cache.mu.Lock()
		for key, value := range cache.content {
			//timeNow := time.Now()
			//timeCacheEntryAlive := timeNow.Sub(value.createdAt)
			timeCacheEntryAlive := time.Since(value.createdAt)
			if timeCacheEntryAlive > interval {

				delete(cache.content, key)
			}
		}
		// omit defer from 'cache.mu.Unlock()' in infinite loops, you have lock manually
		cache.mu.Unlock()
	}

}
