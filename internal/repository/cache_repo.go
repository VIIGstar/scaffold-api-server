package repository

import (
	"context"
	"logur.dev/logur"
	"scaffold-api-server/internal/services/cache"
	"time"
)

type CacheRepo interface {
	Ping(ctx context.Context) error
	Set(ctx context.Context, key string, value interface{}, duration time.Duration) error
	Get(ctx context.Context, key string, value interface{}) error
	Delete(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	HMSet(ctx context.Context, hashKey string, data interface{}) error
	HMGet(ctx context.Context, hashKey string, arrValue interface{}, emptyJson string, keys ...string) error
}

func NewCacheImpl(logger logur.LoggerFacade, redis *cache.Client) CacheRepo {
	return cacheImpl{
		logger: logger,
		Client: redis,
	}
}

type cacheImpl struct {
	logger logur.LoggerFacade
	*cache.Client
}
