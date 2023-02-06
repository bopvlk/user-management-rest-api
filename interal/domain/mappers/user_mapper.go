package mappers

import (
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"gorm.io/gorm"
)

type DBUser struct {
	ID        uint   `gorm:"primary_key"`
	UserName  string `gorm:"unique"`
	FirsnName string
	LastName  string
	Password  string
	CreatedAt *time.Time
	UpdatedAt *time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func modelUserToDBUser(u *models.User) *DBUser {
	return &DBUser{
		ID:        u.ID,
		UserName:  u.UserName,
		FirsnName: u.FirsnName,
		LastName:  u.LastName,
		Password:  u.Password,
	}
}


