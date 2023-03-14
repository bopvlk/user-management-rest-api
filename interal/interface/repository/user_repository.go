package repository

import (
	"context"
	"errors"
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
	DeleteUserByID(ctx context.Context, id int) error
	DeleteOwnUser(ctx context.Context, id int) error
	UpdateUserByID(ctx context.Context, id int, user *models.User) (*models.User, error)
	UpdateOwnUser(ctx context.Context, id int, user *models.User) (*models.User, error)
	RateUser(ctx context.Context, username string, rate bool) (*models.User, error)
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
	if err := ur.db.WithContext(ctx).Where("user_name = ?", username).Where("password = ?", password).First(&user).Error; err != nil {
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

func (ur *userRepository) DeleteUserByID(ctx context.Context, id int) error {
	user := &models.User{}
	tx := ur.db.WithContext(ctx).Where("id = ?", id).First(&user)
	if tx.Error != nil {
		return tx.Error
	}
	if user.Role == "admin" {
		return errors.New("admin user not allowed to delete")
	}

	if err := ur.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) DeleteOwnUser(ctx context.Context, id int) error {

	if err := ur.db.WithContext(ctx).Delete(&models.User{}, id).Error; err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) UpdateUserByID(ctx context.Context, id int, user *models.User) (*models.User, error) {
	if err := ur.db.WithContext(ctx).Where("id = ?", id).Updates(&user).First(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (ur *userRepository) UpdateOwnUser(ctx context.Context, id int, user *models.User) (*models.User, error) {

	tx := ur.db.WithContext(ctx).Where("id = ?", id).Updates(&user).First(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return user, nil
}

func (ur *userRepository) RateUser(ctx context.Context, username string, rate bool) (*models.User, error) {
	user := new(models.User)
	err := ur.db.WithContext(ctx).Where("user_name = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}

	if rate {
		user.Rating++
	} else {
		user.Rating--
	}

	if err := ur.db.WithContext(ctx).Where("user_name = ?", username).Updates(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}
