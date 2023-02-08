package repository

import (
	"context"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUsers(ctx context.Context, page int, users []*models.User) ([]*models.User, error)
	FindOneUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) FindUsers(ctx context.Context, page int, users []*models.User) ([]*models.User, error) {
	err := ur.db.WithContext(ctx).Limit(5).Offset((page - 1) * 5).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (ur *userRepository) FindOneUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := ur.db.WithContext(ctx).Where(&user).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := ur.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) DeleteUser(ctx context.Context, user *models.User) error {
	if err := ur.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return err
	}
	return nil
}
