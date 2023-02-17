package test

import "gorm.io/gorm"

type user struct {
	gorm.Model
	Name string
	Age  uint
}
