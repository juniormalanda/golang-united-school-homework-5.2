package cache

import "time"

type cacheDeadline struct {
	set  bool
	till time.Time
}

func (d *cacheDeadline) expired() bool {
	if !d.set {
		return false
	}

	now := time.Now()

	return now.Before(d.till)
}

type cacheItem struct {
	value    string
	deadline cacheDeadline
}

type Cache struct {
	items map[string]cacheItem
}

func NewCache() Cache {
	return Cache{items: make(map[string]cacheItem)}
}

func (cache *Cache) Get(key string) (string, bool) {
	item, ok := cache.items[key]

	if !ok {
		return "", false
	}

	if item.deadline.expired() {
		return "", false
	}

	return item.value, true
}

func (cache *Cache) Put(key, value string) {
	cache.items[key] = cacheItem{value: value, deadline: cacheDeadline{}}
}

func (cache *Cache) Keys() []string {
	keys := make([]string, 0)

	for k, item := range cache.items {
		if !item.deadline.expired() {
			keys = append(keys, k)
		}
	}

	return keys
}

func (cache *Cache) PutTill(key, value string, deadline time.Time) {
	cache.items[key] = cacheItem{value: value, deadline: cacheDeadline{set: true, till: deadline}}
}
