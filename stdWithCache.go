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
	cache.StatsAccessor
}

func (s *StdWithCache[T]) Init(db *gorm.DB, cacheConfig *config.CacheConfig) (err error) {
	s.db, err = forkDB(db)
	if err != nil {
		return err
	}
	err = s.db.AutoMigrate(s.model)
	if err != nil {
		return err
	}
	c, err := cache.NewGorm2Cache(cacheConfig)
	if err != nil {
		return err
	}
	s.StatsAccessor = c
	err = s.Use(c)
	if err != nil {
		return err
	}
	return nil
}

// Pluck queries a single column from a model, returning in the slice dest. E.g.:
// var ages []int64
// db.Model(&users).Pluck("age", &ages)
func (s *StdWithCache[T]) Pluck(column string, dest interface{}) *gorm.DB {
	// since cache not support pluck (because it can't get primary key from query resp)
	// so we have to use struct slice instead of int/string slice.

}
