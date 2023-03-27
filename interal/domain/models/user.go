package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `json:"id" gorm:"primary_key,"`
	Role         string         `json:"role"`
	Rating       int            `json:"rating"`
	RatedByUsers []RatedByUser  `json:"rate_by_users"`
	UserName     string         `json:"user_name" gorm:"unique"`
	FirstName    string         `json:"first_name"`
	LastName     string         `json:"last_name"`
	Password     string         `json:"password"`
	CreatedAt    *time.Time     `json:"created_at"`
	UpdatedAt    *time.Time     `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}

type RatedByUser struct {
	ID            uint           `json:"id"`
	UserID        uint           `json:"user_id"`
	RatedByUserID uint           `json:"rated_by_user_id"`
	Rate          string         `json:"rate"`
	CreatedAt     *time.Time     `json:"created_at"`
	UpdatedAt     *time.Time     `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `json:"deleted_at" gorm:"index"`
}
