package xerrors

import "errors"

var (
	ErrUnhandled                  = errors.New("unhandled error")
	ErrRecordNotFound             = errors.New("record not found")
	ErrDuplicateViolation         = errors.New("duplication violence")
	ErrAdminExist                 = errors.New("admin exists")
	ErrUserExist                  = errors.New("user exists")
	ErrEmailOrUsernameInUse       = errors.New("email or username is used")
	ErrEmailOrPasswordIsIncorrect = errors.New("email or password is incorrect")
	ErrUnauthorized               = errors.New("unauthorized")
)
