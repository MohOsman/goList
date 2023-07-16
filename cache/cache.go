package cache

import (
	"github.com/bluele/gcache"
	"time"
)

var Cache gcache.Cache

func init() {
	// Initialize the cache
	Cache = gcache.New(100).LRU().Build()
}

func Set(key string, value interface{}) {
	Cache.Set(key, value)
}

func Get(key string) ([]byte, error) {
	value, err := Cache.Get(key)
	return value.([]byte), err
}

func SetWithExpire(key string, value interface{}, expire time.Duration) {
	Cache.SetWithExpire(key, value, expire)
}