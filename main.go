package main

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

type ICache interface {
	Cap() int
	Clear()
	Add(key, value interface{})
	//AddWithTTL(key, value interface{}, ttl time.Duration)
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}
type CacheItem struct {
	key   interface{}
	value interface{}
	time  time.Time
}
type LRUCache struct {
	cache    map[interface{}]*list.Element
	maxItems int
	list     *list.List
	mu       sync.Mutex
}

func NewLRUCache(maxItems int) *LRUCache {
	return &LRUCache{
		maxItems: maxItems,
		cache:    make(map[interface{}]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key interface{}) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.cache[key]; ok {
		c.list.MoveToFront(e)
		return e.Value.(*CacheItem).value, true
	}
	return nil, false
}

func (c *LRUCache) removeLast() {
	element := c.list.Back()
	if element != nil {
		item := c.list.Remove(element).(*CacheItem)
		delete(c.cache, item.key)
	}
}

func (c *LRUCache) Add(key interface{}, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if e, ok := c.cache[key]; ok {
		c.list.MoveToFront(e)
		e.Value.(*CacheItem).value = value
		return
	}

	if c.list.Len() >= c.maxItems {
		c.removeLast()
	}
	item := &CacheItem{key, value, time.Now()}
	element := c.list.PushFront(item)
	c.cache[key] = element
}

func (c *LRUCache) Remove(key interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if e, ok := c.cache[key]; ok {
		c.list.Remove(e)
		delete(c.cache, key)
	}
}

func (c *LRUCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache, c.list = make(map[interface{}]*list.Element), list.New()
}

func (c *LRUCache) Cap() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.cache)
}

func main() {
	c := NewLRUCache(5)
	c.Add(1, "sdfsd")
	fmt.Println(c.Cap())
	c.Add(2, "sdfs2d")
	c.Add(3, "sdfs4d")
	fmt.Println(c.Cap())
	c.Add(4, "sdfs5d")
	c.Add(5, "sdfs6d")
	c.Add(6, "sdfs6d")
	fmt.Println(c.Cap())
	c.Remove(6)
	fmt.Println(c.Cap())
	c.Clear()
	fmt.Println(c.Cap())

	v, _ := c.Get(6)
	fmt.Println(v)
}
