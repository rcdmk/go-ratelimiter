package cache_test

import (
	"errors"
	"testing"
	"time"

	"github.com/rcdmk/go-ratelimiter/cache"
)

func Test_InMemory_Cache_Can_Store_And_Retrieve_Values_For_A_Given_Key(t *testing.T) {
	key1 := "test-key1"
	value1 := 42
	key2 := "test-key2"
	value2 := 84

	memCache := cache.NewInMemory()

	_ = memCache.Set(key1, value1)
	_ = memCache.Set(key2, value2)

	retrievedValue, err := memCache.Get(key1)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedValue != value1 {
		t.Errorf("Expected value %d, got %d", value1, retrievedValue)
	}

	retrievedValue, err = memCache.Get(key2)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedValue != value2 {
		t.Errorf("Expected value %d, got %d", value1, retrievedValue)
	}
}

func Test_InMemory_Cache_Can_Store_And_Retrieve_Values_For_A_Given_Key_Within_Expiration(t *testing.T) {
	key := "test-key"
	value := 42

	memCache := cache.NewInMemory()

	_ = memCache.SetWithExpiration(key, value, 2*time.Millisecond)
	time.Sleep(1 * time.Millisecond)
	retrievedValue, err := memCache.Get(key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Expected value %d, got %d", value, retrievedValue)
	}
}

func Test_InMemory_Cache_Cant_Retrieve_Values_For_A_Given_Expired_Key(t *testing.T) {
	key := "test-key"
	value := 42

	memCache := cache.NewInMemory()

	_ = memCache.SetWithExpiration(key, value, 2*time.Millisecond)
	time.Sleep(3 * time.Millisecond)

	retrievedValue, err := memCache.Get(key)

	if !errors.Is(err, cache.ErrCacheMiss) {
		t.Errorf("Expected error %v, got %v", cache.ErrCacheMiss, err)
	}

	if retrievedValue != 0 {
		t.Errorf("Expected value to be zero, got %d", retrievedValue)
	}
}
