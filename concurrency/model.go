package concurrency

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string
	HashPassword string
	Phone        *string
	Email        *string
}

func (User) TableName() string {
	return "users"
}
