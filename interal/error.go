package interal

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidAccessToken = errors.New("invalid access token")
	ErrUsersMoreThanOne   = errors.New("users more than one")
	ErrCannotCreateUser   = errors.New("it's can't create user")
)
