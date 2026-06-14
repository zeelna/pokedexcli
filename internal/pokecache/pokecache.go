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

// Create a new Cache, to store HTTP Request URL as key, and HTTP Response Body as value. A cache entry is deleted if older than 'interval' seconds
func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		content: make(map[string]cacheEntry),
	}
	go (&cache).reapLoop(interval)
	return &cache
}

// Add Cache entry, key is the URL, value is the HTTP Response Body
func (cache *Cache) Add(key string, value []byte) {
	cache.mu.Lock()
	defer cache.mu.Unlock()

	cache.content[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
}

// Get a Cache entry by a URL as key passed by argument
func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mu.RLock()
	defer cache.mu.RUnlock()

	entry, ok := cache.content[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

// Delete ('reap') an entry in the cache once it's createAt time has become larger than 'interval' (5seconds by config currently)
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
