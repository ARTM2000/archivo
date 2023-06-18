package auth

import "errors"

var (
	ErrUnhandled                  = errors.New("unhandled error")
	ErrAdminExist                 = errors.New("admin exists")
	ErrUserExist                  = errors.New("user exists")
	ErrEmailOrUsernameInUse       = errors.New("email or username is used")
	ErrEmailOrPasswordIsIncorrect = errors.New("email or password is incorrect")
)
