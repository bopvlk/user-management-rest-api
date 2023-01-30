package repository

import (
	"context"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
)

type UserRepository interface {
	CreateUserData(ctx context.Context, user *models.User) (*models.User, error)
	FindAllUsersData(ctx context.Context, users []*models.User) ([]*models.User, error)
	FindOneUserData(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUserData(ctx context.Context, user *models.User) error
}
