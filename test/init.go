package test

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open(":memory:"))
	if err != nil {
		panic(err)
	}
}
