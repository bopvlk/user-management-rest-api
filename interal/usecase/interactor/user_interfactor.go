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
	SignIn(ctx context.Context, user *models.User) (string, error)
	DeleteSigner(ctx context.Context, user *models.User) error
	FindSignerByID(ctx context.Context, user *models.User) (*models.User, error)
	FindAllSigners(ctx context.Context, users []*models.User) ([]*models.User, error)
	// ParseSignToken(accessToken string) (*models.User, error)
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

	user, err := uI.userRepo.CreateUserData(ctx, user)
	if err != nil {
		return nil, "", interal.ErrCannotCreateUser
	}

	token, err := uI.newAuthclainms(user)

	return &uI.expireDuration, token, err
}

func (uI *userInteractor) SignIn(ctx context.Context, user *models.User) (string, error) {
	user.Password = uI.hashing(user.Password)

	user, err := uI.userRepo.FindOneUserData(ctx, user)
	if err != nil {
		return "", interal.ErrUserNotFound
	}

	return uI.newAuthclainms(user)
}

func (uI *userInteractor) DeleteSigner(ctx context.Context, user *models.User) error {
	//TODO
	if err := uI.userRepo.DeleteUserData(ctx, user); err != nil {
		return err
	}
	return nil
}

func (uI *userInteractor) FindSignerByID(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := uI.userRepo.FindOneUserData(ctx, user)
	if err != nil {
		return nil, err
	}
	return user, err
}

func (uI *userInteractor) FindAllSigners(ctx context.Context, users []*models.User) ([]*models.User, error) {
	user, err := uI.userRepo.FindAllUsersData(ctx, users)
	if err != nil {
		return nil, err
	}
	return user, err
}

// func (uI *userInteractor) ParseSignToken(ctx context.Context, accessToken string) (*models.User, error) {
// 	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{
// 		StandardClaims: jwt.StandardClaims{},
// 		User:           &models.User{},
// 	}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return uI.signingKey, nil
// 	})

// 	if err != nil {
// 		return nil, err
// 	}

// 	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
// 		return claims.User, nil
// 	}

// 	return nil, interal.ErrInvalidAccessToken
// }

func (uI *userInteractor) hashing(password string) string {
	pwd := sha1.New()
	pwd.Write([]byte(password))
	pwd.Write([]byte(uI.hashSalt))
	return fmt.Sprintf("%x", pwd.Sum(nil))
}

func (uI *userInteractor) newAuthclainms(user *models.User) (string, error) {
	claims := AuthClaims{
		User: user,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(uI.expireDuration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(uI.signingKey)
}
