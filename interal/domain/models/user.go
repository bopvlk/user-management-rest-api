package models

import "time"

type User struct {
	ID        uint       `json:"id"`
	UserName  string     `json:"user_name"`
	FirsnName string     `json:"first_name"`
	LastName  string     `json:"last_name"`
	Password  string     `json:"password"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
