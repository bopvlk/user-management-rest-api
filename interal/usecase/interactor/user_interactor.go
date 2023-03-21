package interactor

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/interface/repository"
	"github.com/golang-jwt/jwt/v4"
)

type UserInteractor interface {
	SignUp(ctx context.Context, user *models.User) (int, string, error)
	SignIn(ctx context.Context, name, password string) (int, string, error)
	FindOneSigner(ctx context.Context, id uint) (*models.User, error)
	FindSigners(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, error)
	DeleteSignerByID(ctx context.Context, id int) error
	DeleteOwnSignIn(ctx context.Context, id int) error
	UpdateSignersByID(ctx context.Context, id int, user *models.User) (*models.User, error)
	UpdateOwnSignIn(ctx context.Context, id int, user *models.User) (*models.User, error)
	RateUser(ctx context.Context, myID uint, username, rate string) (*models.User, error)
}

type AuthClaims struct {
	jwt.RegisteredClaims
	User *models.User `json:"user"`
}

type userInteractor struct {
	userRepo       repository.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration int
}

func NewUserInteractor(userRepo repository.UserRepository, hashSalt string, signingKey []byte, tokenTTL int) *userInteractor {
	return &userInteractor{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: tokenTTL,
	}
}

func (uI *userInteractor) SignUp(ctx context.Context, user *models.User) (int, string, error) {
	var err error
	user.Password, err = uI.hashing(user.Password)
	if err != nil {
		return 0, "", apperrors.HashingPasswordErr.AppendMessage(err)
	}

	user, err = uI.userRepo.CreateUser(ctx, user)
	if err != nil {
		return 0, "", apperrors.CanNotCreateUserErr.AppendMessage(err)
	}

	token, err := uI.makeSignedToken(user)
	if err != nil {
		return 0, "", apperrors.CanNotCreateTokenErr.AppendMessage(err)
	}

	return uI.expireDuration, token, nil
}

func (uI *userInteractor) SignIn(ctx context.Context, name, password string) (int, string, error) {
	password, err := uI.hashing(password)
	if err != nil {
		return 0, "", apperrors.HashingPasswordErr.AppendMessage(err)
	}

	user, err := uI.userRepo.FindOneUserByUserNameAndPassword(ctx, name, password)
	if err != nil {
		return 0, "", apperrors.UserNotFoundErr.AppendMessage(err)
	}

	token, err := uI.makeSignedToken(user)
	if err != nil {
		return 0, "", apperrors.CanNotCreateTokenErr.AppendMessage(err)
	}

	return uI.expireDuration, token, nil
}

func (uI *userInteractor) DeleteSignerByID(ctx context.Context, id int) error {

	if err := uI.userRepo.DeleteUserByID(ctx, id); err != nil {
		if err == &apperrors.WrongRoleErr {
			return err
		}
		return apperrors.CanNotDeleteUserErr.AppendMessage(err)
	}
	return nil
}

func (uI *userInteractor) DeleteOwnSignIn(ctx context.Context, id int) error {

	if err := uI.userRepo.DeleteOwnUser(ctx, id); err != nil {
		return apperrors.CanNotDeleteUserErr.AppendMessage(err)
	}
	return nil
}

func (uI *userInteractor) FindOneSigner(ctx context.Context, id uint) (*models.User, error) {
	user, err := uI.userRepo.FindOneUserByID(ctx, id)
	if err != nil {
		return nil, apperrors.UserNotFoundErr.AppendMessage(err)
	}
	return user, nil
}

func (uI *userInteractor) FindSigners(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, error) {
	pagination, users, err := uI.userRepo.FindUsers(ctx, pagination)
	if err != nil {
		if err == &apperrors.WrongRoleErr {
			return nil, nil, err
		}
		return nil, nil, apperrors.PaginationErr.AppendMessage(err)
	}

	return pagination, users, nil
}

func (uI *userInteractor) UpdateSignersByID(ctx context.Context, id int, user *models.User) (*models.User, error) {
	user, err := uI.userRepo.UpdateUserByID(ctx, id, user)
	if err != nil {
		return nil, apperrors.CanNotUpdateErr.AppendMessage(err)
	}

	return user, nil
}

func (uI *userInteractor) UpdateOwnSignIn(ctx context.Context, id int, user *models.User) (*models.User, error) {
	user, err := uI.userRepo.UpdateOwnUser(ctx, id, user)
	if err != nil {
		return nil, apperrors.CanNotUpdateErr.AppendMessage(err)
	}

	return user, nil
}

func (uI *userInteractor) RateUser(ctx context.Context, myID uint, username, rate string) (*models.User, error) {
	user, err := uI.userRepo.RateUserByUsername(ctx, myID, username, rate)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (uI *userInteractor) hashing(password string) (string, error) {
	if password == "" {
		return "", errors.New("empty pasword field")
	}
	pwd := sha1.New()
	if _, err := pwd.Write([]byte(password)); err != nil {
		return "", err
	}

	if uI.hashSalt == "" {
		return "", errors.New("empty hashSalt field")
	}
	if _, err := pwd.Write([]byte(uI.hashSalt)); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", pwd.Sum(nil)), nil
}

func (uI *userInteractor) makeSignedToken(user *models.User) (string, error) {
	claims := AuthClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Second * (time.Duration(uI.expireDuration)))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uI.signingKey)
}
