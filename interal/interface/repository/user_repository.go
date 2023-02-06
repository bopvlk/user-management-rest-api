package repository

import (
	"context"
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/repository"
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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) FindAllUsers(ctx context.Context, users []*models.User) ([]*models.User, error) {
	dbuser := make([]*DBUser, len(users))
	for i := 0; i < len(users); i++ {
		dbuser[i] = modelUserToDBUser(users[i])
	}

	err := ur.db.WithContext(ctx).Find(&dbuser).Error
	if err != nil {
		return nil, err
	}

	users = make([]*models.User, len(dbuser))
	for i := 0; i < len(dbuser); i++ {
		users[i] = dbUserToModelUser(dbuser[i])
	}

	return users, nil
}

func (ur *userRepository) FindOneUser(ctx context.Context, user *models.User) (*models.User, error) {
	dbuser := modelUserToDBUser(user)

	if err := ur.db.WithContext(ctx).Where(&dbuser).First(&dbuser).Error; err != nil {
		return nil, err
	}

	user = dbUserToModelUser(dbuser)

	return user, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	dbuser := modelUserToDBUser(user)

	if err := ur.db.WithContext(ctx).Create(&dbuser).Error; err != nil {
		return nil, err
	}
	user = dbUserToModelUser(dbuser)
	return user, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, user *models.User) error {
	dbuser := modelUserToDBUser(user)

	if err := ur.db.WithContext(ctx).Delete(&dbuser).Error; err != nil {
		return err
	}
	return nil
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

func dbUserToModelUser(u *DBUser) *models.User {
	var deleteTime *time.Time
	if u.DeletedAt.Valid {
		deleteTime = &u.DeletedAt.Time
	} else {
		deleteTime = nil
	}
	return &models.User{
		ID:        u.ID,
		UserName:  u.UserName,
		FirsnName: u.FirsnName,
		LastName:  u.LastName,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
		DeletedAt: deleteTime,
	}

}
