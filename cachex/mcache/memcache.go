package mcache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Default value for cache configuration
//noinspection GoUnusedConst
const (
	//NoExpiration :
	//Indicates cached object will never be expired
	NoExpiration = gocache.NoExpiration
	// DefaultExpiration :
	// The default value for expiration
	DefaultExpiration = gocache.DefaultExpiration
)

type Resolver func(key string) interface{}

type MemCache struct {
	cacheInstance *gocache.Cache
	resolver      Resolver
}

func (mc *MemCache) SetOnEvicted(f func(string, interface{})) {
	mc.cacheInstance.OnEvicted(f)
}

func (mc *MemCache) GetStore() *gocache.Cache {
	return mc.cacheInstance
}

func (mc *MemCache) Remove(key string) {
	mc.cacheInstance.Delete(key)
}

func (mc *MemCache) SetResolver(resovler Resolver) {
	mc.resolver = resovler
}

func (mc *MemCache) Set(key string, value interface{}) {
	mc.cacheInstance.SetDefault(key, value)
}

func (mc *MemCache) Get(key string) (interface{}, bool) {
	if mc.resolver == nil {
		return mc.TryGet(key)
	}
	return mc.EnsureGet(key)
}

func (mc *MemCache) TryGet(key string) (interface{}, bool) {
	return mc.cacheInstance.Get(key)
}

func (mc *MemCache) EnsureGet(key string) (interface{}, bool) {
	value, ok := mc.cacheInstance.Get(key)
	if !ok && mc.resolver != nil {
		value = mc.resolver(key)
		ok = value != nil
		if ok {
			mc.Set(key, value)
		}
	}
	return value, ok
}

//export
func NewCache(defaultExpiration, cleanupInterval time.Duration) *MemCache {
	return &MemCache{cacheInstance: gocache.New(defaultExpiration, cleanupInterval)}
}

//export
func NewEnsureCache(defaultExpiration, cleanupInterval time.Duration, resolver Resolver) *MemCache {
	return &MemCache{cacheInstance: gocache.New(defaultExpiration, cleanupInterval), resolver: resolver}
}
