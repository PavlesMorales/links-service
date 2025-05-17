package user

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `json:"name"`
	Email    string `json:"email" grom:"unique" gorm:"uniqueIndex"`
	Password string `json:"password"`
}
