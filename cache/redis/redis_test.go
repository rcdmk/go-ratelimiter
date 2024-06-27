package redis_test

import (
	"errors"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/redis/go-redis/v9"

	"github.com/rcdmk/go-ratelimiter/cache"
	cacheRedis "github.com/rcdmk/go-ratelimiter/cache/redis"
)

func newMockedRedis(t *testing.T) (*redis.Client, *miniredis.Miniredis) {
	miniRedis := miniredis.RunT(t)

	redisClient := redis.NewClient(&redis.Options{
		Addr: miniRedis.Addr(),
	})

	return redisClient, miniRedis
}

func Test_Redis_Cache_Can_Store_And_Retrieve_Values_For_A_Given_Key(t *testing.T) {
	key1 := "test-key1"
	value1 := 42
	key2 := "test-key2"
	value2 := 84

	redisClient, _ := newMockedRedis(t)

	memCache := cacheRedis.New(redisClient)

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

func Test_Redis_Cache_Can_Store_And_Retrieve_Values_For_A_Given_Key_Within_Expiration(t *testing.T) {
	key := "test-key"
	value := 42

	redisClient, miniRedis := newMockedRedis(t)
	memCache := cacheRedis.New(redisClient)

	_ = memCache.SetWithExpiration(key, value, 5*time.Millisecond)
	miniRedis.FastForward(4 * time.Millisecond)
	retrievedValue, err := memCache.Get(key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrievedValue != value {
		t.Errorf("Expected value %d, got %d", value, retrievedValue)
	}
}

func Test_Redis_Cache_Cant_Retrieve_Values_For_A_Given_Expired_Key(t *testing.T) {
	key := "test-key"
	value := 42

	redisClient, miniRedis := newMockedRedis(t)
	memCache := cacheRedis.New(redisClient)

	_ = memCache.SetWithExpiration(key, value, 2*time.Millisecond)
	miniRedis.FastForward(3 * time.Millisecond)

	retrievedValue, err := memCache.Get(key)

	if !errors.Is(err, cache.ErrCacheMiss) {
		t.Errorf("Expected error %v, got %v", cache.ErrCacheMiss, err)
	}

	if retrievedValue != 0 {
		t.Errorf("Expected value to be zero, got %d", retrievedValue)
	}
}
