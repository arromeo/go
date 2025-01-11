package cache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkCacheSet(b *testing.B) {
	cache := NewCache(1000)
	for i := 0; i < b.N; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, time.Minute)
	}
}

func BenchmarkCacheGet(b *testing.B) {
	cache := NewCache(1000)
	for i := 0; i < 1000; i++ {
		cache.Set(fmt.Sprintf("key%d", i), i, time.Minute)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key%d", i%1000))
	}
}

func TestCacheSetAndGet(t *testing.T) {
	cache := NewCache(2)
	cache.Set("key1", "value1", time.Minute)

	value, found := cache.Get("key1")
	if !found || value != "value1" {
		t.Errorf("expected value1, got %v", value)
	}
}

func TestCacheEviction(t *testing.T) {
	cache := NewCache(2)
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value2", time.Minute)
	cache.Set("key3", "value3", time.Minute)

	_, exists := cache.Get("key1")
	if exists {
		t.Error("expected key1 to be evicted, but it was found")
	}
}

func TestCacheTtl(t *testing.T) {
	cache := NewCache(2)
	cache.Set("key1", "value1", time.Second)
	time.Sleep(2 * time.Second)

	_, found := cache.Get("key1")
	if found {
		t.Error("expected key1 to expire, but it was found")
	}
}

func TestCacheUpdateKey(t *testing.T) {
	cache := NewCache(2)
	cache.Set("key1", "value1", time.Minute)
	cache.Set("key2", "value1", time.Minute)
	cache.Set("key1", "value1-updated", time.Minute)

	value, _ := cache.Get("key1")
	if value != "value1-updated" {
		t.Errorf("expected value1-updated, got %v", value)
	}
}

func TestCacheThreadSafety(t *testing.T) {
	cache := NewCache(10)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key%d", i), i, time.Minute)
			cache.Get(fmt.Sprintf("key%d", i))
		}(i)
	}
	wg.Wait()
}
