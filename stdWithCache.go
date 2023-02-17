package stdao

import (
	"github.com/Pacific73/gorm-cache/cache"
	"github.com/Pacific73/gorm-cache/config"
	"gorm.io/gorm"
)

func CreateWithCache[T any](m T) StdWithCache[T] {
	return StdWithCache[T]{model: m}
}

type StdWithCache[T any] struct {
	// Std is a standard struct for all dao structs
	db    *gorm.DB
	model T
}

func (s *StdWithCache[T]) Init(db *gorm.DB, cacheConfig *config.CacheConfig) error {
	s.db = db
	err := s.db.AutoMigrate(s.model)
	if err != nil {
		return err
	}
	c, err := cache.NewGorm2Cache(cacheConfig)
	if err != nil {
		return err
	}
	err = s.db.Use(c)
	if err != nil {
		return err
	}
	return nil
}
