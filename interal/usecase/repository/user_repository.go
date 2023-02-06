package repository

import (
	"context"

	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	FindAllUsers(ctx context.Context, users []*models.User) ([]*models.User, error)
	FindOneUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, user *models.User) error
}
