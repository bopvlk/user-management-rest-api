package interactor

import (
	"context"
	"crypto/sha1"
	"fmt"
	"time"

	"git.foxminded.com.ua/3_REST_API/interal"
	"git.foxminded.com.ua/3_REST_API/interal/domain/models"
	"git.foxminded.com.ua/3_REST_API/interal/usecase/repository"
	"github.com/golang-jwt/jwt/v4"
)

const CtxUserKey = "user"

type UserInteractor interface {
	SignUp(ctx context.Context, user *models.User) (*time.Duration, string, error)
	SignIn(ctx context.Context, user *models.User) (*time.Duration, string, error)
	DeleteSigner(ctx context.Context, user *models.User) error
	FindSignerByID(ctx context.Context, user *models.User) (*models.User, error)
	FindAllSigners(ctx context.Context, users []*models.User) ([]*models.User, error)
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
	user.Password = uI.hashing(user.Password)

	user, err := uI.userRepo.CreateUser(ctx, user)
	if err != nil {
		return nil, "", interal.ErrCannotCreateUser
	}

	token, err := uI.makeSignedToken(user)

	return &uI.expireDuration, token, err
}

func (uI *userInteractor) SignIn(ctx context.Context, user *models.User) (*time.Duration, string, error) {
	user.Password = uI.hashing(user.Password)

	user, err := uI.userRepo.FindOneUser(ctx, user)
	if err != nil {
		return nil, "", interal.ErrUserNotFound
	}

	token, err := uI.makeSignedToken(user)

	return &uI.expireDuration, token, err
}

func (uI *userInteractor) DeleteSigner(ctx context.Context, user *models.User) error {

	if err := uI.userRepo.DeleteUser(ctx, user); err != nil {
		return err
	}
	return nil
}

func (uI *userInteractor) FindSignerByID(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := uI.userRepo.FindOneUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (uI *userInteractor) FindAllSigners(ctx context.Context, users []*models.User) ([]*models.User, error) {
	user, err := uI.userRepo.FindAllUsers(ctx, users)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (uI *userInteractor) hashing(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(uI.hashSalt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
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
