package db

import (
	"fmt"
	"sync"
	"testing"
)

func TestCacheLRU_Eviction(t *testing.T) {
	cache := NewCacheLRU(2)

	cache.Set("a", "1")
	cache.Set("b", "2")

	// Access "a" to make it Most Recently Used (MRU)
	if val, _ := cache.Get("a"); val == "" {
		t.Fatal("Expected to find key 'a'")
	}

	// Adding "c" should evict "b" (since "a" was recently promoted)
	cache.Set("c", "3")

	if val, exists := cache.Get("b"); exists || val != "" {
		t.Error("Key 'b' should have been evicted")
	}

	if val, exists := cache.Get("a"); !exists {
		t.Errorf("Expected key 'a' to persist with value '1', got val=%s", val)
	}
}

func TestCacheLRU_Concurrency(t *testing.T) {
	cache := NewCacheLRU(10)
	var wg sync.WaitGroup
	workers := 50

	// Concurrent Writers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key-%d", id), fmt.Sprintf("val-%d", id))
		}(i)
	}

	// Concurrent Readers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			cache.Get(fmt.Sprintf("key-%d", id))
		}(i)
	}

	wg.Wait() // Run with `go test -race` to verify thread safety
}
