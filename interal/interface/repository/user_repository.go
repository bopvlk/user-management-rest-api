package repository

import (
	"context"
	"math"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"gorm.io/gorm"
)

//go:generate mockgen -destination=../../../gen/mocks/mock_user_repository.go -package=mocks . UserRepository

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindUsers(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, error)
	FindOneUserByID(ctx context.Context, id uint) (*models.User, error)
	FindOneUserByUserNameAndPassword(ctx context.Context, username, password string) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) FindUsers(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, error) {

	offset := (pagination.Page - 1) * pagination.Limit
	users := []*models.User{}
	if err := ur.db.WithContext(ctx).Limit(pagination.Limit).Offset(offset).Order(pagination.Sort).Find(&users).Error; err != nil {
		return nil, nil, err
	}

	if err := ur.db.Model(&models.User{}).Count(&pagination.TotalRows).Error; err != nil {
		return nil, nil, err
	}

	pagination.TotalPages = int(math.Ceil(float64(pagination.TotalRows) / float64(pagination.Limit)))
	return pagination, users, nil
}

func (ur *userRepository) FindOneUserByID(ctx context.Context, id uint) (*models.User, error) {
	user := models.User{}
	if err := ur.db.WithContext(ctx).Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (ur *userRepository) FindOneUserByUserNameAndPassword(ctx context.Context, username, password string) (*models.User, error) {
	user := models.User{}
	if err := ur.db.WithContext(ctx).Where("username = ?", username).Where("password = ?", password).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
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
