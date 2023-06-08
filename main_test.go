package main

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	c := NewLRUCache(5)
	c.Add(4, 4)
	if r, _ := c.Get(4); r != 4 {
		t.Error("Get function error 4!=4")
	}
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			e, _ := c.Get(4)
			if e != 4 {
				t.Error("Get function sync error")
			}
		}(i)
	}
	wg.Wait()
	if _, e := c.Get(4); !e {
		t.Error("Get check error")
	}
	if _, e := c.Get(8); e {
		t.Error("Get check error")
	}
}

func TestAdd(t *testing.T) {
	c := NewLRUCache(50)
	if e := c.Add([]int{1, 4}, []int{1, 4}); e == nil {
		t.Error("Add function error uncomparable slice passed as key")
	}
	if e := c.Add(func() {}, 8); e == nil {
		t.Error("Add function error uncomparable func passed as key")
	}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Add(i, i)
		}(i)
	}
	wg.Wait()
	for i := 0; i < 50; i++ {
		if _, e := c.Get(i); !e {
			t.Error("Add sync error")
		}
	}
}

func TestTTL(t *testing.T) {
	c := NewLRUCache(5)
	c.AddWithTTL("test", 5, time.Second)
	c.Remove("test")
	c.Add("test", 20)
	<-time.After(2 * time.Second)
	val, ok := c.Get("test")
	assert.True(t, ok)
	v, ok := val.(int)
	assert.True(t, ok)
	assert.Equal(t, 20, v)
}

func TestAddWithTTL(t *testing.T) {
	c := NewLRUCache(50)
	if e := c.AddWithTTL([]int{1, 4}, []int{1, 4}, 1); e == nil {
		t.Error("AddWithTTL function error uncomparable slice passed as key")
	}
	if e := c.AddWithTTL(func() {}, 8, 1); e == nil {
		t.Error("AddWithTTL function error uncomparable func passed as key")
	}
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.AddWithTTL(i, i, 1)
		}(i)
	}
	wg.Wait()
	time.Sleep(time.Second * 1)
	if c.Cap() != 0 {
		t.Error("AddWithTTL sync error")
	}
	c.AddWithTTL(1, 1, 1)
	if c.Cap() != 1 {
		t.Error("AddWithTTL cap error")
	}
}

func TestRemove(t *testing.T) {
	c := NewLRUCache(10)
	c.Add(1, 1)
	c.Add(2, 2)
	c.Add(3, 3)
	c.Remove(1)
	if _, e := c.Get(1); e {
		t.Error("Remove function error")
	}
	c.Remove(2)
	c.Remove(3)
	if c.Cap() != 0 {
		t.Error("Remove function error")
	}
}

func TestRemoveLast(t *testing.T) {
	c := NewLRUCache(10)
	c.Add(1, 1)
	c.Add(2, 2)
	c.Add(3, 3)
	c.removeLast()
	if _, e := c.Get(1); e {
		t.Error("removeLast function error")
	}
	for i := 0; i <= 10; i++ {
		c.Add(i, i)
	}
	c.removeLast()
	if _, e := c.Get(1); e {
		t.Error("removeLast function error")
	}
}

func TestClear(t *testing.T) {
	c := NewLRUCache(10)
	c.Add(1, 1)
	c.Add(2, 2)
	c.Add(3, 3)
	c.Clear()
	if e := c.Cap(); e != 0 {
		t.Error("removeLast function error")
	}
	for i := 0; i < 10; i++ {
		c.Add(i, i)
	}
	for i := 0; i < 10; i++ {
		if _, e := c.Get(i); !e {
			t.Error("Clear reinit error")
		}
	}
}

func CapClear(t *testing.T) {
	c := NewLRUCache(10)
	for i := 1; i <= 10; i++ {
		c.Add(i, i)
		if c.Cap() != i {
			t.Error("Cap func error")
		}
	}
	c.Add(1, 1)
	c.Add(1, 1)
	c.Add(1, 1)
	if c.Cap() != 10 {
		t.Error("Cap func error")
	}
}
