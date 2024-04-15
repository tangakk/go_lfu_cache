package go_lfu_cache

import (
	"cmp"
	"slices"
	"sync"
)

// Cache value structure
type CacheValue struct {
	value any
	freq  int
	key   string
}

// Cache structure
// if len > UpperBound, cache will be cleaned to the LowerBound
// if UpperBound or LowerBound is set to 0, nothing will happen
type Cache struct {
	UpperBound int
	LowerBound int
	lock       sync.Mutex
	values     map[string]*CacheValue
}

// Creates new Cache instance
func New() *Cache {
	return &Cache{
		0, 0, sync.Mutex{}, make(map[string]*CacheValue),
	}
}

// Checks if cache has this key
func (c *Cache) Has(key string) bool {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.values[key]

	return ok
}

// Return value with key if it exists and increases freq for it
// or nil if it doesn't
func (c *Cache) Get(key string) interface{} {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		c.increment(e)
		return e.value
	}
	return nil
}

// idk that was in big Ya's course
func (c *Cache) increment(e *CacheValue) {
	e.freq++
}

// Saves values with key
func (c *Cache) Set(key string, value interface{}) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if e, ok := c.values[key]; ok {
		e.value = value
	} else {
		c.lock.Unlock()
		if len(c.values) >= c.UpperBound && c.UpperBound != 0 && c.LowerBound != 0 {
			c.Evict(c.UpperBound - c.LowerBound)
		}
		c.lock.Lock()
		c.values[key] = &CacheValue{value: value, key: key, freq: 1}
	}
}

// Returns cache len
func (c *Cache) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.values)
}

// Returns freq of cache value with key.
// Freq is increased every time value is being read.
// Freq is set to 1 when elements from cache were deleted.
// Returns 0 if element doesn't exist
func (c *Cache) GetFrequency(key string) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	if val, ok := c.values[key]; ok {
		return val.freq
	}
	return 0
}

// Returns all keys from Cache
func (c *Cache) Keys() []string {
	c.lock.Lock()
	defer c.lock.Unlock()
	keys := make([]string, len(c.values))
	for key := range c.values {
		keys = append(keys, key)
	}
	return keys
}

// Delets least frequently used elements from cache
// Returns number of deleted elements
func (c *Cache) Evict(count int) int {
	c.lock.Lock()
	defer c.lock.Unlock()
	count = min(len(c.values), count)
	temp := make([]*CacheValue, len(c.values))
	i := 0
	for _, value := range c.values {
		temp[i] = value
		i++
	}
	slices.SortStableFunc(temp, func(a *CacheValue, b *CacheValue) int {
		return -cmp.Compare[int](a.freq, b.freq)
	})

	to_delete := temp[len(temp)-count:]
	for _, el := range to_delete {
		delete(c.values, el.key)
	}
	for _, value := range c.values {
		value.freq = 1
	}
	return count
}
