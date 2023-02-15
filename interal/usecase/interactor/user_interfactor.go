package interactor

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"git.foxminded.com.ua/3_REST_API/interal/apperrors"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/interface/repository"
	"github.com/golang-jwt/jwt/v4"
)

const CtxUserKey = "user"

type UserInteractor interface {
	SignUp(ctx context.Context, user *models.User) (*time.Duration, string, error)
	SignIn(ctx context.Context, name, password string) (*time.Duration, string, error)
	DeleteSigner(ctx context.Context, user *models.User) error
	FindOneSigner(ctx context.Context, id uint) (*models.User, error)
	FindSigners(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, error)
}

type AuthClaims struct {
	jwt.RegisteredClaims
	User *models.User `json:"user"`
}

type userInteractor struct {
	userRepo       repository.UserRepository
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewUserInteractor(
	userRepo repository.UserRepository,
	hashSalt string,
	signingKey []byte,
	tokenTTL time.Duration) *userInteractor {
	return &userInteractor{
		userRepo:       userRepo,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTL,
	}
}

func (uI *userInteractor) SignUp(ctx context.Context, user *models.User) (*time.Duration, string, error) {
	var err error
	user.Password, err = uI.hashing(user.Password)
	if err != nil {
		return nil, "", apperrors.HashingPasswordErr.AppendMessage(err)
	}

	user, err = uI.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, "", apperrors.CanNotCreateUserErr.AppendMessage(err)
	}

	token, err := uI.makeSignedToken(user)
	if err != nil {
		return nil, "", apperrors.CanNotCreateTokenErr.AppendMessage(err)
	}

	return &uI.expireDuration, token, nil
}

func (uI *userInteractor) SignIn(ctx context.Context, name, password string) (*time.Duration, string, error) {
	password, err := uI.hashing(password)
	if err != nil {
		return nil, "", apperrors.HashingPasswordErr.AppendMessage(err)
	}

	user, err := uI.userRepo.FindOneUserByUserNameAndPassword(ctx, name, password)
	if err != nil {
		return nil, "", apperrors.UserNotFoundErr.AppendMessage(err)
	}

	token, err := uI.makeSignedToken(user)
	if err != nil {
		return nil, "", apperrors.CanNotCreateTokenErr.AppendMessage(err)
	}

	return &uI.expireDuration, token, nil
}

func (uI *userInteractor) DeleteSigner(ctx context.Context, user *models.User) error {

	if err := uI.userRepo.DeleteUser(ctx, user); err != nil {
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
		return nil, nil, apperrors.PaginationErr.AppendMessage(err)
	}
	return pagination, users, nil
}

func (uI *userInteractor) hashing(password string) (string, error) {
	pwd := sha1.New()
	if _, err := pwd.Write([]byte(password)); err != nil {
		return "", err
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(uI.expireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uI.signingKey)
}
