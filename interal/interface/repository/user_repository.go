package repository

import (
	"context"
	"errors"
	"math"
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
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
	RateUserByUsername(ctx context.Context, userWhoRateID uint, username, rate string) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (ur *userRepository) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	if err := ur.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
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

func (ur *userRepository) RateUserByUsername(ctx context.Context, rateUserID uint, username, rate string) (*models.User, error) {
	user := &models.User{}
	if err := ur.db.WithContext(ctx).Where("user_name = ?", username).Preload("RatedByUsers").First(&user).Error; err != nil {
		return nil, apperrors.UserNotFoundErr.AppendMessage(err)
	}

	existingRateUser, rating, err := determineUserRate(user.Rating, user.RatedByUsers, rateUserID, rate)
	if err != nil {
		return nil, err
	}
	user.Rating = rating

	if existingRateUser != nil {
		existingRateUser.Rate = rate
		if err := ur.db.WithContext(ctx).Where("id = ?", existingRateUser.ID).Updates(&existingRateUser).Error; err != nil {
			return nil, apperrors.CanNotUpdateErr.AppendMessage(err)
		}
	} else {
		if err := ur.db.WithContext(ctx).Model(&user).Association("RatedByUsers").Append(&models.RatedByUser{RatedByUserID: rateUserID, Rate: rate}); err != nil {
			return nil, apperrors.CanNotCreateTableErr.AppendMessage(err)
		}
	}

	if err := ur.db.WithContext(ctx).Where("user_name = ?", username).UpdateColumns(&models.User{Rating: user.Rating}).Error; err != nil {
		return nil, apperrors.CanNotUpdateErr.AppendMessage(err)
	}

	return user, nil
}

func determineUserRate(userRating int, ratedUsers []models.RatedByUser, ratedUserId uint, rate string) (*models.RatedByUser, int, error) {
	now := time.Now()

	for _, u := range ratedUsers {

		if u.RatedByUserID != ratedUserId {
			continue
		}

		if u.CreatedAt.Add(time.Hour * 1).After(now) {
			return nil, 0, &apperrors.ProblemWithGivingRating
		}

		switch u.Rate {
		case rate:
			return nil, 0, &apperrors.CanNotRateAgain
		case "up":
			if rate == "rm" {
				userRating--
				return &u, userRating, nil
			}

			userRating -= 2
			return &u, userRating, nil
		case "rm":
			if rate == "up" {
				userRating++
				return &u, userRating, nil
			}

			userRating--
			return &u, userRating, nil
		case "down":
			if rate == "up" {
				userRating += 2
				return &u, userRating, nil
			}

			userRating++
			return &u, userRating, nil
		default:
			return nil, 0, &apperrors.UnkownRateErr
		}
	}

	if rate == "up" {
		userRating++
		return nil, userRating, nil
	}

	userRating--
	return nil, userRating, nil
}
