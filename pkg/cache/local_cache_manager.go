package cache

import (
	"time"
)

const LimitKeys = 500

var cacheManager *localCacheManager

func init() {
	if cacheManager == nil {
		cacheManager = &localCacheManager{
			storage: map[string]*storedValue{},
		}
	}
}

type storedValue struct {
	ticker       *time.Timer
	value        interface{}
	isCompressed bool
}

type localCacheManager struct {
	storage map[string]*storedValue
}

func Get(key string) (interface{}, error) {
	if cacheManager.storage[key] == nil {
		return nil, NilCacheError
	}

	return cacheManager.storage[key].value, nil
}

func Set(key string, value interface{}, expiry time.Duration) {
	cacheManager.storage[key] = &storedValue{
		value:        value,
		isCompressed: false,
	}
	go setExpired(key, expiry)
}

func Delete(key string) (existed bool) {
	existed = cacheManager.storage[key] != nil
	delete(cacheManager.storage, key)
	return existed
}

func setExpired(key string, expiry time.Duration) {
	ref := cacheManager.storage[key]
	if ref != nil {
		ref.ticker = time.AfterFunc(expiry, func() {
			Delete(key)
		})
	}
}
