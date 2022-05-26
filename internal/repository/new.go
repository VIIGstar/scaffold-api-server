package repository

import (
	"logur.dev/logur"
	"scaffold-api-server/internal/services/cache"
	"scaffold-api-server/internal/services/database"
)

// Registry ...
type Registry interface {
	Database() DatabaseRepo
	Cache() CacheRepo
}

// New ...
func New(logger logur.LoggerFacade, db *database.DB, redis *cache.Client) Registry {
	return impl{
		dbRepo:    NewDBImpl(logger, db),
		cacheRepo: NewCacheImpl(logger, redis),
	}
}

type impl struct {
	cacheRepo CacheRepo
	dbRepo    DatabaseRepo
}

func (i impl) Database() DatabaseRepo {
	return i.dbRepo
}

func (i impl) Cache() CacheRepo {
	return i.cacheRepo
}
