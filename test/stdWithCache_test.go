package test

import (
	"github.com/Pacific73/gorm-cache/config"
	"github.com/asjdf/stdao"
	"testing"
)

func TestStdWithCacheDAO(t *testing.T) {
	UserDAO := stdao.CreateWithCache(&user{})
	UserDAO.Init(db, &config.CacheConfig{
		CacheLevel:           config.CacheLevelAll,
		CacheStorage:         config.CacheStorageMemory,
		InvalidateWhenUpdate: true,
		CacheTTL:             10000,
		DebugMode:            true,
	})
}
