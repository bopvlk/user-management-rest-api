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

//go:generate mockgen -source=user_interfactor.go -destination=moks/mock.go

const CtxUserKey = "user"

type UserInteractor interface {
	SignUp(ctx context.Context, user *models.User) (int, string, *apperrors.AppError)
	SignIn(ctx context.Context, name, password string) (int, string, *apperrors.AppError)
	DeleteSigner(ctx context.Context, user *models.User) *apperrors.AppError
	FindOneSigner(ctx context.Context, id uint) (*models.User, *apperrors.AppError)
	FindSigners(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, *apperrors.AppError)
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

func (uI *userInteractor) SignUp(ctx context.Context, user *models.User) (int, string, *apperrors.AppError) {
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

func (uI *userInteractor) SignIn(ctx context.Context, name, password string) (int, string, *apperrors.AppError) {
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

func (uI *userInteractor) DeleteSigner(ctx context.Context, user *models.User) *apperrors.AppError {

	if err := uI.userRepo.DeleteUser(ctx, user); err != nil {
		return apperrors.CanNotDeleteUserErr.AppendMessage(err)
	}
	return nil
}

func (uI *userInteractor) FindOneSigner(ctx context.Context, id uint) (*models.User, *apperrors.AppError) {
	user, err := uI.userRepo.FindOneUserByID(ctx, id)
	if err != nil {
		return nil, apperrors.UserNotFoundErr.AppendMessage(err)
	}
	return user, nil
}

func (uI *userInteractor) FindSigners(ctx context.Context, pagination *models.Pagination) (*models.Pagination, []*models.User, *apperrors.AppError) {

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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(uI.expireDuration))),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uI.signingKey)
}
