package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID              uint           `json:"id" gorm:"primary_key,"`
	Role            string         `json:"role"`
	Rating          int            `json:"rating"`
	WhenIGaveRating *time.Time     `json:"my_last_vote"`
	WhoRateds       []WhoRated     `json:"who_rated"`
	UserName        string         `json:"user_name" gorm:"unique"`
	FirstName       string         `json:"first_name"`
	LastName        string         `json:"last_name"`
	Password        string         `json:"password"`
	CreatedAt       *time.Time     `json:"created_at"`
	UpdatedAt       *time.Time     `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type WhoRated struct {
	ID         uint           `json:"rating_id"`
	UserID     uint           `json:"user_id"`
	WhoRatedID uint           `json:"who_rated"`
	Rate       string         `json:"vote"`
	CreatedAt  *time.Time     `json:"created_at"`
	UpdatedAt  *time.Time     `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
