package stdao

import (
	"github.com/asjdf/stdao/page"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create[T any](m T) Std[T] {
	return Std[T]{model: m}
}

type Std[T any] struct {
	// Std is a standard struct for all dao structs
	db    *gorm.DB
	model T
}

func (s *Std[T]) Init(db *gorm.DB) (err error) {
	s.db, err = forkDB(db)
	return s.db.AutoMigrate(s.model)
}

// Create the model to the database.
func (s *Std[T]) Create(model T, tx ...*gorm.DB) (result *gorm.DB) {
	if len(tx) > 0 {
		return tx[0].Create(model)
	}
	return s.db.Create(model)
}

// Save the model to the database.
func (s *Std[T]) Save(model T, tx ...*gorm.DB) (result *gorm.DB) {
	if len(tx) > 0 {
		return tx[0].Save(model)
	}
	return s.db.Save(model)
}

// Update the model to the database.
func (s *Std[T]) Update(model T, tx ...*gorm.DB) (result *gorm.DB) {
	if len(tx) > 0 {
		return tx[0].Updates(model)
	}
	return s.db.Updates(model)
}

func (s *Std[T]) Updates(model T, m map[string]interface{}, tx ...*gorm.DB) (result *gorm.DB) {
	if len(tx) > 0 {
		return tx[0].Model(model).Updates(m)
	}
	return s.db.Model(model).Updates(m)
}

// Delete the model from the database.
func (s *Std[T]) Delete(model T, tx ...*gorm.DB) (result *gorm.DB) {
	if len(tx) > 0 {
		return tx[0].Delete(model)
	}
	return s.db.Delete(model)
}

// Find the model from the database.
func (s *Std[T]) Find(model T, tx ...*gorm.DB) (result *gorm.DB) {
	if len(tx) > 0 {
		return tx[0].Where(model).Find(model)
	}
	return s.db.Find(model)
}

func (s *Std[T]) List(where []clause.Expression, order []clause.OrderByColumn, page page.Paginate, tx ...*gorm.DB) (list []T, result *gorm.DB) {
	list = make([]T, 0)
	var db = s.db
	if len(tx) > 0 {
		db = tx[0]
	}
	result = db.Clauses(clause.Where{Exprs: where}, clause.OrderBy{Columns: order}).Scopes(page.Paginate()).Find(list)
	return
}

func (s *Std[T]) Count() (count int64) {
	s.db.Model(s.model).Count(&count)
	return
}
