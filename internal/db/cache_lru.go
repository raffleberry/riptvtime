package db

import (
	"container/list"
	"sync"
)

type cacheItem struct {
	key   string
	value string
}

type CacheLRU struct {
	mu        sync.RWMutex
	limit     int
	evictList *list.List
	items     map[string]*list.Element
}

func NewCacheLRU(limit int) *CacheLRU {
	return &CacheLRU{
		limit:     limit,
		evictList: list.New(),
		items:     make(map[string]*list.Element),
	}
}

func (c *CacheLRU) Set(key string, value string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.evictList.MoveToFront(elem)
		elem.Value.(*cacheItem).value = value
	}

	if c.evictList.Len() >= c.limit {
		c.evictOldest()
	}

	item := &cacheItem{key: key, value: value}
	elem := c.evictList.PushFront(item)
	c.items[key] = elem
}

func (c *CacheLRU) Get(key string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, exists := c.items[key]; exists {
		c.evictList.MoveToFront(elem)
		return elem.Value.(*cacheItem).value, true
	}

	return "", false
}

func (c *CacheLRU) evictOldest() {
	elem := c.evictList.Back()
	if elem != nil {
		c.evictList.Remove(elem)
		item := elem.Value.(*cacheItem)
		delete(c.items, item.key)
	}
}
