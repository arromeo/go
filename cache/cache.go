package cache

import (
	"container/list"
	"sync"
	"time"
)

type CacheItem struct {
	Key        string
	Value      interface{}
	Expiration time.Time
}

type Cache struct {
	capacity int
	items    map[string]*list.Element
	order    *list.List
	lock     sync.RWMutex
}

func NewCache(capacity int) *Cache {
	return &Cache{
		capacity: capacity,
		items:    make(map[string]*list.Element),
		order:    list.New(),
	}
}

func (c *Cache) Set(key string, value interface{}, ttl time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, exists := c.items[key]; exists {
		c.order.MoveToFront(element)
		element.Value.(*CacheItem).Value = value
		element.Value.(*CacheItem).Expiration = time.Now().Add(ttl)
		return
	}

	if c.order.Len() >= c.capacity {
		c.evict()
	}

	item := &CacheItem{
		Key:        key,
		Value:      value,
		Expiration: time.Now().Add(ttl),
	}
	element := c.order.PushFront(item)
	c.items[key] = element
}

func (c *Cache) Get(key string) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	element, exists := c.items[key]
	if !exists {
		return nil, false
	}

	item := element.Value.(*CacheItem)
	if time.Now().After(item.Expiration) {
		c.lock.RUnlock()
		c.lock.Lock()
		c.removeElement(element)
		c.lock.Unlock()
		return nil, false
	}

	c.lock.RUnlock()
	c.lock.Lock()
	c.order.MoveToFront(element)
	c.lock.Unlock()
	c.lock.RLock()

	return item.Value, true
}

func (c *Cache) evict() {
	element := c.order.Back()
	if element != nil {
		c.removeElement(element)
	}
}

func (c *Cache) removeElement(element *list.Element) {
	item := element.Value.(*CacheItem)
	delete(c.items, item.Key)
	c.order.Remove(element)
}

func (c *Cache) Delete(key string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if element, exists := c.items[key]; exists {
		c.removeElement(element)
	}
}

func (c *Cache) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items = make(map[string]*list.Element)
	c.order.Init()
}
