package model

import (
	"gorm.io/gorm"
)

//go:generate  paramgen
type User struct {
	gorm.Model
	Name     string
	Password string
	Email    string
	IsAdmin  int
}
