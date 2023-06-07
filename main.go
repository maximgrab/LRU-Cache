package main

import (
	"container/list"
	"fmt"
	"reflect"
	"sync"
	"time"
)

type ICache interface {
	Cap() int
	Clear()
	Add(key, value interface{})
	AddWithTTL(key, value interface{}, ttl time.Duration)
	Get(key interface{}) (value interface{}, ok bool)
	Remove(key interface{})
}

type LRUCache struct {
	cache    map[interface{}]*list.Element
	maxItems int
	list     *list.List
	mu       sync.Mutex
}

type CacheItem struct {
	key        interface{}
	value      interface{}
	time       time.Time
	timeToLive time.Duration
}

func NewLRUCache(maxItems int) *LRUCache {
	return &LRUCache{
		maxItems: maxItems,
		cache:    make(map[interface{}]*list.Element),
		list:     list.New(),
	}
}

func (c *LRUCache) Get(key interface{}) (value interface{}, ok bool) {
	if k := !reflect.ValueOf(key).Comparable(); k {
		return nil, false
	}
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

func (c *LRUCache) AddWithTTL(key interface{}, value interface{}, timeToLive time.Duration) {
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
	item := &CacheItem{key, value, time.Now(), timeToLive}
	if timeToLive > 0 {
		time.AfterFunc(timeToLive, func() { c.Remove(key) })
	}
	element := c.list.PushFront(item)
	c.cache[key] = element
}
func (c *LRUCache) Add(key interface{}, value interface{}) (err error) {
	if k := !reflect.ValueOf(key).Comparable(); k {
		return fmt.Errorf("Error, key type is uncomparable")
	}
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
	item := &CacheItem{key, value, time.Now(), 0}
	element := c.list.PushFront(item)
	c.cache[key] = element
	return nil
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
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Add(i, "dsb")
		}(i)
	}
	wg.Wait()
	time.Sleep(1 * time.Second)
	c.AddWithTTL(8, 6, 1*time.Second)
	fmt.Println(c.Cap())
	e := c.Add([]int{1, 3}, 5)
	fmt.Println(e.Error())
	fmt.Println(c.Get([]int{1, 3}))
	fmt.Println(c.Cap())
	fmt.Println(c.Get(2))
	fmt.Println(c.Get(8))
	time.Sleep(2 * time.Second)
	fmt.Println(c.Get(8))
}
