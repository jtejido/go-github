package cache

import (
	"fmt"
	"hash/fnv"
	"sync"
	"time"
)

const (
	tempScale     = time.Second
	flushInterval = tempScale
)

var hasher = fnv.New64a()

var (
	keyNotFoundError = "Item with key: '%s' not found"
)

// any message above, and corresponding arguments
func errorf(msg string, args ...interface{}) error {
	return fmt.Errorf(msg, args...)
}

// Cache is intended to be a general purpose K-V store.
type Cache struct {
	sync.RWMutex
	cache map[uint64]*Item
	stop  chan int
}

type Item struct {
	Value  []byte
	Expire *time.Time
}

// Returns a New Cache instance that checks for items about to expire every second.
func New() *Cache {
	c := new(Cache)
	c.stop = make(chan int)
	c.cache = make(map[uint64]*Item)
	c.tempFlush()
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				c.tempFlush()
			}
		}
	}()
	return c
}

func (c *Cache) tempFlush() {
	go func() {
		for {
			select {
			case <-c.stop:
				return
			case <-time.After(flushInterval):
				if len(c.cache) == 0 {
					continue
				}

				c.Lock()
				for id, u := range c.cache {
					if u.Expire != nil {
						if time.Now().After(*u.Expire) {
							delete(c.cache, id)
						}
					}
				}
				c.Unlock()
			}
		}
	}()
}

// Returns the data from a given key
func (c *Cache) Get(key string) (item *Item, err error) {
	c.RLock()
	defer c.RUnlock()

	k := getIndexKey(key)
	item, ok := c.cache[k]

	if !ok {
		return nil, errorf(keyNotFoundError, key)
	}

	return item, nil
}

// Sets the data with a given key
func (c *Cache) Set(key string, item *Item) {
	c.Lock()
	defer c.Unlock()

	k := getIndexKey(key)

	c.cache[k] = item
}

// Deletes a data given a key
func (c *Cache) Delete(key string) (err error) {
	c.Lock()
	defer c.Unlock()

	k := getIndexKey(key)

	if _, ok := c.cache[k]; !ok {
		return errorf(keyNotFoundError, key)
	}

	delete(c.cache, k)

	return nil
}

// Clears the cache's and filter's (if available) items
func (c *Cache) Clear() {
	c.Lock()
	c.cache = make(map[uint64]*Item)
	c.Unlock()
}

func getIndexKey(key string) uint64 {
	hasher.Reset()
	hasher.Write([]byte(key))
	return hasher.Sum64()
}
