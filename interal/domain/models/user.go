package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primary_key,"`
	Role      string         `json:"role"`
	Rating    Rating         `json:"user_rating"`
	UserName  string         `json:"user_name" gorm:"unique"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	Password  string         `json:"password"`
	CreatedAt *time.Time     `json:"created_at"`
	UpdatedAt *time.Time     `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type Rating struct {
	UserID   uint   `json:"rating_id"`
	Rating   int    `json:"rating"`
	WhoRated string `json:"who_rated"`
}
