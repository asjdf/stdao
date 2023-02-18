package stdao

import (
	"github.com/asjdf/gorm-cache/cache"
	"github.com/asjdf/gorm-cache/config"
	"gorm.io/gorm"
)

func CreateWithCache[T any](m T) StdWithCache[T] {
	return StdWithCache[T]{Std: Std[T]{model: m}}
}

type StdWithCache[T any] struct {
	Std[T]
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
